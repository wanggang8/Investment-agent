package workflow

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"

	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/idgen"
)

func TestFixtureSentimentProxyCollectorNormalizesP34Payload(t *testing.T) {
	collector := FixtureSentimentProxyCollector{Fixtures: map[string]SentimentProxyPoint{
		"000300": {SourceName: "eastmoney_sentiment_proxy", SourceLevel: model.SourceLevelB, DataDate: "2026-06-05", HeatScore: 64, SentimentState: model.SentimentNeutral, Raw: map[string]any{"rank": "沪深300"}},
	}}

	point, err := collector.FetchMarketData(context.Background(), "000300")
	if err != nil {
		t.Fatalf("FetchMarketData: %v", err)
	}
	if point.SourceName != "eastmoney_sentiment_proxy" || point.SourceType != "sentiment_proxy" || point.SourceLevel != model.SourceLevelB || point.TradeDate != "2026-06-05" {
		t.Fatalf("unexpected source fields: %+v", point)
	}
	if point.SentimentState != model.SentimentNeutral || point.Metadata["data_category"] != "sentiment_proxy" || point.Metadata["heat_score"] != float64(64) {
		t.Fatalf("unexpected sentiment proxy metadata: %+v", point)
	}
	health, ok := point.Metadata["p34_source_health"].(map[string]any)
	sentimentHealth, _ := health["sentiment_proxy"].(map[string]any)
	if !ok || sentimentHealth["freshness"] != "stubbed" || sentimentHealth["source_level"] != string(model.SourceLevelB) || sentimentHealth["failure_category"] != "stubbed" || sentimentHealth["last_success_at"] != nil || sentimentHealth["last_failure_at"] == nil {
		t.Fatalf("expected explicit stubbed sentiment proxy health, got %+v", point.Metadata)
	}
}

