package idgen

import "testing"

func TestGeneratorCreatesReadableIDs(t *testing.T) {
	g := NewGenerator()
	id := g.New("audit")
	if len(id) <= len("audit_") || id[:len("audit_")] != "audit_" {
		t.Fatalf("expected audit_ prefix, got %q", id)
	}
}

func TestFixedGeneratorReturnsDeterministicIDs(t *testing.T) {
	g := NewFixedGenerator(map[string][]string{
		"audit": {"audit_one", "audit_two"},
	})
	if got := g.New("audit"); got != "audit_one" {
		t.Fatalf("first id = %q", got)
	}
	if got := g.New("audit"); got != "audit_two" {
		t.Fatalf("second id = %q", got)
	}
}
