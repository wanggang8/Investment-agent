package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/repository"
)

func TestAuditRepositoryAppendReadAndConstraint(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewAuditRepository(db)

	event := repository.AuditEvent{
		AuditEventID: "audit1", RequestID: "req1", Actor: "system", Action: "generate_decision",
		Status: "success", CreatedAt: testTime,
	}
	if err := repo.AppendAuditEvent(ctx, event); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetAuditEvent(ctx, "audit1")
	if err != nil {
		t.Fatal(err)
	}
	if got.Action != "generate_decision" {
		t.Fatalf("got %#v", got)
	}

	bad := repository.AuditEvent{AuditEventID: "audit_bad", Actor: "system", Action: "generate_decision", Status: "failed", CreatedAt: testTime}
	if err := repo.AppendAuditEvent(ctx, bad); err == nil {
		t.Fatal("expected failed event error_code constraint")
	}
	if _, err := repo.GetAuditEvent(ctx, "audit_bad"); err == nil {
		t.Fatal("bad audit event persisted")
	}
}
