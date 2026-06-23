package workflow

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWorkflowBusinessIDsUseGenerator(t *testing.T) {
	for _, name := range []string{"evidence_verification_graph.go", "market_refresh_graph.go", "evolution_proposal_graph.go", "gatekeeper_audit_graph.go"} {
		source := readWorkflowSource(t, name)
		for _, marker := range []string{"\"chunk_\" +", "\"group_\" +", "\"event_\" +", "\"market_\" +", "\"proposal_\" +", "\"gatekeeper_\" +"} {
			if strings.Contains(source, marker) {
				t.Fatalf("%s must use workflowID generator, found %s", name, marker)
			}
		}
	}
}

func readWorkflowSource(t *testing.T, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(".", name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(content)
}
