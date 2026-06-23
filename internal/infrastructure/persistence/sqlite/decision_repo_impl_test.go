package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

func TestDecisionRepositoryWriteReadAndRollback(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewDecisionRepository(db)

	decision := repository.DecisionRecord{
		DecisionID: "dec1", RequestID: "req1", WorkflowType: "daily_discipline",
		WorkflowStatus: "completed", RecordType: "formal_trade_advice", DashboardState: "normal",
		FinalVerdictStatus: "hold", FinalVerdictText: "持有", ConfirmationStatus: "pending",
		RuleVersion: "v3.0", CreatedAt: testTime,
	}
	refs := []repository.EvidenceRef{{
		EvidenceRefID: "eref1", EvidenceID: "ev1", DecisionID: "dec1", SummaryID: "sum1",
		SourceName: "source", SourceLevel: "A", EvidenceRole: "formal", Summary: "summary",
		IndependentSourceCount: 3, HighGradeIndependentSourceCount: 2, CreatedAt: testTime,
	}}
	if err := repo.SaveDecisionRecord(ctx, decision, refs); err != nil {
		t.Fatal(err)
	}
	got, gotRefs, err := repo.GetDecisionRecord(ctx, "dec1")
	if err != nil {
		t.Fatal(err)
	}
	if got.DecisionID != "dec1" || len(gotRefs) != 1 || gotRefs[0].IndependentSourceCount != 3 || gotRefs[0].HighGradeIndependentSourceCount != 2 {
		t.Fatalf("unexpected read: %#v %#v", got, gotRefs)
	}

	bad := decision
	bad.DecisionID = "dec_bad"
	bad.DashboardState = "invalid"
	if err := repo.SaveDecisionRecord(ctx, bad, nil); err == nil {
		t.Fatal("expected rollback error")
	}
	if _, _, err := repo.GetDecisionRecord(ctx, "dec_bad"); err == nil {
		t.Fatal("decision persisted after rollback")
	}
}

func TestOperationConfirmationWriteRead(t *testing.T) {
	db := testDB(t)
	repo := NewDecisionRepository(db)
	confirmation := repository.OperationConfirmation{
		ConfirmationID: "conf1", DecisionID: "dec1", ConfirmationType: "executed_manually", OperationType: "buy", Symbol: "510300", Quantity: 10, Price: 2.5, Fees: 1.2, CreatedAt: testTime,
	}
	if err := repo.SaveOperationConfirmation(context.Background(), confirmation); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetOperationConfirmation(context.Background(), "conf1")
	if err != nil {
		t.Fatal(err)
	}
	if got.ConfirmationType != "executed_manually" || !floatCloseLocal(got.Fees, 1.2) {
		t.Fatalf("got %#v", got)
	}
}

