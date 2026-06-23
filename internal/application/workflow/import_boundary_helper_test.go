package workflow

import (
	"encoding/json"
	"os/exec"
	"strings"
	"testing"
)

func assertPackageDoesNotImport(t *testing.T, pkg string, forbidden string) {
	t.Helper()
	cmd := exec.Command("go", "list", "-json", pkg)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list %s: %v", pkg, err)
	}
	var data struct {
		Imports []string
	}
	if err := json.Unmarshal(out, &data); err != nil {
		t.Fatalf("parse go list output for %s: %v", pkg, err)
	}
	for _, imp := range data.Imports {
		if imp == forbidden || strings.HasPrefix(imp, forbidden+"/") {
			t.Fatalf("%s imports forbidden package %s", pkg, imp)
		}
	}
}
