package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"investment-agent/internal/application/handler"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
	"investment-agent/internal/infrastructure/wiring"
	"investment-agent/pkg/httputil"
	"investment-agent/pkg/logger"
)

// main 启动本地 HTTP 服务。P0 阶段只暴露健康检查，后续 API 在此注册路由。
func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	cfg, err := config.Load("")
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("validate config: %w", err)
	}

	log := logger.New(cfg.Log.Level)
	log.Info("starting http server", "addr", cfg.Server.Addr())

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/v1/health", healthHandler)

	store, err := appsqlite.Open(cfg.SQLite.Path)
	if err != nil {
		panic(fmt.Errorf("open sqlite: %w", err))
	}
	defer store.Close()
	if err := appsqlite.Migrate(context.Background(), store.DB); err != nil {
		panic(fmt.Errorf("migrate sqlite: %w", err))
	}
	transactor := appsqlite.NewTransactor(store.DB)
	repos := repository.Repositories{
		DecisionRepo:                  appsqlite.NewDecisionRepository(store.DB),
		AuditRepo:                     appsqlite.NewAuditRepository(store.DB),
		RuleRepo:                      appsqlite.NewRuleRepository(store.DB),
		MarketRepo:                    appsqlite.NewMarketRepository(store.DB),
		SettingsRepo:                  appsqlite.NewSettingsRepository(store.DB),
		IntelligenceRepo:              appsqlite.NewIntelligenceRepository(store.DB),
		NotificationRepo:              appsqlite.NewNotificationRepository(store.DB),
		PortfolioRepo:                 appsqlite.NewPortfolioRepository(store.DB),
		DailyAutoRunRepo:              appsqlite.NewDailyAutoRunRepository(store.DB),
		DailyDisciplineReportRepo:     appsqlite.NewDailyDisciplineReportRepository(store.DB),
		RiskAlertRepo:                 appsqlite.NewRiskAlertRepository(store.DB),
		RuleEffectRepo:                appsqlite.NewRuleEffectRepository(store.DB),
		DataQualityGateResolutionRepo: appsqlite.NewDataQualityGateResolutionRepository(store.DB),
	}
	deps := wiring.NewWorkflowDependencies(cfg, repos, transactor)
	handler.NewApp(deps, repos, transactor).RegisterRoutes(mux)
	stopDailyAutoRun := startDailyAutoRunScheduler(context.Background(), workflow.NewDailyAutoRunner(cfg.DailyAutoRun, deps), cfg.DailyAutoRun, log)
	defer stopDailyAutoRun()

	if err := http.ListenAndServe(cfg.Server.Addr(), mux); err != nil {
		return fmt.Errorf("listen: %w", err)
	}
	return nil
}

// healthHandler 用于本地开发和部署探活。
// P0 按契约返回简单 JSON，业务接口仍使用统一信封。
func healthHandler(w http.ResponseWriter, r *http.Request) {
	httputil.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func startDailyAutoRunScheduler(ctx context.Context, runner *workflow.DailyAutoRunner, cfg config.DailyAutoRunConfig, log *slog.Logger) func() {
	if !cfg.Enabled {
		return func() {}
	}
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		for {
			nextRun := nextDailyAutoRunTime(time.Now(), cfg)
			runner.SetNextRunAt(nextRun.Format(time.RFC3339))
			timer := time.NewTimer(time.Until(nextRun))
			select {
			case <-ctx.Done():
				timer.Stop()
				return
			case <-timer.C:
				runDailyAutoRun(ctx, runner, log)
			}
		}
	}()
	return cancel
}

func nextDailyAutoRunTime(now time.Time, cfg config.DailyAutoRunConfig) time.Time {
	loc, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		loc = time.Local
	}
	hour, minute := 0, 0
	if parsed, err := time.Parse("15:04", cfg.RunTime); err == nil {
		hour, minute = parsed.Hour(), parsed.Minute()
	}
	localNow := now.In(loc)
	next := time.Date(localNow.Year(), localNow.Month(), localNow.Day(), hour, minute, 0, 0, loc)
	if !next.After(localNow) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

func runDailyAutoRun(ctx context.Context, runner *workflow.DailyAutoRunner, log *slog.Logger) {
	if _, err := runner.RunOnce(ctx, time.Now()); err != nil {
		log.Warn("daily auto-run failed", "error", err)
	}
}