func TestMarketRefreshGraphPersistsP34SentimentProxy(t *testing.T) {
	deps := WorkflowDependencies{MarketDataSource: FixtureSentimentProxyCollector{Fixtures: map[string]SentimentProxyPoint{
		"000300": {SourceName: "eastmoney_sentiment_proxy", SourceLevel: model.SourceLevelB, DataDate: "2026-06-05", HeatScore: 64, SentimentState: model.SentimentNeutral},
	}}}

	out, err := NewMarketRefreshGraphWithDependencies(deps).Run(context.Background(), MarketRefreshInput{RequestID: "req_graph_market_p34_sentiment", Symbol: "000300"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	var metrics map[string]any
	if err := json.Unmarshal([]byte(out.MarketSnapshot.MarketMetricsJSON), &metrics); err != nil {
		t.Fatalf("unmarshal metrics: %v", err)
	}
	metadata, ok := metrics["metadata"].(map[string]any)
	if !ok || metadata["data_category"] != "sentiment_proxy" || metadata["heat_score"] != float64(64) {
		t.Fatalf("expected P34 sentiment metadata, got %+v", metrics)
	}
	if out.MarketSnapshot.SentimentState != model.SentimentNeutral || out.MarketSnapshot.ClosePrice != 0 {
		t.Fatalf("sentiment proxy must not fabricate price data: %+v", out.MarketSnapshot)
	}
}

type failingMarketRepo struct{}

func (f failingMarketRepo) SaveMarketSnapshot(context.Context, model.MarketSnapshot, string) error {
	return apperr.New(apperr.CodeInternalError, apperr.CategoryInternal, "write failed")
}

func (f failingMarketRepo) GetMarketSnapshot(context.Context, string) (model.MarketSnapshot, error) {
	return model.MarketSnapshot{}, sql.ErrNoRows
}
func (f failingMarketRepo) GetLatestMarketSnapshot(context.Context) (model.MarketSnapshot, error) {
	return model.MarketSnapshot{}, sql.ErrNoRows
}
func (f failingMarketRepo) GetLatestMarketSnapshotBySymbol(context.Context, string) (model.MarketSnapshot, error) {
	return model.MarketSnapshot{}, sql.ErrNoRows
}

func (f failingMarketRepo) MarketSnapshotExists(context.Context, string) (bool, error) {
	return false, nil
}

func TestMarketRefreshGraphPersistsP27SourceMetadata(t *testing.T) {
	deps := WorkflowDependencies{MarketDataSource: testMarketDataSource{point: MarketDataPoint{ClosePrice: 4.321, SourceName: "eastmoney_fund", SourceLevel: model.SourceLevelB, SourceType: "fund_nav", TradeDate: "2026-06-05", CapturedAt: "2026-06-05T21:30:00Z", ContentHash: "sha256:test", Metadata: map[string]any{"fund_name": "沪深300ETF", "accumulated_nav": 5.678}}}}

	out, err := NewMarketRefreshGraphWithDependencies(deps).Run(context.Background(), MarketRefreshInput{RequestID: "req_graph_market_p27", Symbol: "510300"})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	var metrics map[string]any
	if err := json.Unmarshal([]byte(out.MarketSnapshot.MarketMetricsJSON), &metrics); err != nil {
		t.Fatalf("unmarshal metrics: %v", err)
	}
	if metrics["source_name"] != "eastmoney_fund" || metrics["source_level"] != string(model.SourceLevelB) || metrics["source_type"] != "fund_nav" || metrics["trade_date"] != "2026-06-05" || metrics["captured_at"] != "2026-06-05T21:30:00Z" || metrics["content_hash"] != "sha256:test" {
		t.Fatalf("source metadata not persisted: %+v", metrics)
	}
	metadata, ok := metrics["metadata"].(map[string]any)
	if !ok || metadata["fund_name"] != "沪深300ETF" || metadata["accumulated_nav"] != 5.678 {
		t.Fatalf("collector metadata not persisted: %+v", metrics)
	}
}

func TestMarketRefreshGraphDoesNotBackfillPercentilesForP27Source(t *testing.T) {
	deps := WorkflowDependencies{MarketDataSource: testMarketDataSource{point: MarketDataPoint{ClosePrice: 4.321, SourceName: "eastmoney_fund", SourceLevel: model.SourceLevelB, SourceType: "fund_nav", TradeDate: "2026-06-05"}}}

	out, err := NewMarketRefreshGraphWithDependencies(deps).Run(context.Background(), MarketRefreshInput{RequestID: "req_graph_market_no_backfill", Symbol: "510300", PEPercentile: 50, PBPercentile: 50})
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if out.MarketSnapshot.PEPercentile != 0 || out.MarketSnapshot.PBPercentile != 0 {
		t.Fatalf("P27 source must not backfill valuation percentiles: %+v", out.MarketSnapshot)
	}
}

func TestMarketRefreshGraphDedupesP27SourceTradeDate(t *testing.T) {
	store, err := appsqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer store.Close()
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	repos := repository.Repositories{MarketRepo: appsqlite.NewMarketRepository(store.DB), AuditRepo: appsqlite.NewAuditRepository(store.DB)}
	deps := NewWorkflowDependencies(repos, appsqlite.NewTransactor(store.DB))
	deps.MarketDataSource = testMarketDataSource{point: MarketDataPoint{ClosePrice: 4.321, SourceName: "eastmoney_fund", SourceLevel: model.SourceLevelB, SourceType: "fund_nav", TradeDate: "2026-06-05", Metadata: map[string]any{"fund_name": "沪深300ETF"}}}

	for _, requestID := range []string{"req_graph_market_dedupe_1", "req_graph_market_dedupe_2"} {
		if _, err := NewMarketRefreshGraphWithDependencies(deps).Run(context.Background(), MarketRefreshInput{RequestID: requestID, Symbol: "510300"}); err != nil {
			t.Fatalf("Run %s: %v", requestID, err)
		}
	}

	var snapshots int
	if err := store.DB.QueryRow(`SELECT COUNT(*) FROM market_snapshots WHERE symbol='510300'`).Scan(&snapshots); err != nil {
		t.Fatalf("count snapshots: %v", err)
	}
	if snapshots != 1 {
		t.Fatalf("expected one deduped snapshot, got %d", snapshots)
	}
}

func TestMarketRefreshGraphWritesFailureAuditOnSnapshotWriteError(t *testing.T) {
	deps := WorkflowDependencies{MarketRepo: failingMarketRepo{}}
	out, err := NewMarketRefreshGraphWithDependencies(deps).Run(context.Background(), MarketRefreshInput{RequestID: "req_graph_market_failed", Symbol: "510300"})
	if !apperr.IsCode(err, apperr.CodeMarketSnapshotWriteFailed) {
		t.Fatalf("expected MARKET_SNAPSHOT_WRITE_FAILED, got %v", err)
	}
	if len(out.AuditEvents) != 1 || out.AuditEvents[0].Status != model.AuditStatusFailed || out.AuditEvents[0].ErrorCode != string(apperr.CodeMarketSnapshotWriteFailed) {
		t.Fatalf("expected write failure audit, got %+v", out.AuditEvents)
	}
}

func TestMarketRefreshGraphPersistsSingleFailureAuditOnSnapshotWriteError(t *testing.T) {
	store, err := appsqlite.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	defer store.Close()
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	transactor := appsqlite.NewTransactor(store.DB)
	repos := repository.Repositories{MarketRepo: appsqlite.NewMarketRepository(store.DB), AuditRepo: appsqlite.NewAuditRepository(store.DB)}
	deps := NewWorkflowDependencies(repos, transactor)
	SetWorkflowIDGenerator(idgen.NewFixedGenerator(map[string][]string{"market": {"market_conflict"}, "audit": {"audit_failed"}}))
	defer SetWorkflowIDGenerator(idgen.NewGenerator())
	if err := repos.MarketRepo.SaveMarketSnapshot(context.Background(), model.MarketSnapshot{MarketSnapshotID: "market_conflict", Symbol: "510300", LiquidityState: model.LiquidityNormal, SentimentState: model.SentimentNeutral}, "2026-01-01T00:00:00Z"); err != nil {
		t.Fatalf("seed market: %v", err)
	}

	_, err = NewMarketRefreshGraphWithDependencies(deps).Run(context.Background(), MarketRefreshInput{RequestID: "req_graph_market_persist_failed", Symbol: "510300"})

	if !apperr.IsCode(err, apperr.CodeMarketSnapshotWriteFailed) {
		t.Fatalf("expected MARKET_SNAPSHOT_WRITE_FAILED, got %v", err)
	}
	var total, failed, success int
	if err := store.DB.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE request_id='req_graph_market_persist_failed' AND action='refresh_market_data'`).Scan(&total); err != nil {
		t.Fatalf("count audits: %v", err)
	}
	if err := store.DB.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE request_id='req_graph_market_persist_failed' AND action='refresh_market_data' AND status='failed' AND error_code=?`, string(apperr.CodeMarketSnapshotWriteFailed)).Scan(&failed); err != nil {
		t.Fatalf("count failed audits: %v", err)
	}
	if err := store.DB.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE request_id='req_graph_market_persist_failed' AND action='refresh_market_data' AND status='success'`).Scan(&success); err != nil {
		t.Fatalf("count success audits: %v", err)
	}
	if total != 1 || failed != 1 || success != 0 {
		t.Fatalf("expected one failed persisted audit and no success, total=%d failed=%d success=%d", total, failed, success)
	}
}
