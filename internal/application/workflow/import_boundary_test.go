package workflow

import "testing"

func TestWorkflowPackageDoesNotImportSQLite(t *testing.T) {
	assertPackageDoesNotImport(t, "investment-agent/internal/application/workflow", "investment-agent/internal/infrastructure/persistence/sqlite")
	assertPackageDoesNotImport(t, "investment-agent/internal/application/workflow", "database/sql")
}

func TestDomainPackagesDoNotImportApplicationOrInfrastructure(t *testing.T) {
	for _, pkg := range []string{
		"investment-agent/internal/domain/model",
		"investment-agent/internal/domain/repository",
		"investment-agent/internal/domain/rule",
	} {
		assertPackageDoesNotImport(t, pkg, "investment-agent/internal/application")
		assertPackageDoesNotImport(t, pkg, "investment-agent/internal/infrastructure")
		assertPackageDoesNotImport(t, pkg, "database/sql")
		assertPackageDoesNotImport(t, pkg, "net/http")
	}
}

func TestLLMInfrastructureDoesNotImportApplication(t *testing.T) {
	assertPackageDoesNotImport(t, "investment-agent/internal/infrastructure/llm/deepseek", "investment-agent/internal/application")
}