func TestDecisionRepositoryListsConfirmationsAndTransactions(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewDecisionRepository(db)

	for _, confirmation := range []repository.OperationConfirmation{
		{ConfirmationID: "conf_late", DecisionID: "dec_loop", ConfirmationType: "watch", CreatedAt: "2026-01-02T10:00:00Z"},
		{ConfirmationID: "conf_early", DecisionID: "dec_loop", ConfirmationType: "executed_manually", OperationType: "buy", Symbol: "510300", Quantity: 10, Price: 2.5, Fees: 1, ExecutedAt: "2026-01-02T09:30:00Z", PayloadJSON: `{"raw":"sensitive"}`, CreatedAt: "2026-01-02T09:00:00Z"},
		{ConfirmationID: "conf_other", DecisionID: "dec_other", ConfirmationType: "planned", CreatedAt: "2026-01-01T09:00:00Z"},
	} {
		if err := repo.SaveOperationConfirmation(ctx, confirmation); err != nil {
			t.Fatalf("save confirmation %s: %v", confirmation.ConfirmationID, err)
		}
	}
	for _, tx := range []repository.PositionTransaction{
		{TransactionID: "tx_late", ConfirmationID: "conf_early", Symbol: "510300", OperationType: "buy", Quantity: 4, Price: 2.6, Fees: 0.5, OccurredAt: "2026-01-02T10:30:00Z", BeforePositionJSON: `{"raw":"before"}`, AfterPositionJSON: `{"raw":"after"}`, CreatedAt: "2026-01-02T10:31:00Z"},
		{TransactionID: "tx_early", ConfirmationID: "conf_early", Symbol: "510300", OperationType: "buy", Quantity: 6, Price: 2.5, Fees: 0.5, OccurredAt: "2026-01-02T09:30:00Z", CreatedAt: "2026-01-02T09:31:00Z"},
		{TransactionID: "tx_other", ConfirmationID: "conf_other", Symbol: "510500", OperationType: "sell", Quantity: 1, Price: 3, OccurredAt: "2026-01-01T09:30:00Z", CreatedAt: "2026-01-01T09:31:00Z"},
	} {
		if err := repo.SavePositionTransaction(ctx, tx); err != nil {
			t.Fatalf("save transaction %s: %v", tx.TransactionID, err)
		}
	}

	confirmations, err := repo.ListOperationConfirmations(ctx, "dec_loop")
	if err != nil {
		t.Fatal(err)
	}
	if len(confirmations) != 2 {
		t.Fatalf("expected 2 confirmations, got %#v", confirmations)
	}
	if confirmations[0].ConfirmationID != "conf_early" || confirmations[1].ConfirmationID != "conf_late" {
		t.Fatalf("confirmations not ordered by created_at: %#v", confirmations)
	}
	if confirmations[0].PayloadJSON != "" {
		t.Fatalf("list confirmation should not load raw payload json: %#v", confirmations[0])
	}

	transactions, err := repo.ListPositionTransactionsByConfirmation(ctx, "conf_early")
	if err != nil {
		t.Fatal(err)
	}
	if len(transactions) != 2 {
		t.Fatalf("expected 2 transactions, got %#v", transactions)
	}
	if transactions[0].TransactionID != "tx_early" || transactions[1].TransactionID != "tx_late" {
		t.Fatalf("transactions not ordered by occurred_at: %#v", transactions)
	}
	if transactions[1].BeforePositionJSON != "" || transactions[1].AfterPositionJSON != "" {
		t.Fatalf("list transaction should not load raw position json: %#v", transactions[1])
	}
}

func TestDecisionRepositoryUpdateConfirmationStatusIfCurrent(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewDecisionRepository(db)
	decision := repository.DecisionRecord{
		DecisionID: "dec_confirm", RequestID: "req_confirm", WorkflowType: "daily_discipline",
		WorkflowStatus: "completed", RecordType: "formal_trade_advice", DashboardState: "normal",
		FinalVerdictStatus: "hold", FinalVerdictText: "持有", ConfirmationStatus: "pending",
		RuleVersion: "v3.0", CreatedAt: testTime,
	}
	if err := repo.SaveDecisionRecord(ctx, decision, nil); err != nil {
		t.Fatal(err)
	}
	updated, err := repo.UpdateDecisionConfirmationStatusIfCurrent(ctx, "dec_confirm", "watch", "executed_manually")
	if err != nil {
		t.Fatal(err)
	}
	if updated {
		t.Fatal("expected stale status update to be rejected")
	}
	_, status, err := repo.GetDecisionConfirmationState(ctx, "dec_confirm")
	if err != nil {
		t.Fatal(err)
	}
	if status != "pending" {
		t.Fatalf("status changed after stale update: %s", status)
	}
	updated, err = repo.UpdateDecisionConfirmationStatusIfCurrent(ctx, "dec_confirm", "pending", "planned")
	if err != nil {
		t.Fatal(err)
	}
	if !updated {
		t.Fatal("expected current status update to succeed")
	}
	_, status, err = repo.GetDecisionConfirmationState(ctx, "dec_confirm")
	if err != nil {
		t.Fatal(err)
	}
	if status != "planned" {
		t.Fatalf("unexpected status after update: %s", status)
	}
}

func floatCloseLocal(got, want float64) bool {
	if got > want {
		return got-want < 1e-9
	}
	return want-got < 1e-9
}

func TestDecisionRepositoryClassifiesErrors(t *testing.T) {
	db := testDB(t)
	repo := NewDecisionRepository(db)
	if _, _, err := repo.GetDecisionRecord(context.Background(), "missing_decision"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found decision error, got %v", err)
	}
	if _, err := repo.GetOperationConfirmation(context.Background(), "missing_confirmation"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found confirmation error, got %v", err)
	}
}
