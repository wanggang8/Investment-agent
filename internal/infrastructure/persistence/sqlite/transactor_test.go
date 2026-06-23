package sqlite

import (
	"context"
	"errors"
	"testing"

	"investment-agent/internal/domain/repository"
)

func TestTransactorCommitsWrites(t *testing.T) {
	db := testDB(t)
	tr := NewTransactor(db)
	ctx := context.Background()

	err := tr.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		return repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{
			AuditEventID: "audit_tx_commit",
			Actor:        "system",
			Action:       "generate_decision",
			Status:       "success",
			CreatedAt:    "2026-05-29T04:00:00Z",
		})
	})
	if err != nil {
		t.Fatalf("WithinTx error: %v", err)
	}

	if _, err := NewAuditRepository(db).GetAuditEvent(ctx, "audit_tx_commit"); err != nil {
		t.Fatalf("expected committed audit event: %v", err)
	}
}

func TestTransactorRollsBackWrites(t *testing.T) {
	db := testDB(t)
	tr := NewTransactor(db)
	ctx := context.Background()
	want := errors.New("force rollback")

	err := tr.WithinTx(ctx, func(ctx context.Context, repos repository.Repositories) error {
		if err := repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{
			AuditEventID: "audit_tx_rollback",
			Actor:        "system",
			Action:       "generate_decision",
			Status:       "success",
			CreatedAt:    "2026-05-29T04:00:00Z",
		}); err != nil {
			return err
		}
		return want
	})
	if !errors.Is(err, want) {
		t.Fatalf("WithinTx error = %v", err)
	}

	if _, err := NewAuditRepository(db).GetAuditEvent(ctx, "audit_tx_rollback"); err == nil {
		t.Fatal("expected audit event to be rolled back")
	}
}
