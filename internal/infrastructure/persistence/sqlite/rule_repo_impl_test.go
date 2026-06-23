package sqlite

import (
	"context"
	"testing"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

func TestRuleRepositoryWriteReadAndRollback(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewRuleRepository(db)

	version := repository.RuleVersion{RuleVersion: "v2.0", Status: "archived", RulesJSON: "{}", EffectiveAt: testTime, CreatedAt: testTime}
	if err := repo.SaveRuleVersion(ctx, version); err != nil {
		t.Fatal(err)
	}
	gotVersion, err := repo.GetRuleVersion(ctx, "v2.0")
	if err != nil {
		t.Fatal(err)
	}
	if gotVersion.Status != "archived" {
		t.Fatalf("got %#v", gotVersion)
	}

	bad := repository.RuleVersion{RuleVersion: "bad", Status: "invalid", RulesJSON: "{}", EffectiveAt: testTime, CreatedAt: testTime}
	if err := repo.SaveRuleVersion(ctx, bad); err == nil {
		t.Fatal("expected constraint error")
	}
	if _, err := repo.GetRuleVersion(ctx, "bad"); err == nil {
		t.Fatal("bad rule version persisted")
	}
}

func TestRuleProposalAndGatekeeperAuditWriteRead(t *testing.T) {
	ctx := context.Background()
	db := testDB(t)
	repo := NewRuleRepository(db)

	proposal := repository.RuleProposal{
		ProposalID: "prop1", ProposalType: "threshold", Status: "draft", Title: "调整阈值",
		ProposalVersion: "v3.1-proposal", SampleCount: 1, CreatedAt: testTime,
	}
	if err := repo.SaveRuleProposal(ctx, proposal); err != nil {
		t.Fatal(err)
	}
	gotProposal, err := repo.GetRuleProposal(ctx, "prop1")
	if err != nil {
		t.Fatal(err)
	}
	if gotProposal.ProposalType != "threshold" {
		t.Fatalf("got %#v", gotProposal)
	}

	audit := repository.GatekeeperAudit{
		GatekeeperAuditID: "ga1", ProposalID: "prop1", AuditResult: "needs_user_review",
		ViolatesFundamentalRule: false, HasRuleConflict: false, AllowApply: false,
		AuditedRuleVersion: "v3.0", CreatedAt: testTime,
	}
	if err := repo.SaveGatekeeperAudit(ctx, audit); err != nil {
		t.Fatal(err)
	}
	gotAudit, err := repo.GetGatekeeperAudit(ctx, "ga1")
	if err != nil {
		t.Fatal(err)
	}
	if gotAudit.AuditResult != "needs_user_review" {
		t.Fatalf("got %#v", gotAudit)
	}
}

func TestRuleRepositoryClassifiesErrors(t *testing.T) {
	db := testDB(t)
	repo := NewRuleRepository(db)
	if _, err := repo.GetRuleVersion(context.Background(), "missing_rule"); !apperr.IsCode(err, apperr.CodeNotFound) {
		t.Fatalf("expected not found rule version error, got %v", err)
	}
	bad := repository.RuleVersion{RuleVersion: "bad_classified", Status: "invalid", RulesJSON: "{}", EffectiveAt: testTime, CreatedAt: testTime}
	if err := repo.SaveRuleVersion(context.Background(), bad); !apperr.IsCode(err, apperr.CodeConflict) {
		t.Fatalf("expected conflict rule version error, got %v", err)
	}
}
