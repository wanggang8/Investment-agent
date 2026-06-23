package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"investment-agent/internal/application/dto"
)

func TestAdjustPortfolioPersistsAdjustReasonInAudit(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/adjustments", bytes.NewBufferString(`{"cash":80,"total_assets":80,"adjust_reason":"与券商账户手动核对后校准","positions":[]}`))
	req.Header.Set("X-Request-ID", "req_portfolio_adjust")
	w := httptest.NewRecorder()

	app.AdjustPortfolio(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var inputRefType, inputRef string
	if err := db.QueryRow(`SELECT COALESCE(input_ref_type,''),COALESCE(input_ref,'') FROM audit_events WHERE request_id='req_portfolio_adjust' ORDER BY created_at DESC LIMIT 1`).Scan(&inputRefType, &inputRef); err != nil {
		t.Fatalf("read audit: %v", err)
	}
	if inputRefType != "adjust_reason" || inputRef != "与券商账户手动核对后校准" {
		t.Fatalf("expected adjust reason audit ref, got type=%q ref=%q", inputRefType, inputRef)
	}
}

func TestAdjustPortfolioReplacesCurrentPositions(t *testing.T) {
	app, db := testApp(t)
	seed := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/init", bytes.NewBufferString(`{"cash":50,"total_assets":100,"positions":[{"symbol":"510300","name":"沪深300","quantity":10,"cost_price":2,"current_price":3,"buy_reason":"低估配置"},{"symbol":"159915","name":"创业板","quantity":5,"cost_price":4,"current_price":4,"buy_reason":"分散配置"}]}`))
	seed.Header.Set("X-Request-ID", "req_portfolio_seed")
	seedW := httptest.NewRecorder()
	app.InitPortfolio(seedW, seed)
	if seedW.Code != http.StatusOK {
		t.Fatalf("expected seed 200, got %d body=%s", seedW.Code, seedW.Body.String())
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/adjustments", bytes.NewBufferString(`{"cash":100,"total_assets":100,"adjust_reason":"清仓后校准","positions":[]}`))
	req.Header.Set("X-Request-ID", "req_portfolio_replace")
	w := httptest.NewRecorder()
	app.AdjustPortfolio(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var current int
	if err := db.QueryRow(`SELECT COUNT(*) FROM positions`).Scan(&current); err != nil {
		t.Fatalf("count current positions: %v", err)
	}
	if current != 0 {
		t.Fatalf("expected current positions to be replaced, got %d", current)
	}
}

func TestAdjustPortfolioPersistsBuyDateAndPositionState(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/adjustments", bytes.NewBufferString(`{"cash":1000,"total_assets":1800,"adjust_reason":"P87 状态校准","positions":[{"symbol":"510300","name":"沪深300ETF","quantity":100,"cost_price":3,"current_price":4,"buy_date":"2026-01-05","position_state":"sell_only","buy_reason":"买入逻辑破坏后只卖不买","asset_tag":"core"},{"symbol":"159915","name":"创业板ETF","quantity":100,"cost_price":2,"current_price":4,"buy_date":"2026-01-06","position_state":"frozen_watch","buy_reason":"多源验证不足冻结观察","asset_tag":"satellite"}]}`))
	req.Header.Set("X-Request-ID", "req_p87_state")
	w := httptest.NewRecorder()

	app.AdjustPortfolio(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var sellOnly, frozen int
	if err := db.QueryRow(`SELECT COUNT(*) FROM positions WHERE symbol='510300' AND buy_date='2026-01-05' AND position_state='sell_only' AND asset_tag='core'`).Scan(&sellOnly); err != nil {
		t.Fatalf("count sell only position: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM position_snapshots WHERE symbol='159915' AND buy_date='2026-01-06' AND position_state='frozen_watch' AND asset_tag='satellite'`).Scan(&frozen); err != nil {
		t.Fatalf("count frozen snapshot: %v", err)
	}
	if sellOnly != 1 || frozen != 1 {
		t.Fatalf("expected persisted position state/buy date, sellOnly=%d frozen=%d", sellOnly, frozen)
	}
	var highRiskRatio float64
	if err := db.QueryRow(`SELECT high_risk_ratio FROM portfolio_snapshots ORDER BY snapshot_time DESC LIMIT 1`).Scan(&highRiskRatio); err != nil {
		t.Fatalf("read high risk ratio: %v", err)
	}
	if highRiskRatio < 0.44 || highRiskRatio > 0.45 {
		t.Fatalf("expected high risk ratio from sell_only/frozen_watch market value, got %.6f", highRiskRatio)
	}
}

func TestAdjustPortfolioRejectsInvalidPositionState(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/adjustments", bytes.NewBufferString(`{"cash":70,"total_assets":100,"adjust_reason":"bad state","positions":[{"symbol":"510300","name":"沪深300ETF","quantity":10,"cost_price":2,"current_price":3,"position_state":"auto_trade","buy_reason":"低估配置"}]}`))
	req.Header.Set("X-Request-ID", "req_bad_position_state")
	w := httptest.NewRecorder()

	app.AdjustPortfolio(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d body=%s", w.Code, w.Body.String())
	}
	var positions int
	if err := db.QueryRow(`SELECT COUNT(*) FROM positions`).Scan(&positions); err != nil {
		t.Fatalf("count positions: %v", err)
	}
	if positions != 0 {
		t.Fatalf("invalid state should not write positions, got %d", positions)
	}
}

func TestPortfolioOfflineTransactionHandlerWritesLocalFacts(t *testing.T) {
	app, db := testApp(t)
	seed := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/init", bytes.NewBufferString(`{"cash":100,"total_assets":100,"positions":[]}`))
	seed.Header.Set("X-Request-ID", "req_seed_cash")
	seedW := httptest.NewRecorder()
	app.InitPortfolio(seedW, seed)
	if seedW.Code != http.StatusOK {
		t.Fatalf("expected seed 200, got %d body=%s", seedW.Code, seedW.Body.String())
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/offline-transactions", bytes.NewBufferString(`{"operation_type":"buy","symbol":"510300","name":"沪深300ETF","quantity":10,"price":3,"fees":1,"executed_at":"2026-05-29T03:00:00Z","buy_reason":"低估配置","note":"仅记录线下动作"}`))
	req.Header.Set("X-Request-ID", "req_offline_handler")
	w := httptest.NewRecorder()
	app.RecordOfflineTransaction(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var txCount, auditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM position_transactions WHERE symbol='510300'`).Scan(&txCount); err != nil {
		t.Fatalf("count transactions: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE request_id='req_offline_handler'`).Scan(&auditCount); err != nil {
		t.Fatalf("count audit: %v", err)
	}
	if txCount != 1 || auditCount != 1 {
		t.Fatalf("expected transaction and audit facts, tx=%d audit=%d", txCount, auditCount)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("不连接券商")) {
		t.Fatalf("expected non-trading safety text, body=%s", w.Body.String())
	}
}

func TestPortfolioHoldingEditAndRemoveHandlers(t *testing.T) {
	app, db := testApp(t)
	seed := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/init", bytes.NewBufferString(`{"cash":70,"total_assets":100,"positions":[{"symbol":"510300","name":"沪深300ETF","quantity":10,"cost_price":2,"current_price":3,"buy_reason":"低估配置"}]}`))
	seed.Header.Set("X-Request-ID", "req_seed_holding")
	seedW := httptest.NewRecorder()
	app.InitPortfolio(seedW, seed)
	if seedW.Code != http.StatusOK {
		t.Fatalf("expected seed 200, got %d body=%s", seedW.Code, seedW.Body.String())
	}
	var positionID string
	if err := db.QueryRow(`SELECT position_id FROM positions WHERE symbol='510300'`).Scan(&positionID); err != nil {
		t.Fatalf("read position id: %v", err)
	}

	editReq := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/holdings", bytes.NewBufferString(`{"position_id":"`+positionID+`","reason":"本地校准","confirmation":"confirmed","position":{"symbol":"510300","name":"沪深300ETF","quantity":8,"cost_price":2,"current_price":4,"buy_reason":"低估配置"}}`))
	editReq.Header.Set("X-Request-ID", "req_edit_holding")
	editW := httptest.NewRecorder()
	app.EditHolding(editW, editReq)
	if editW.Code != http.StatusOK {
		t.Fatalf("expected edit 200, got %d body=%s", editW.Code, editW.Body.String())
	}

	removeReq := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/holdings/remove", bytes.NewBufferString(`{"position_id":"`+positionID+`","reason":"清仓后校准","confirmation":"confirmed"}`))
	removeReq.Header.Set("X-Request-ID", "req_remove_holding")
	removeW := httptest.NewRecorder()
	app.RemoveHolding(removeW, removeReq)
	if removeW.Code != http.StatusOK {
		t.Fatalf("expected remove 200, got %d body=%s", removeW.Code, removeW.Body.String())
	}
	var current int
	if err := db.QueryRow(`SELECT COUNT(*) FROM positions`).Scan(&current); err != nil {
		t.Fatalf("count positions: %v", err)
	}
	if current != 0 {
		t.Fatalf("expected current holding removed, got %d", current)
	}
}

func TestPortfolioBatchImportValidateDoesNotWriteFacts(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/imports/validate", bytes.NewBufferString(`{"rows":[{"row_number":1,"row_type":"holding","symbol":"510300","name":"沪深300ETF","quantity":10,"cost_price":2,"current_price":3,"buy_reason":"低估配置"},{"row_number":2,"row_type":"holding","symbol":"","name":"坏数据","quantity":1,"cost_price":1,"current_price":1}]}`))
	req.Header.Set("X-Request-ID", "req_import_validate")
	w := httptest.NewRecorder()
	app.ValidatePortfolioImport(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	if !bytes.Contains(w.Body.Bytes(), []byte(`"invalid_count":1`)) {
		t.Fatalf("expected row validation errors, body=%s", w.Body.String())
	}
	var current int
	if err := db.QueryRow(`SELECT COUNT(*) FROM positions`).Scan(&current); err != nil {
		t.Fatalf("count positions: %v", err)
	}
	if current != 0 {
		t.Fatalf("validation should not write positions, got %d", current)
	}
	var batches int
	if err := db.QueryRow(`SELECT COUNT(*) FROM local_account_import_batches WHERE request_id='req_import_validate' AND status='validated'`).Scan(&batches); err != nil {
		t.Fatalf("count import batches: %v", err)
	}
	if batches != 1 {
		t.Fatalf("expected validation batch metadata, got %d", batches)
	}
}

func TestPortfolioBatchImportConfirmWritesFacts(t *testing.T) {
	app, db := testApp(t)
	validateReq := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/imports/validate", bytes.NewBufferString(`{"rows":[{"row_number":1,"row_type":"holding","symbol":"510300","name":"沪深300ETF","quantity":10,"cost_price":2,"current_price":3,"buy_date":"2026-01-05","position_state":"sell_only","buy_reason":"低估配置"}]}`))
	validateReq.Header.Set("X-Request-ID", "req_import_validate_confirm")
	validateW := httptest.NewRecorder()
	app.ValidatePortfolioImport(validateW, validateReq)
	if validateW.Code != http.StatusOK {
		t.Fatalf("expected validate 200, got %d body=%s", validateW.Code, validateW.Body.String())
	}
	var validateBody struct {
		Data struct {
			ImportBatchID string `json:"import_batch_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(validateW.Body.Bytes(), &validateBody); err != nil {
		t.Fatalf("decode validate response: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/imports/confirm", bytes.NewBufferString(`{"import_batch_id":"`+validateBody.Data.ImportBatchID+`","confirm_reason":"导入初始持仓","rows":[{"row_number":1,"row_type":"holding","symbol":"510300","name":"沪深300ETF","quantity":10,"cost_price":2,"current_price":3,"buy_date":"2026-01-05","position_state":"sell_only","buy_reason":"低估配置"}]}`))
	req.Header.Set("X-Request-ID", "req_import_confirm")
	w := httptest.NewRecorder()
	app.ConfirmPortfolioImport(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var positions, audits int
	if err := db.QueryRow(`SELECT COUNT(*) FROM positions WHERE symbol='510300' AND buy_date='2026-01-05' AND position_state='sell_only'`).Scan(&positions); err != nil {
		t.Fatalf("count positions: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE request_id='req_import_confirm'`).Scan(&audits); err != nil {
		t.Fatalf("count audits: %v", err)
	}
	if positions != 1 || audits != 1 {
		t.Fatalf("expected imported position and audit, positions=%d audits=%d", positions, audits)
	}
}

func TestPortfolioBatchImportConfirmWritesTransactionRows(t *testing.T) {
	app, db := testApp(t)
	seed := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/init", bytes.NewBufferString(`{"cash":100,"total_assets":100,"positions":[]}`))
	seed.Header.Set("X-Request-ID", "req_import_tx_seed")
	seedW := httptest.NewRecorder()
	app.InitPortfolio(seedW, seed)
	if seedW.Code != http.StatusOK {
		t.Fatalf("expected seed 200, got %d body=%s", seedW.Code, seedW.Body.String())
	}

	validateReq := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/imports/validate", bytes.NewBufferString(`{"rows":[{"row_number":1,"row_type":"transaction","operation_type":"buy","symbol":"510300","name":"沪深300ETF","quantity":10,"price":3,"fees":1,"occurred_at":"2026-05-29T03:00:00Z","buy_reason":"低估配置"}]}`))
	validateReq.Header.Set("X-Request-ID", "req_import_tx_validate")
	validateW := httptest.NewRecorder()
	app.ValidatePortfolioImport(validateW, validateReq)
	if validateW.Code != http.StatusOK {
		t.Fatalf("expected validate 200, got %d body=%s", validateW.Code, validateW.Body.String())
	}
	var validateBody struct {
		Data struct {
			ImportBatchID string `json:"import_batch_id"`
		} `json:"data"`
	}
	if err := json.Unmarshal(validateW.Body.Bytes(), &validateBody); err != nil {
		t.Fatalf("decode validate response: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/imports/confirm", bytes.NewBufferString(`{"import_batch_id":"`+validateBody.Data.ImportBatchID+`","confirm_reason":"导入线下交易","rows":[{"row_number":1,"row_type":"transaction","operation_type":"buy","symbol":"510300","name":"沪深300ETF","quantity":10,"price":3,"fees":1,"occurred_at":"2026-05-29T03:00:00Z","buy_reason":"低估配置"}]}`))
	req.Header.Set("X-Request-ID", "req_import_tx_confirm")
	w := httptest.NewRecorder()
	app.ConfirmPortfolioImport(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var txCount, positionCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM position_transactions WHERE symbol='510300'`).Scan(&txCount); err != nil {
		t.Fatalf("count transactions: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM positions WHERE symbol='510300'`).Scan(&positionCount); err != nil {
		t.Fatalf("count positions: %v", err)
	}
	if txCount != 1 || positionCount != 1 {
		t.Fatalf("expected transaction row import to write transaction and position, tx=%d positions=%d", txCount, positionCount)
	}
}

func TestP88QuarterlyRebalanceReviewCalculatesManualActionsOnly(t *testing.T) {
	app, db := testApp(t)
	seed := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/init", bytes.NewBufferString(`{"cash":100,"total_assets":1000,"positions":[{"symbol":"510300","name":"沪深300ETF","quantity":100,"cost_price":3,"current_price":6,"buy_reason":"核心配置","asset_tag":"core"},{"symbol":"159915","name":"创业板ETF","quantity":100,"cost_price":2,"current_price":3,"buy_reason":"卫星配置","asset_tag":"satellite"}]}`))
	seed.Header.Set("X-Request-ID", "req_p88_rebalance_seed")
	seedW := httptest.NewRecorder()
	app.InitPortfolio(seedW, seed)
	if seedW.Code != http.StatusOK {
		t.Fatalf("expected seed 200, got %d body=%s", seedW.Code, seedW.Body.String())
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/rebalance-review", bytes.NewBufferString(`{"target_core_ratio":0.2,"target_satellite_ratio":0.5,"target_cash_ratio":0.3,"drift_threshold":0.15,"review_date":"2026-06-22"}`))
	req.Header.Set("X-Request-ID", "req_p88_rebalance")
	w := httptest.NewRecorder()
	app.ReviewQuarterlyRebalance(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var body struct {
		Data dto.RebalanceReviewResponse `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Data.Items) != 3 || body.Data.DriftThreshold != 0.15 || body.Data.SafetyStatement == "" {
		t.Fatalf("unexpected rebalance response: %+v", body.Data)
	}
	byBucket := map[string]dto.RebalanceReviewItem{}
	for _, item := range body.Data.Items {
		byBucket[item.Bucket] = item
	}
	if byBucket["core"].Recommendation != "manual_sell_or_reduce" || byBucket["core"].ManualAmount <= 0 {
		t.Fatalf("expected manual core sell/reduce recommendation: %+v", byBucket["core"])
	}
	if byBucket["satellite"].Recommendation != "manual_buy_or_add" || byBucket["satellite"].ManualAmount <= 0 {
		t.Fatalf("expected manual satellite buy/add recommendation: %+v", byBucket["satellite"])
	}
	if byBucket["cash"].Recommendation != "manual_raise_cash" || byBucket["cash"].ManualAmount <= 0 {
		t.Fatalf("expected manual cash raise recommendation: %+v", byBucket["cash"])
	}
	var txCount, auditCount int
	if err := db.QueryRow(`SELECT COUNT(*) FROM position_transactions`).Scan(&txCount); err != nil {
		t.Fatalf("count transactions: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE request_id='req_p88_rebalance' AND action='run_local_task' AND input_ref_type='rebalance_review'`).Scan(&auditCount); err != nil {
		t.Fatalf("count audit: %v", err)
	}
	if txCount != 0 || auditCount != 1 {
		t.Fatalf("rebalance review must only write audit, tx=%d audit=%d", txCount, auditCount)
	}
}

func TestPortfolioCorrectionHandlerWritesCorrectionAndAudit(t *testing.T) {
	app, db := testApp(t)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/portfolio/corrections", bytes.NewBufferString(`{"target_type":"position","target_id":"pos_manual","before_json":"{\"quantity\":10}","after_json":"{\"quantity\":8}","correction_reason":"录入数量修正"}`))
	req.Header.Set("X-Request-ID", "req_correction")
	w := httptest.NewRecorder()
	app.CorrectPortfolioFact(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	var corrections, audits int
	if err := db.QueryRow(`SELECT COUNT(*) FROM local_account_corrections WHERE target_id='pos_manual'`).Scan(&corrections); err != nil {
		t.Fatalf("count corrections: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM audit_events WHERE request_id='req_correction'`).Scan(&audits); err != nil {
		t.Fatalf("count audits: %v", err)
	}
	if corrections != 1 || audits != 1 {
		t.Fatalf("expected correction and audit, corrections=%d audits=%d", corrections, audits)
	}
}
