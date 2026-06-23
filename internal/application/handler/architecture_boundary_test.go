package handler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandlersDoNotDependOnSQLDBOrQuerySQL(t *testing.T) {
	files := map[string]string{}
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read handler dir: %v", err)
	}
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		files[name] = readHandlerSource(t, name)
	}
	for name, source := range files {
		if strings.Contains(source, "database/sql") || strings.Contains(source, "*sql.DB") || strings.Contains(source, "db *sql.DB") {
			t.Fatalf("%s must not depend on database/sql", name)
		}
		for _, marker := range []string{".ExecContext(", ".QueryContext(", ".QueryRowContext("} {
			if strings.Contains(source, marker) {
				t.Fatalf("%s must not use SQL directly: %s", name, marker)
			}
		}
	}
}

func readHandlerSource(t *testing.T, name string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(".", name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(content)
}
