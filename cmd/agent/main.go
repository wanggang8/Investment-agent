package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"investment-agent/internal/application/dto"
	"investment-agent/internal/application/service"
	"investment-agent/internal/application/workflow"
	"investment-agent/internal/domain/model"
	"investment-agent/internal/domain/repository"
	"investment-agent/internal/infrastructure/config"
	appsqlite "investment-agent/internal/infrastructure/persistence/sqlite"
	"investment-agent/internal/infrastructure/wiring"
	"investment-agent/internal/pkg/apperr"
	"investment-agent/internal/pkg/clock"
	"investment-agent/internal/pkg/idgen"
	"investment-agent/pkg/logger"
)

var supportedTasks = map[string]model.AuditAction{
	"daily":                                model.AuditActionGenerateDecision,
	"market-refresh":                       model.AuditActionRefreshMarketData,
	"evidence-index":                       model.AuditActionRebuildIndex,
	"public-evidence-refresh":              model.AuditActionRunLocalTask,
	"p34-expanded-refresh":                 model.AuditActionRunLocalTask,
	"llm-smoke":                            model.AuditActionRunLocalTask,
	"retrieval-quality-smoke":              model.AuditActionRunLocalTask,
	"data-source-quality-regression":       model.AuditActionRunLocalTask,
	"data-source-quality-resolution-check": model.AuditActionRunLocalTask,
	"review":                               model.AuditActionRunLocalTask,
}

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("agent", flag.ContinueOnError)
	fs.SetOutput(stderr)
	configPath := fs.String("config", "", "配置文件路径，默认读取 INVESTMENT_AGENT_CONFIG 或 configs/config.example.yaml")
	task := fs.String("task", "", "手动任务：daily、market-refresh、evidence-index、public-evidence-refresh、p34-expanded-refresh、llm-smoke、retrieval-quality-smoke、data-source-quality-regression、data-source-quality-resolution-check、review")
	period := fs.String("period", "monthly", "复盘周期：monthly 或 quarterly")
	source := fs.String("source", "", "P34 扩展数据源或 P48 回归模式：sentiment_proxy_fixture、configured、fixture、current")
	symbol := fs.String("symbol", "510300", "本地任务标的代码，用于 market-refresh、public-evidence-refresh 与 p34-expanded-refresh")
	strictQualityGate := fs.Bool("strict-quality-gate", false, "data-source-quality-regression current 模式严格发布门禁；policy gate=block 时返回失败")
	startDateText := fs.String("start-date", "", "public-evidence-refresh / p34-expanded-refresh 起始日期，格式 YYYY-MM-DD")
	endDateText := fs.String("end-date", "", "public-evidence-refresh / p34-expanded-refresh 结束日期，格式 YYYY-MM-DD")
	showSchedule := fs.Bool("schedule", false, "显示本地调度配置说明；默认不会自动执行任务")
	validateConfig := fs.Bool("validate-config", false, "校验本地配置并输出诊断，不执行任务")
	preflight := fs.Bool("preflight", false, "执行 P40 本地部署预检，检查依赖、路径和配置，不执行任务")
	diagnosticsPath := fs.String("diagnostics", "", "将预检诊断写入本地 JSON 文件，不包含密钥原文")
	releaseUpgradeCheck := fs.Bool("release-upgrade-check", false, "执行 P49 本地发布/升级检查，输出版本、备份、迁移和 smoke 计划，不执行升级")
	targetVersion := fs.String("target-version", "", "P49 发布/升级检查的目标版本或 release label")
	backupDir := fs.String("backup", "", "将 SQLite 数据库备份到指定目录，不执行任务")
	restorePath := fs.String("restore", "", "从指定备份文件恢复 SQLite 数据库，默认拒绝覆盖现有数据库")
	restoreConfirm := fs.Bool("restore-confirm", false, "确认恢复操作允许写入 sqlite.path")
	recoverySmokePath := fs.String("recovery-smoke", "", "将指定备份恢复到 sqlite.path 并验证本地事实可读，目标库存在时拒绝覆盖")
	help := fs.Bool("help", false, "显示帮助")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *help || len(args) == 0 {
		printHelp(stdout)
		return 0
	}
	if *showSchedule {
		// 调度配置只做显式说明，避免用户误以为本工具会在后台自动运行或执行交易。
		fmt.Fprintln(stdout, "本地调度配置：默认不自动运行；需要用户显式安装或编辑系统计划任务；不会执行交易；不会自动应用规则。详见 docs/ops-local-scheduler.md。")
		return 0
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(stderr, "load config: %v\n", err)
		return 1
	}
	if *preflight {
		report := buildPreflightReport(*configPath, cfg)
		if strings.TrimSpace(*diagnosticsPath) != "" {
			if err := writePreflightDiagnostics(*diagnosticsPath, report); err != nil {
				fmt.Fprintf(stderr, "write diagnostics: %v\n", err)
				return 1
			}
		}
		printPreflightReport(stdout, report)
		if preflightHasFailed(report) {
			return 1
		}
		return 0
	}
	if *releaseUpgradeCheck {
		report := buildReleaseUpgradeReport(cfg, *targetVersion)
		if strings.TrimSpace(*diagnosticsPath) != "" {
			if err := writeReleaseUpgradeDiagnostics(*diagnosticsPath, report); err != nil {
				fmt.Fprintf(stderr, "write diagnostics: %v\n", err)
				return 1
			}
		}
		printReleaseUpgradeReport(stdout, report)
		if releaseUpgradeHasFailed(report) {
			return 1
		}
		return 0
	}
	if err := cfg.Validate(); err != nil {
		fmt.Fprintf(stderr, "validate config: %v\n", err)
		return 1
	}
	if *validateConfig {
		fmt.Fprintf(stdout, "config validation passed: sqlite=%s veclite=%s；不会执行交易。\n", cfg.SQLite.Path, cfg.VecLite.Path)
		return 0
	}
	if strings.TrimSpace(*backupDir) != "" {
		backupPath, err := backupSQLite(cfg.SQLite.Path, *backupDir)
		if err != nil {
			fmt.Fprintf(stderr, "backup sqlite: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "backup created:%s\n", backupPath)
		return 0
	}
	if strings.TrimSpace(*restorePath) != "" {
		if err := restoreSQLite(*restorePath, cfg.SQLite.Path, *restoreConfirm); err != nil {
			fmt.Fprintf(stderr, "restore sqlite: %v\n", err)
			return 1
		}
		fmt.Fprintln(stdout, "restore completed；不会执行交易。")
		return 0
	}
	if strings.TrimSpace(*recoverySmokePath) != "" {
		summary, err := runRecoverySmoke(ctx, cfg, *recoverySmokePath)
		if err != nil {
			fmt.Fprintf(stderr, "recovery smoke: %v\n", err)
			return 1
		}
		fmt.Fprintf(stdout, "recovery smoke completed:%s；不会执行交易；不会外部推送；不会自动应用规则。\n", summary)
		return 0
	}
	if strings.TrimSpace(*task) == "" {
		fmt.Fprintln(stderr, "missing --task；可选 daily、market-refresh、evidence-index、public-evidence-refresh、p34-expanded-refresh、llm-smoke、retrieval-quality-smoke、data-source-quality-regression、data-source-quality-resolution-check、review")
		return 2
	}
	if _, ok := supportedTasks[*task]; !ok {
		fmt.Fprintf(stderr, "unsupported task %q；可选 daily、market-refresh、evidence-index、public-evidence-refresh、p34-expanded-refresh、llm-smoke、retrieval-quality-smoke、data-source-quality-regression、data-source-quality-resolution-check、review\n", *task)
		return 2
	}
	if *task == "review" && *period != "monthly" && *period != "quarterly" {
		fmt.Fprintln(stderr, "review --period 只支持 monthly 或 quarterly")
		return 2
	}
	var startDate, endDate time.Time
	if *task == "public-evidence-refresh" || *task == "p34-expanded-refresh" {
		startDate, endDate, err = parsePublicEvidenceDateWindow(*startDateText, *endDateText)
		if err != nil {
			fmt.Fprintf(stderr, "%s date window: %v\n", *task, err)
			return 2
		}
	} else if strings.TrimSpace(*startDateText) != "" || strings.TrimSpace(*endDateText) != "" {
		fmt.Fprintf(stdout, "%s 任务会忽略 --start-date/--end-date；该参数只用于 public-evidence-refresh 与 p34-expanded-refresh。\n", *task)
	}

	if *task != "review" && *period != "monthly" {
		fmt.Fprintf(stdout, "%s 任务会忽略 --period；该参数只用于 review。\n", *task)
	}

	log := logger.New(cfg.Log.Level)
	log.Info("starting local agent task", "task", *task, "sqlite", cfg.SQLite.Path)

	outputRef, err := runTask(ctx, cfg, *task, *period, *source, *symbol, *strictQualityGate, startDate, endDate)
	if err != nil {
		fmt.Fprintf(stderr, "run task %s: %v\n", *task, err)
		return 1
	}
	if *task == "data-source-quality-regression" {
		fmt.Fprintf(stdout, "data source quality regression completed:%s；不会执行交易。\n", outputRef)
		return 0
	}
	if *task == "data-source-quality-resolution-check" {
		fmt.Fprintf(stdout, "data source quality resolution check completed:%s；不会执行交易。\n", outputRef)
		return 0
	}
	fmt.Fprintf(stdout, "task %s completed；已写入 audit_events；不会执行交易。\n", *task)
	return 0
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "Investment Agent 本地任务入口")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "用法：")
	fmt.Fprintln(w, "  go run ./cmd/agent --task daily")
	fmt.Fprintln(w, "  go run ./cmd/agent --task market-refresh")
	fmt.Fprintln(w, "  go run ./cmd/agent --task evidence-index")
	fmt.Fprintln(w, "  go run ./cmd/agent --task public-evidence-refresh --symbol 510300 --start-date YYYY-MM-DD --end-date YYYY-MM-DD")
	fmt.Fprintln(w, "  go run ./cmd/agent --task p34-expanded-refresh --source sentiment_proxy_fixture --symbol 000300 --start-date YYYY-MM-DD --end-date YYYY-MM-DD")
	fmt.Fprintln(w, "  go run ./cmd/agent --task llm-smoke --symbol 510300")
	fmt.Fprintln(w, "  go run ./cmd/agent --task retrieval-quality-smoke --symbol 510300")
	fmt.Fprintln(w, "  go run ./cmd/agent --task data-source-quality-regression --source fixture|current --symbol 000300")
	fmt.Fprintln(w, "  go run ./cmd/agent --task data-source-quality-regression --source current --symbol 000300 --strict-quality-gate")
	fmt.Fprintln(w, "  go run ./cmd/agent --task data-source-quality-resolution-check --symbol 000300")
	fmt.Fprintln(w, "  go run ./cmd/agent --task review --period monthly|quarterly")
	fmt.Fprintln(w, "  go run ./cmd/agent --validate-config")
	fmt.Fprintln(w, "  go run ./cmd/agent --schedule")
	fmt.Fprintln(w, "  go run ./cmd/agent --backup ./data/backups")
	fmt.Fprintln(w, "  go run ./cmd/agent --restore ./data/backups/agent-YYYYMMDDTHHMMSSZ.db --restore-confirm")
	fmt.Fprintln(w, "  go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json")
	fmt.Fprintln(w, "  go run ./cmd/agent --release-upgrade-check --target-version vNEXT --diagnostics ./tmp/release-upgrade.json")
	fmt.Fprintln(w, "  go run ./cmd/agent --recovery-smoke ./data/backups/agent-YYYYMMDDTHHMMSSZ.db")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "说明：本入口只触发本地分析、刷新、索引和复盘任务，不会执行交易，不会自动应用规则。")
	fmt.Fprintln(w, "本地调度：仅提供 launchd/cron 示例，默认不自动运行，需要用户显式安装；任务结果写入 audit_events。")
}

type preflightCheck struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Detail      string `json:"detail,omitempty"`
	Remediation string `json:"remediation,omitempty"`
}

type preflightReport struct {
	GeneratedAt string           `json:"generated_at"`
	ConfigPath  string           `json:"config_path"`
	SafetyNote  string           `json:"safety_note"`
	Checks      []preflightCheck `json:"checks"`
}

func buildPreflightReport(configPath string, cfg *config.Config) preflightReport {
	report := preflightReport{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		ConfigPath:  resolveConfigPath(configPath),
		SafetyNote:  "本地预检只读取配置和依赖状态；不会执行交易；不会外部推送；不会自动应用规则。",
	}
	report.Checks = append(report.Checks,
		preflightConfigCheck(cfg),
		executableCheck("go", "go_binary"),
		executableCheck("node", "node_binary"),
		executableCheck("npm", "npm_binary"),
		playwrightCheck(),
		readOnlyPathCheck("sqlite_path", cfg.SQLite.Path, false),
		readOnlyPathCheck("veclite_path", cfg.VecLite.Path, true),
		dataSourceCheck(cfg),
		deepSeekCheck(cfg),
	)
	return report
}

func resolveConfigPath(configPath string) string {
	if strings.TrimSpace(configPath) != "" {
		return configPath
	}
	if env := os.Getenv("INVESTMENT_AGENT_CONFIG"); strings.TrimSpace(env) != "" {
		return env
	}
	return "configs/config.example.yaml"
}

func preflightConfigCheck(cfg *config.Config) preflightCheck {
	if err := cfg.Validate(); err != nil {
		return preflightCheck{Name: "config_validation", Status: "failed", Detail: err.Error(), Remediation: "修复配置文件中的必填项或 URL，再重新运行 --preflight。"}
	}
	return preflightCheck{Name: "config_validation", Status: "pass", Detail: "配置结构校验通过"}
}

func executableCheck(binary string, name string) preflightCheck {
	path, err := exec.LookPath(binary)
	if err != nil {
		return preflightCheck{Name: name, Status: "failed", Detail: binary + " not found", Remediation: "安装 " + binary + " 并确认 PATH 可见。"}
	}
	return preflightCheck{Name: name, Status: "pass", Detail: path}
}

func playwrightCheck() preflightCheck {
	cliPath, ok := playwrightCLIPath()
	return playwrightCheckWithPaths(cliPath, ok, playwrightBrowserSearchDirs())
}

func playwrightCheckWithPaths(cliPath string, cliOK bool, browserDirs []string) preflightCheck {
	if !cliOK {
		return preflightCheck{Name: "playwright_browser", Status: "skipped", Detail: "Playwright package not installed; browser check skipped", Remediation: "运行 npm --prefix web install && npm --prefix web exec playwright install chromium。"}
	}
	browserPath, ok := playwrightChromiumExecutable(browserDirs)
	if !ok {
		return preflightCheck{Name: "playwright_browser", Status: "warning", Detail: "Playwright CLI found: " + cliPath + "; chromium browser not found", Remediation: "运行 npm --prefix web exec playwright install chromium。"}
	}
	return preflightCheck{Name: "playwright_browser", Status: "pass", Detail: "cli=" + cliPath + "; browser=" + browserPath}
}

func playwrightCLIPath() (string, bool) {
	root := repoRoot()
	candidates := []string{
		filepath.Join(root, "web", "node_modules", ".bin", "playwright"),
		filepath.Join(root, "web", "node_modules", ".bin", "playwright.cmd"),
		filepath.Join(root, "web", "node_modules", "playwright"),
	}
	for _, candidate := range candidates {
		if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
			return candidate, true
		}
	}
	return "", false
}

func playwrightChromiumExecutable(browserDirs []string) (string, bool) {
	executableCandidates := []string{
		filepath.Join("chrome-mac", "Chromium.app", "Contents", "MacOS", "Chromium"),
		filepath.Join("chrome-linux", "chrome"),
		filepath.Join("chrome-win", "chrome.exe"),
		filepath.Join("chromium", "chrome-win", "chrome.exe"),
		filepath.Join("chromium", "chrome"),
	}
	for _, base := range browserDirs {
		if base == "" {
			continue
		}
		entries, err := os.ReadDir(base)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() && strings.HasPrefix(entry.Name(), "chromium") {
				browserRoot := filepath.Join(base, entry.Name())
				for _, candidate := range executableCandidates {
					executable := filepath.Join(browserRoot, candidate)
					if isExecutableFile(executable) {
						return executable, true
					}
				}
			}
		}
	}
	return "", false
}

func isExecutableFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return false
	}
	if strings.HasSuffix(strings.ToLower(path), ".exe") {
		return true
	}
	return info.Mode().Perm()&0o111 != 0
}

func playwrightBrowserSearchDirs() []string {
	root := repoRoot()
	dirs := []string{}
	if custom := strings.TrimSpace(os.Getenv("PLAYWRIGHT_BROWSERS_PATH")); custom != "" {
		if custom == "0" {
			dirs = append(dirs,
				filepath.Join(root, "web", "node_modules", "playwright-core", ".local-browsers"),
				filepath.Join(root, "web", "node_modules", "playwright", ".local-browsers"),
			)
		} else {
			dirs = append(dirs, custom)
		}
	}
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		dirs = append(dirs,
			filepath.Join(home, "Library", "Caches", "ms-playwright"),
			filepath.Join(home, ".cache", "ms-playwright"),
		)
	}
	if localAppData := strings.TrimSpace(os.Getenv("LOCALAPPDATA")); localAppData != "" {
		dirs = append(dirs, filepath.Join(localAppData, "ms-playwright"))
	}
	return dirs
}

func repoRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		if fileExists(filepath.Join(dir, "go.mod")) && fileExists(filepath.Join(dir, "web", "package.json")) {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return wd
		}
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func readOnlyPathCheck(name string, path string, pathMayBeDirectory bool) preflightCheck {
	path = strings.TrimSpace(path)
	if path == "" {
		return preflightCheck{Name: name, Status: "failed", Detail: "path is empty", Remediation: "在配置文件中填写本地路径。"}
	}
	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			if pathMayBeDirectory {
				return readOnlyDirPermissionCheck(name, path, "directory exists")
			}
			return preflightCheck{Name: name, Status: "failed", Detail: "path is directory", Remediation: "配置为具体 SQLite 文件路径。"}
		}
		return readOnlyDirPermissionCheck(name, filepath.Dir(path), "file exists")
	} else if !os.IsNotExist(err) {
		return preflightCheck{Name: name, Status: "failed", Detail: err.Error(), Remediation: "检查路径权限或文件状态。"}
	}
	parent := filepath.Dir(path)
	if info, err := os.Stat(parent); err == nil {
		if !info.IsDir() {
			return preflightCheck{Name: name, Status: "failed", Detail: "parent is not directory: " + parent, Remediation: "调整配置到有效本地目录。"}
		}
		check := readOnlyDirPermissionCheck(name, parent, "path missing; parent exists")
		if check.Status == "pass" {
			check.Status = "warning"
			check.Remediation = "首次启动或迁移会创建本地文件；如非预期，请检查配置路径。"
		}
		return check
	} else if !os.IsNotExist(err) {
		return preflightCheck{Name: name, Status: "failed", Detail: err.Error(), Remediation: "检查上级目录权限或文件状态。"}
	}
	if pathMayBeDirectory {
		return preflightCheck{Name: name, Status: "warning", Detail: "path and parent missing: " + path, Remediation: "确认 VecLite 路径或先创建目标目录；预检不会代为创建。"}
	}
	return preflightCheck{Name: name, Status: "warning", Detail: "path parent missing: " + parent, Remediation: "确认 SQLite 路径或先创建目标目录；预检不会代为创建。"}
}

func readOnlyDirPermissionCheck(name string, dir string, detailPrefix string) preflightCheck {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return preflightCheck{Name: name, Status: "warning", Detail: detailPrefix + "; directory missing: " + dir, Remediation: "确认配置路径；预检不会创建目录。"}
		}
		return preflightCheck{Name: name, Status: "failed", Detail: err.Error(), Remediation: "检查目录权限或文件状态。"}
	}
	if !info.IsDir() {
		return preflightCheck{Name: name, Status: "failed", Detail: detailPrefix + "; not a directory: " + dir, Remediation: "调整配置到有效本地目录。"}
	}
	if info.Mode().Perm()&0o222 == 0 {
		return preflightCheck{Name: name, Status: "failed", Detail: detailPrefix + "; directory not writable by mode: " + dir, Remediation: "修复目录写权限后重试。"}
	}
	return preflightCheck{Name: name, Status: "pass", Detail: detailPrefix + "; directory mode allows write: " + dir}
}

func dataSourceCheck(cfg *config.Config) preflightCheck {
	if cfg.DataSources.UseStub {
		return preflightCheck{Name: "data_sources", Status: "pass", Detail: "stub enabled"}
	}
	if cfg.DataSources.PublicEvidence.Enabled || cfg.DataSources.MarketCollectors.Enabled || len(cfg.DataSources.Enabled) > 0 {
		return preflightCheck{Name: "data_sources", Status: "pass", Detail: "configured read-only sources"}
	}
	return preflightCheck{Name: "data_sources", Status: "warning", Detail: "no explicit source enabled", Remediation: "启用 stub 或配置只读公开数据源。"}
}

func deepSeekCheck(cfg *config.Config) preflightCheck {
	if strings.TrimSpace(cfg.DeepSeek.APIKey) == "" {
		return preflightCheck{Name: "deepseek", Status: "warning", Detail: "api key missing", Remediation: "需要真实 LLM smoke 时在本地配置中填写 key；预检不会输出密钥原文。"}
	}
	return preflightCheck{Name: "deepseek", Status: "pass", Detail: "api key configured; model=" + strings.TrimSpace(cfg.DeepSeek.Model)}
}

func writePreflightDiagnostics(path string, report preflightReport) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func printPreflightReport(w io.Writer, report preflightReport) {
	fmt.Fprintf(w, "preflight generated_at=%s config=%s\n", report.GeneratedAt, report.ConfigPath)
	for _, check := range report.Checks {
		if check.Remediation != "" {
			fmt.Fprintf(w, "%s:%s:%s；修复：%s\n", check.Name, check.Status, check.Detail, check.Remediation)
		} else {
			fmt.Fprintf(w, "%s:%s:%s\n", check.Name, check.Status, check.Detail)
		}
	}
	fmt.Fprintln(w, report.SafetyNote)
}

func preflightHasFailed(report preflightReport) bool {
	for _, check := range report.Checks {
		if check.Status == "failed" {
			return true
		}
	}
	return false
}

type releaseUpgradeCheckItem struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Detail      string `json:"detail,omitempty"`
	Remediation string `json:"remediation,omitempty"`
}

type releaseUpgradeReport struct {
	GeneratedAt               string                    `json:"generated_at"`
	CurrentVersion            string                    `json:"current_version"`
	TargetVersion             string                    `json:"target_version,omitempty"`
	Status                    string                    `json:"status"`
	Checks                    []releaseUpgradeCheckItem `json:"checks"`
	BackupReminder            string                    `json:"backup_reminder"`
	PreUpgradeCommands        []string                  `json:"pre_upgrade_commands"`
	PostUpgradeSmokeCommands  []string                  `json:"post_upgrade_smoke_commands"`
	GeneratedArtifactBoundary string                    `json:"generated_artifact_boundary"`
	SafetyNote                string                    `json:"safety_note"`
}

func buildReleaseUpgradeReport(cfg *config.Config, targetVersion string) releaseUpgradeReport {
	sanitizedTargetVersion, targetVersionSafe := sanitizeReleaseTargetVersion(targetVersion)
	checks := []releaseUpgradeCheckItem{
		releaseVersionCheck(targetVersion, targetVersionSafe),
		releaseConfigCheck(cfg),
		releaseSQLiteBackupCheck(cfg.SQLite.Path),
		releaseMigrationPrecheck(),
		releaseSmokePlanCheck(),
	}
	report := releaseUpgradeReport{
		GeneratedAt:    time.Now().UTC().Format(time.RFC3339),
		CurrentVersion: releaseCurrentVersion(),
		TargetVersion:  sanitizedTargetVersion,
		Checks:         checks,
		BackupReminder: "升级前请手动运行：go run ./cmd/agent --backup ./data/backups；本检查不会创建备份。",
		PreUpgradeCommands: []string{
			"go run ./cmd/agent --preflight --diagnostics ./tmp/preflight.json",
			"go run ./cmd/agent --backup ./data/backups",
		},
		PostUpgradeSmokeCommands: []string{
			"go run ./cmd/agent --preflight --diagnostics ./tmp/preflight-after-upgrade.json",
			"bash scripts/recovery-smoke.sh",
			"bash scripts/e2e-smoke.sh",
			"bash scripts/local-install-diagnostics.sh --skip-e2e",
		},
		GeneratedArtifactBoundary: "诊断建议写入 tmp/；不要提交本地数据库、备份、日志或私密配置。",
		SafetyNote:                "P49 检查只读取本地配置和迁移文件状态；不会执行升级；不会运行迁移；不会创建备份；不会恢复或覆盖数据库；不会执行交易；不会外部推送；不会自动确认；不会自动应用规则；不会自动修复。",
	}
	report.Status = releaseUpgradeStatus(checks)
	return report
}

func releaseCurrentVersion() string {
	return "local-dev"
}

func sanitizeReleaseTargetVersion(targetVersion string) (string, bool) {
	value := strings.TrimSpace(targetVersion)
	if value == "" {
		return "", true
	}
	if len(value) > 64 {
		return "<redacted-target-version>", false
	}
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '_' || r == '-' || r == '+' {
			continue
		}
		return "<redacted-target-version>", false
	}
	return value, true
}

func releaseVersionCheck(targetVersion string, targetVersionSafe bool) releaseUpgradeCheckItem {
	if strings.TrimSpace(targetVersion) == "" {
		return releaseUpgradeCheckItem{Name: "version_check", Status: "warning", Detail: "target version missing", Remediation: "使用 --target-version 指定目标版本或 release label。"}
	}
	if !targetVersionSafe {
		return releaseUpgradeCheckItem{Name: "version_check", Status: "warning", Detail: "target version redacted due to unsafe characters", Remediation: "使用只包含字母、数字、点、下划线、加号或短横线的 release label。"}
	}
	return releaseUpgradeCheckItem{Name: "version_check", Status: "pass", Detail: "target version provided"}
}

func releaseConfigCheck(cfg *config.Config) releaseUpgradeCheckItem {
	if err := cfg.Validate(); err != nil {
		return releaseUpgradeCheckItem{Name: "config_validation", Status: "failed", Detail: err.Error(), Remediation: "先修复配置结构，再执行升级。"}
	}
	return releaseUpgradeCheckItem{Name: "config_validation", Status: "pass", Detail: "配置结构校验通过"}
}

func releaseSQLiteBackupCheck(dbPath string) releaseUpgradeCheckItem {
	if strings.TrimSpace(dbPath) == "" {
		return releaseUpgradeCheckItem{Name: "backup_reminder", Status: "failed", Detail: "sqlite.path is empty", Remediation: "先在配置中填写 sqlite.path，再执行升级前备份。"}
	}
	info, err := os.Stat(dbPath)
	if err == nil {
		if info.IsDir() {
			return releaseUpgradeCheckItem{Name: "backup_reminder", Status: "failed", Detail: "sqlite path is directory", Remediation: "sqlite.path 必须指向具体 SQLite 文件。"}
		}
		return releaseUpgradeCheckItem{Name: "backup_reminder", Status: "warning", Detail: "sqlite database exists; backup required before upgrade", Remediation: "手动运行 go run ./cmd/agent --backup ./data/backups。"}
	}
	if os.IsNotExist(err) {
		return releaseUpgradeCheckItem{Name: "backup_reminder", Status: "warning", Detail: "sqlite database missing; first install or path needs confirmation", Remediation: "确认配置路径；如果是旧库升级，先指向现有 SQLite 并备份。"}
	}
	return releaseUpgradeCheckItem{Name: "backup_reminder", Status: "failed", Detail: "sqlite path stat failed", Remediation: "检查 SQLite 路径权限或文件状态。"}
}

func releaseMigrationPrecheck() releaseUpgradeCheckItem {
	migrationDir := filepath.Join(repoRoot(), "internal", "infrastructure", "persistence", "sqlite", "migration")
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		return releaseUpgradeCheckItem{Name: "migration_precheck", Status: "failed", Detail: "migration directory unreadable", Remediation: "确认仓库包含 sqlite migration 文件后再升级。"}
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	if len(names) == 0 {
		return releaseUpgradeCheckItem{Name: "migration_precheck", Status: "failed", Detail: "no migration files found", Remediation: "确认 sqlite migration 文件已随版本发布。"}
	}
	return releaseUpgradeCheckItem{Name: "migration_precheck", Status: "pass", Detail: fmt.Sprintf("migration_files=%d last=%s", len(names), names[len(names)-1])}
}

func releaseSmokePlanCheck() releaseUpgradeCheckItem {
	return releaseUpgradeCheckItem{Name: "smoke_plan", Status: "pass", Detail: "post-upgrade smoke commands listed for manual execution"}
}

func releaseUpgradeStatus(checks []releaseUpgradeCheckItem) string {
	status := "ready"
	for _, check := range checks {
		if check.Status == "failed" {
			return "blocked"
		}
		if check.Status == "warning" {
			status = "warning"
		}
	}
	return status
}

func writeReleaseUpgradeDiagnostics(path string, report releaseUpgradeReport) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func printReleaseUpgradeReport(w io.Writer, report releaseUpgradeReport) {
	fmt.Fprintf(w, "release upgrade generated_at=%s current=%s target=%s status=%s\n", report.GeneratedAt, report.CurrentVersion, missingValueText(report.TargetVersion), report.Status)
	for _, check := range report.Checks {
		if check.Remediation != "" {
			fmt.Fprintf(w, "%s:%s:%s；处理：%s\n", check.Name, check.Status, check.Detail, check.Remediation)
		} else {
			fmt.Fprintf(w, "%s:%s:%s\n", check.Name, check.Status, check.Detail)
		}
	}
	fmt.Fprintln(w, "pre-upgrade commands:")
	for _, command := range report.PreUpgradeCommands {
		fmt.Fprintf(w, "- %s\n", command)
	}
	fmt.Fprintln(w, "post-upgrade smoke commands:")
	for _, command := range report.PostUpgradeSmokeCommands {
		fmt.Fprintf(w, "- %s\n", command)
	}
	fmt.Fprintln(w, report.BackupReminder)
	fmt.Fprintln(w, report.SafetyNote)
}

func missingValueText(value string) string {
	if strings.TrimSpace(value) == "" {
		return "<missing>"
	}
	return value
}

func releaseUpgradeHasFailed(report releaseUpgradeReport) bool {
	return report.Status == "blocked"
}

type agentRuntime struct {
	store      *appsqlite.Store
	repos      repository.Repositories
	transactor repository.Transactor
	deps       workflow.WorkflowDependencies
}

func runTask(ctx context.Context, cfg *config.Config, task string, period string, source string, symbol string, strictQualityGate bool, startDate, endDate time.Time) (string, error) {
	rt, err := openRuntime(ctx, cfg)
	if err != nil {
		return "", err
	}
	defer rt.store.Close()

	outputRef := "no_auto_trading"
	switch task {
	case "retrieval-quality-smoke":
		outputRef, err = runRetrievalQualitySmoke(ctx, rt, symbol)
	case "data-source-quality-regression":
		outputRef, err = runDataSourceQualityRegression(ctx, rt, source, symbol, strictQualityGate)
	case "data-source-quality-resolution-check":
		outputRef, err = runDataSourceQualityResolutionCheck(ctx, rt, symbol)
	default:
		err = executeTask(ctx, cfg, rt, task, period, source, symbol, startDate, endDate)
	}
	if err != nil {
		auditOutputRef := "task_failed"
		if task == "data-source-quality-regression" && strictQualityGate && strings.HasPrefix(strings.TrimSpace(outputRef), "data_source_quality:") {
			auditOutputRef = outputRef
		}
		if task == "data-source-quality-resolution-check" && strings.HasPrefix(strings.TrimSpace(outputRef), "data_quality_gate_resolution:") {
			auditOutputRef = outputRef
		}
		if auditErr := appendTaskAudit(ctx, rt.repos.AuditRepo, task, period, source, symbol, cfg.DeepSeek.Model, startDate, endDate, string(model.AuditStatusFailed), errorCode(err), auditOutputRef); auditErr != nil {
			return "", fmt.Errorf("%w; write failed audit: %v", err, auditErr)
		}
		return "", err
	}
	return outputRef, appendTaskAudit(ctx, rt.repos.AuditRepo, task, period, source, symbol, cfg.DeepSeek.Model, startDate, endDate, string(model.AuditStatusSuccess), "", outputRef)
}

func openRuntime(ctx context.Context, cfg *config.Config) (agentRuntime, error) {
	store, err := appsqlite.Open(cfg.SQLite.Path)
	if err != nil {
		return agentRuntime{}, fmt.Errorf("open sqlite: %w", err)
	}
	if err := appsqlite.Migrate(ctx, store.DB); err != nil {
		_ = store.Close()
		return agentRuntime{}, fmt.Errorf("migrate sqlite: %w", err)
	}
	transactor := appsqlite.NewTransactor(store.DB)
	repos := repository.Repositories{
		DecisionRepo:                  appsqlite.NewDecisionRepository(store.DB),
		AuditRepo:                     appsqlite.NewAuditRepository(store.DB),
		RuleRepo:                      appsqlite.NewRuleRepository(store.DB),
		MarketRepo:                    appsqlite.NewMarketRepository(store.DB),
		SettingsRepo:                  appsqlite.NewSettingsRepository(store.DB),
		IntelligenceRepo:              appsqlite.NewIntelligenceRepository(store.DB),
		PortfolioRepo:                 appsqlite.NewPortfolioRepository(store.DB),
		DailyAutoRunRepo:              appsqlite.NewDailyAutoRunRepository(store.DB),
		DailyDisciplineReportRepo:     appsqlite.NewDailyDisciplineReportRepository(store.DB),
		RiskAlertRepo:                 appsqlite.NewRiskAlertRepository(store.DB),
		DataQualityGateResolutionRepo: appsqlite.NewDataQualityGateResolutionRepository(store.DB),
	}
	return agentRuntime{store: store, repos: repos, transactor: transactor, deps: wiring.NewWorkflowDependencies(cfg, repos, transactor)}, nil
}

func backupSQLite(dbPath, backupDir string) (string, error) {
	if strings.TrimSpace(dbPath) == "" {
		return "", fmt.Errorf("sqlite.path is required")
	}
	info, err := os.Stat(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("sqlite.path does not exist: %s", dbPath)
		}
		return "", err
	}
	if info.IsDir() {
		return "", fmt.Errorf("sqlite.path is not a file: %s", dbPath)
	}
	if err := os.MkdirAll(backupDir, 0o700); err != nil {
		return "", err
	}
	name := "agent-" + time.Now().UTC().Format("20060102T150405Z") + ".db"
	backupPath := filepath.Join(backupDir, name)
	if err := backupSQLiteConsistent(dbPath, backupPath); err != nil {
		return "", err
	}
	return backupPath, nil
}

func backupSQLiteConsistent(srcPath, dstPath string) error {
	db, err := sql.Open("sqlite", srcPath)
	if err != nil {
		return err
	}
	defer db.Close()
	quoted, err := sqliteStringLiteral(filepath.Clean(dstPath))
	if err != nil {
		return err
	}
	_, err = db.Exec("VACUUM INTO " + quoted)
	return err
}

func sqliteStringLiteral(value string) (string, error) {
	if strings.ContainsRune(value, 0) {
		return "", fmt.Errorf("path contains NUL byte")
	}
	return "'" + strings.ReplaceAll(value, "'", "''") + "'", nil
}

func restoreSQLite(backupPath, dbPath string, confirmed bool) error {
	if !confirmed {
		return fmt.Errorf("refuse to restore without --restore-confirm")
	}
	if strings.TrimSpace(backupPath) == "" || strings.TrimSpace(dbPath) == "" {
		return fmt.Errorf("restore source and sqlite.path are required")
	}
	if _, err := os.Lstat(dbPath); err == nil {
		return fmt.Errorf("refuse to overwrite existing sqlite; move it aside before restore")
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o700); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(dbPath), filepath.Base(dbPath)+".restore-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	if err := copyFile(backupPath, tmpPath, 0o600); err != nil {
		_ = os.Remove(tmpPath)
		return err
	}
	return os.Rename(tmpPath, dbPath)
}

func runRecoverySmoke(ctx context.Context, cfg *config.Config, backupPath string) (string, error) {
	if err := restoreSQLite(backupPath, cfg.SQLite.Path, true); err != nil {
		return "", err
	}
	rt, err := openRuntime(ctx, cfg)
	if err != nil {
		return "", err
	}
	defer rt.store.Close()
	counts, err := recoverySmokeCounts(rt.store.DB)
	if err != nil {
		return "", err
	}
	if counts.totalRestoredFacts() == 0 {
		return "", fmt.Errorf("restored database contains no readable local facts")
	}
	outputRef := fmt.Sprintf("recovery_smoke:decisions=%d:audits=%d:portfolio=%d:reports=%d:intelligence=%d:no_auto_trading", counts.Decisions, counts.Audits, counts.Portfolios, counts.Reports, counts.Intelligence)
	if err := rt.repos.AuditRepo.AppendAuditEvent(ctx, repository.AuditEvent{
		AuditEventID:  idgen.NewGenerator().New("audit"),
		RequestID:     idgen.NewGenerator().New("req"),
		WorkflowType:  "recovery_smoke",
		NodeName:      "cmd_agent",
		Actor:         string(model.AuditActorUser),
		Action:        string(model.AuditActionRunLocalTask),
		NodeAction:    "recovery_smoke",
		Status:        string(model.AuditStatusSuccess),
		InputRefType:  "backup_file",
		InputRef:      filepath.Base(backupPath),
		OutputRefType: "safety_boundary",
		OutputRef:     outputRef,
		CreatedAt:     clock.SystemClock{}.NowRFC3339(),
	}); err != nil {
		return "", err
	}
	return outputRef, nil
}

type recoverySmokeFactCounts struct {
	Decisions    int
	Audits       int
	Portfolios   int
	Reports      int
	Intelligence int
}

func (c recoverySmokeFactCounts) totalRestoredFacts() int {
	return c.Decisions + c.Audits + c.Portfolios + c.Reports + c.Intelligence
}

func recoverySmokeCounts(db *sql.DB) (recoverySmokeFactCounts, error) {
	var counts recoverySmokeFactCounts
	for _, item := range []struct {
		table string
		value *int
	}{
		{table: "decision_records", value: &counts.Decisions},
		{table: "audit_events", value: &counts.Audits},
		{table: "portfolio_snapshots", value: &counts.Portfolios},
		{table: "daily_discipline_reports", value: &counts.Reports},
		{table: "intelligence_summary", value: &counts.Intelligence},
	} {
		if err := db.QueryRow("SELECT COUNT(*) FROM " + item.table).Scan(item.value); err != nil {
			return recoverySmokeFactCounts{}, err
		}
	}
	return counts, nil
}

func copyFile(src, dst string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func executeTask(ctx context.Context, cfg *config.Config, rt agentRuntime, task string, period string, source string, symbol string, startDate, endDate time.Time) error {
	requestID := idgen.NewGenerator().New("req")
	switch task {
	case "daily":
		dailyCtx, err := buildDailyWorkflowContext(ctx, rt, requestID)
		if err != nil {
			return err
		}
		out, err := workflow.NewDailyDisciplineGraphWithDependencies(rt.deps).Run(ctx, dailyCtx)
		if err != nil {
			return err
		}
		return upsertManualDailyDisciplineReport(ctx, rt.repos, rt.transactor, cfg.DailyAutoRun.Timezone, requestID, out)
	case "market-refresh":
		_, err := workflow.NewMarketRefreshGraphWithDependencies(rt.deps).Run(ctx, workflow.MarketRefreshInput{RequestID: requestID, Symbol: symbol})
		return err
	case "evidence-index":
		evidenceSvc := service.NewEvidenceService(rt.transactor)
		chunkCount, err := evidenceSvc.CountRAGChunks(ctx)
		if err != nil {
			return err
		}
		if chunkCount == 0 {
			if _, err := workflow.NewEvidenceVerificationGraphWithDependencies(rt.deps).Run(ctx, workflow.EvidenceVerificationInput{RequestID: requestID, Symbol: "510300", Sources: []string{"stub-a", "stub-b"}}); err != nil {
				return err
			}
		}
		if _, err := evidenceSvc.RebuildVectorIndexWithStats(ctx, service.NewFileVectorIndex(cfg.VecLite.Path)); err != nil {
			return err
		}
		_, err = evidenceSvc.AppendRebuildAudit(ctx, requestID)
		return err
	case "public-evidence-refresh":
		return runPublicEvidenceRefresh(ctx, cfg, rt, requestID, symbol, startDate, endDate)
	case "p34-expanded-refresh":
		return runP34ExpandedRefresh(ctx, cfg, rt, source, symbol, startDate, endDate)
	case "llm-smoke":
		return runLLMSmoke(ctx, cfg, rt, symbol)
	case "review":
		_, err := service.NewQueryService(rt.repos).ReviewSummary(ctx, period)
		return err
	default:
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "未知本地任务")
	}
}

func runPublicEvidenceRefresh(ctx context.Context, cfg *config.Config, rt agentRuntime, requestID string, symbol string, start, end time.Time) error {
	if !cfg.DataSources.PublicEvidence.Enabled {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "public evidence collector is disabled by data_sources.public_evidence.enabled")
	}
	collector := &workflow.CompositePublicEvidenceCollector{Collectors: publicEvidenceCollectors(cfg)}
	service := workflow.PublicEvidenceIngestionService{Collector: collector, IntelligenceRepo: rt.repos.IntelligenceRepo, AuditRepo: rt.repos.AuditRepo, GenerateAuditID: func() string { return idgen.NewGenerator().New("audit") }, RequestID: requestID}
	if end.IsZero() {
		end = time.Now().UTC()
	}
	if start.IsZero() {
		start = end.AddDate(0, 0, -90)
	}
	return service.IngestPublicEvidence(ctx, symbol, start, end)
}

func runP34ExpandedRefresh(ctx context.Context, cfg *config.Config, rt agentRuntime, source string, symbol string, start, end time.Time) error {
	deps := rt.deps
	source = strings.TrimSpace(source)
	if source == "" {
		source = "configured"
	}
	switch source {
	case "configured":
	case "csindex_extended":
		deps.MarketDataSource = workflow.CsindexCollector{BaseURL: cfg.DataSources.MarketCollectors.CSIndexBaseURL, IncludeExtended: true}
	case "sentiment_proxy_fixture":
		dataDate := time.Now().UTC().Format(time.DateOnly)
		if !end.IsZero() {
			dataDate = end.Format(time.DateOnly)
		} else if !start.IsZero() {
			dataDate = start.Format(time.DateOnly)
		}
		deps.MarketDataSource = workflow.FixtureSentimentProxyCollector{Fixtures: map[string]workflow.SentimentProxyPoint{
			symbol: {SourceName: "sentiment_proxy_fixture", SourceLevel: model.SourceLevelC, DataDate: dataDate, HeatScore: 62, SentimentState: model.SentimentNeutral, Raw: map[string]any{"window_start": formatOptionalDate(start), "window_end": formatOptionalDate(end)}},
		}}
	default:
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "unsupported P34 source")
	}
	_, err := workflow.NewMarketRefreshGraphWithDependencies(deps).Run(ctx, workflow.MarketRefreshInput{RequestID: idgen.NewGenerator().New("req"), Symbol: symbol})
	return err
}

func runLLMSmoke(ctx context.Context, cfg *config.Config, rt agentRuntime, symbol string) error {
	if strings.TrimSpace(cfg.DeepSeek.APIKey) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "llm-smoke requires deepseek.api_key in config file")
	}
	if strings.TrimSpace(symbol) == "" {
		return apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "llm-smoke requires --symbol")
	}
	resp, err := rt.deps.AnalystService.Analyze(ctx, workflow.AnalystRequest{
		AgentName:       "value",
		Symbol:          symbol,
		EvidenceSummary: "P37 real LLM smoke: 本地最小样本，仅验证模型调用、解析、质量门禁和审计记录。",
		PositionContext: "本 smoke 不读取、不写入账户或持仓，不创建确认单或交易流水。",
		RuleBoundary:    "LLM 只生成分析材料，最终裁决由规则引擎负责；不得输出交易指令、收益承诺或最终裁决。",
	})
	if err != nil {
		return err
	}
	if resp.Metadata["parse_status"] != "parsed" || resp.Metadata["quality_status"] != "passed" {
		return apperr.New(apperr.CodeAnalystUnavailable, apperr.CategoryInternal, "llm-smoke metadata did not pass parse and quality checks")
	}
	return nil
}

func runRetrievalQualitySmoke(ctx context.Context, rt agentRuntime, symbol string) (string, error) {
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return "", apperr.New(apperr.CodeBadRequest, apperr.CategoryBadRequest, "retrieval-quality-smoke requires --symbol")
	}
	retrieval := rt.deps.RetrievalService
	if retrieval == nil {
		retrieval = service.NewRetrievalAdapter(rt.transactor, nil)
	}
	result, err := retrieval.RetrieveEvidence(ctx, workflow.RetrievalRequest{Symbol: symbol})
	if err != nil {
		return "", err
	}
	return retrievalQualitySmokeOutputRef(result.QualitySummary), nil
}

func runDataSourceQualityRegression(ctx context.Context, rt agentRuntime, source string, symbol string, strictQualityGate bool) (string, error) {
	mode := strings.TrimSpace(source)
	if mode == "" {
		mode = service.DataSourceQualityModeFixture
	}
	out, err := service.NewDataSourceQualityService(rt.repos).Run(ctx, service.DataSourceQualityRegressionRequest{Mode: mode, Symbol: symbol})
	if err != nil {
		return "", err
	}
	outputRef := service.DataSourceQualityAuditOutputRef(out)
	if strictQualityGate && out.Policy.ReleaseGate == service.DataSourceQualityReleaseGateBlock {
		return outputRef, fmt.Errorf("current data quality policy gate blocked: %s", outputRef)
	}
	return outputRef, nil
}

func runDataSourceQualityResolutionCheck(ctx context.Context, rt agentRuntime, symbol string) (string, error) {
	check, err := service.NewDataSourceQualityService(rt.repos).CheckGateResolution(ctx, service.DataQualityGateResolutionCheckRequest{Symbol: symbol})
	if err != nil {
		return "", err
	}
	outputRef := dataQualityGateResolutionAuditOutputRef(check)
	if check.ReleaseClaimState == service.DataQualityReleaseClaimRequiresResolution {
		return outputRef, fmt.Errorf("current data quality resolution required: %s", outputRef)
	}
	return outputRef, nil
}

func dataQualityGateResolutionAuditOutputRef(check dto.DataQualityGateResolutionCheck) string {
	resolutionType := "none"
	if check.ActiveResolution != nil {
		resolutionType = check.ActiveResolution.ResolutionType
	}
	return strings.Join([]string{
		"data_quality_gate_resolution",
		"claim_state=" + firstNonEmpty(check.ReleaseClaimState, service.DataQualityReleaseClaimRequiresResolution),
		"policy=" + firstNonEmpty(check.Policy.Verdict, service.DataSourceQualityPolicyBlocked),
		"gate=" + firstNonEmpty(check.Policy.ReleaseGate, service.DataSourceQualityReleaseGateBlock),
		"fingerprint=" + safeAuditFingerprint(check.PolicyFingerprint),
		"resolution=" + firstNonEmpty(resolutionType, "none"),
		fmt.Sprintf("clean_data_claim=%t", check.CleanDataClaimAllowed),
		"no_auto_trading",
	}, ":")
}

func retrievalQualitySmokeOutputRef(summary workflow.RetrievalQualitySummary) string {
	return strings.Join([]string{
		"retrieval_quality",
		"status=" + firstNonEmpty(summary.Status, "unknown"),
		fmt.Sprintf("topk=%d", summary.TopK),
		"fallback=" + firstNonEmpty(summary.FallbackSource, "unknown"),
		"index=" + firstNonEmpty(summary.IndexHealth, "unknown"),
		"consistency=" + firstNonEmpty(summary.SourceConsistencyStatus, "unknown"),
		"no_auto_trading",
	}, ":")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func formatOptionalDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(time.DateOnly)
}

func parsePublicEvidenceDateWindow(startText, endText string) (time.Time, time.Time, error) {
	startText = strings.TrimSpace(startText)
	endText = strings.TrimSpace(endText)
	var start, end time.Time
	var err error
	if startText != "" {
		start, err = time.ParseInLocation(time.DateOnly, startText, time.UTC)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("--start-date must use YYYY-MM-DD")
		}
	}
	if endText != "" {
		end, err = time.ParseInLocation(time.DateOnly, endText, time.UTC)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("--end-date must use YYYY-MM-DD")
		}
	}
	if !start.IsZero() && !end.IsZero() && start.After(end) {
		return time.Time{}, time.Time{}, fmt.Errorf("--start-date must not be after --end-date")
	}
	return start, end, nil
}

func publicEvidenceCollectors(cfg *config.Config) []workflow.PublicEvidenceCollector {
	publicCfg := cfg.DataSources.PublicEvidence
	if len(publicCfg.Sources) == 0 {
		publicCfg.Sources = []string{"cninfo", "szse", "csrc"}
	}
	collectors := make([]workflow.PublicEvidenceCollector, 0, len(publicCfg.Sources))
	for _, source := range publicCfg.Sources {
		switch strings.TrimSpace(source) {
		case "cninfo":
			collectors = append(collectors, &workflow.CninfoCollector{BaseURL: publicCfg.CNInfoBaseURL, OrgIDBySymbol: publicCfg.CNInfoOrgIDs})
		case "szse":
			collectors = append(collectors, &workflow.SzseCollector{BaseURL: publicCfg.SZSEBaseURL})
		case "csrc":
			collectors = append(collectors, &workflow.CsrcCollector{BaseURL: publicCfg.CSRCBaseURL})
		case "csindex_index":
			collectors = append(collectors, workflow.CsindexIndexEvidenceCollector{BaseURL: cfg.DataSources.MarketCollectors.CSIndexBaseURL})
		case "eastmoney_fund":
			collectors = append(collectors, workflow.EastmoneyFundEvidenceCollector{BaseURL: cfg.DataSources.MarketCollectors.EastmoneyFundBaseURL})
		}
	}
	return collectors
}

func buildDailyWorkflowContext(ctx context.Context, rt agentRuntime, requestID string) (workflow.WorkflowContext, error) {
	portfolio, err := rt.repos.PortfolioRepo.GetLatestPortfolioSnapshot(ctx)
	if err != nil {
		return workflow.WorkflowContext{}, apperr.New(apperr.CodeDataRequired, apperr.CategoryConflict, "每日纪律缺少账户快照")
	}
	market, err := rt.repos.MarketRepo.GetLatestMarketSnapshot(ctx)
	if err != nil {
		return workflow.WorkflowContext{}, apperr.New(apperr.CodeDataStale, apperr.CategoryConflict, "每日纪律缺少市场快照")
	}
	rule, err := rt.repos.RuleRepo.GetActiveRuleVersion(ctx)
	if err != nil {
		return workflow.WorkflowContext{}, apperr.New(apperr.CodeRuleVersionMissing, apperr.CategoryConflict, "每日纪律缺少生效规则版本")
	}
	positions, err := rt.repos.PortfolioRepo.ListPositions(ctx)
	if err != nil {
		return workflow.WorkflowContext{}, err
	}
	dailyPositionSnapshots := dailyPositions(positions)
	return workflow.WorkflowContext{
		RequestID:                 requestID,
		WorkflowType:              workflow.WorkflowDailyDiscipline,
		Symbol:                    market.Symbol,
		RuleVersion:               rule.RuleVersion,
		CapabilityStatus:          workflow.CapabilityInScope,
		PortfolioSnapshot:         model.PortfolioSnapshot{SnapshotID: portfolio.SnapshotID, Cash: portfolio.Cash, TotalAssets: portfolio.TotalAssets, CashRatio: portfolio.CashRatio, HighRiskRatio: portfolio.HighRiskRatio, PositionCount: portfolio.PositionCount},
		PositionSnapshots:         dailyPositionSnapshots,
		MarketSnapshot:            market,
		AnalystUnavailable:        true,
		ExpectedReturnSampleCount: workflow.ExpectedReturnSampleCountFromWorkflowData(dailyPositionSnapshots, market),
	}, nil
}

func dailyPositions(items []repository.Position) []model.Position {
	out := make([]model.Position, 0, len(items))
	for _, item := range items {
		out = append(out, model.Position{PositionID: item.PositionID, Symbol: item.Symbol, Name: item.Name, Quantity: item.Quantity, CostPrice: item.CostPrice, CurrentPrice: item.CurrentPrice, MarketValue: item.MarketValue, UnrealizedProfitRatio: item.UnrealizedProfitRatio, PositionState: model.PositionState(item.PositionState), AssetTag: item.AssetTag})
	}
	return out
}

var manualDailyDisciplineReportNow = time.Now

func upsertManualDailyDisciplineReport(ctx context.Context, repos repository.Repositories, tx repository.Transactor, timezone string, requestID string, out workflow.WorkflowContext) error {
	if repos.DailyDisciplineReportRepo == nil {
		return nil
	}
	now := clock.SystemClock{}.NowRFC3339()
	localDate := manualDailyDisciplineReportLocalDate(manualDailyDisciplineReportNow(), timezone)
	report := repository.DailyDisciplineReport{ReportID: manualDailyDisciplineReportID(localDate, requestID), LocalDate: localDate, Scope: "holdings", SymbolSetHash: manualDailyDisciplineSymbolSetHash(out.PositionSnapshots), SourceType: "manual", SourceID: requestID, DecisionID: out.DecisionID, Status: manualDailyDisciplineReportStatus(ctx, repos.DecisionRepo, out), Summary: "今日纪律报告已生成", CreatedAt: now, UpdatedAt: now}
	if err := repos.DailyDisciplineReportRepo.UpsertDailyDisciplineReport(ctx, report); err != nil {
		return err
	}
	return triggerManualDailyRiskAlerts(ctx, repos, tx, requestID, report, out)
}

func triggerManualDailyRiskAlerts(ctx context.Context, repos repository.Repositories, tx repository.Transactor, requestID string, report repository.DailyDisciplineReport, out workflow.WorkflowContext) error {
	if tx == nil || repos.DecisionRepo == nil || repos.RiskAlertRepo == nil || strings.TrimSpace(report.DecisionID) == "" {
		return nil
	}
	decision, _, err := repos.DecisionRepo.GetDecisionRecord(ctx, report.DecisionID)
	if err != nil {
		return err
	}
	riskSvc := service.NewRiskAlertService(tx)
	inputs := riskSvc.BuildRiskAlertTriggers(decision, out.MarketSnapshot, service.SourceHealthRiskInputsFromExpectedReturnJSON(decision.ExpectedReturnScenariosJSON))
	for _, input := range inputs {
		input.ReportID = report.ReportID
		input.RequestID = requestID
		if _, err := riskSvc.TriggerRiskAlert(ctx, input); err != nil {
			return err
		}
	}
	return nil
}

func manualDailyDisciplineReportLocalDate(now time.Time, timezone string) string {
	loc := time.UTC
	if strings.TrimSpace(timezone) != "" {
		loaded, err := time.LoadLocation(timezone)
		if err == nil {
			loc = loaded
		}
	}
	return now.In(loc).Format(time.DateOnly)
}

func manualDailyDisciplineReportStatus(ctx context.Context, repo repository.DecisionRepository, out workflow.WorkflowContext) string {
	if status := workflow.DailyDisciplineReportStatus(out); status == "degraded" {
		return status
	}
	if repo == nil || out.DecisionID == "" {
		return "success"
	}
	decision, _, err := repo.GetDecisionRecord(ctx, out.DecisionID)
	if err == nil && decision.WorkflowStatus == string(model.WorkflowDegraded) {
		return "degraded"
	}
	return "success"
}

func manualDailyDisciplineReportID(localDate string, requestID string) string {
	return "daily_report_manual_" + localDate + "_" + requestID
}

func manualDailyDisciplineSymbolSetHash(positions []model.Position) string {
	if len(positions) == 0 {
		return "manual"
	}
	symbols := make([]string, 0, len(positions))
	for _, position := range positions {
		if strings.TrimSpace(position.Symbol) != "" {
			symbols = append(symbols, strings.TrimSpace(position.Symbol))
		}
	}
	if len(symbols) == 0 {
		return "manual"
	}
	sort.Strings(symbols)
	sum := sha256.Sum256([]byte(strings.Join(symbols, ",")))
	return hex.EncodeToString(sum[:8])
}

func appendTaskAudit(ctx context.Context, repo repository.AuditRepository, task string, period string, source string, symbol string, llmModel string, startDate, endDate time.Time, status string, errCode string, outputRef string) error {
	action := supportedTasks[task]
	inputRef := task
	if task == "review" {
		inputRef = task + ":" + period
	}
	if task == "public-evidence-refresh" {
		inputRef = publicEvidenceAuditInputRef(task, symbol, startDate, endDate)
	}
	if task == "p34-expanded-refresh" {
		inputRef = p34ExpandedAuditInputRef(task, source, symbol, startDate, endDate)
	}
	if task == "llm-smoke" {
		inputRef = llmSmokeAuditInputRef(task, symbol, llmModel)
		if status == string(model.AuditStatusSuccess) && outputRef == "no_auto_trading" {
			outputRef = "llm_smoke:quality=passed:parse=parsed:no_auto_trading"
		}
	}
	if task == "retrieval-quality-smoke" {
		inputRef = retrievalQualitySmokeAuditInputRef(task, symbol)
	}
	if task == "data-source-quality-regression" {
		inputRef = dataSourceQualityAuditInputRef(task, source, symbol)
	}
	if task == "data-source-quality-resolution-check" {
		inputRef = dataQualityGateResolutionAuditInputRef(task, symbol)
	}
	if status == string(model.AuditStatusFailed) && errCode == "" {
		errCode = string(apperr.CodeInternalError)
	}
	// 审计写入记录本地任务输入摘要、执行结果和安全边界；本入口不更新账户、不创建下单请求。
	return repo.AppendAuditEvent(ctx, repository.AuditEvent{
		AuditEventID:  idgen.NewGenerator().New("audit"),
		RequestID:     idgen.NewGenerator().New("req"),
		WorkflowType:  task,
		NodeName:      "cmd_agent",
		Actor:         string(model.AuditActorUser),
		Action:        string(action),
		NodeAction:    "manual_local_task",
		Status:        status,
		ErrorCode:     errCode,
		InputRefType:  "agent_task",
		InputRef:      inputRef,
		OutputRefType: "safety_boundary",
		OutputRef:     outputRef,
		CreatedAt:     clock.SystemClock{}.NowRFC3339(),
	})
}

func llmSmokeAuditInputRef(task string, symbol string, llmModel string) string {
	return strings.Join([]string{task, "symbol=" + strings.TrimSpace(symbol), "model=" + strings.TrimSpace(llmModel)}, ":")
}

func retrievalQualitySmokeAuditInputRef(task string, symbol string) string {
	return strings.Join([]string{task, "symbol=" + strings.TrimSpace(symbol)}, ":")
}

func dataSourceQualityAuditInputRef(task string, source string, symbol string) string {
	mode := strings.TrimSpace(source)
	if mode == "" {
		mode = service.DataSourceQualityModeFixture
	}
	switch mode {
	case service.DataSourceQualityModeFixture, service.DataSourceQualityModeCurrent:
	default:
		mode = "unsupported"
	}
	return strings.Join([]string{task, "source=" + mode, "symbol=" + safeAuditSymbol(symbol)}, ":")
}

func dataQualityGateResolutionAuditInputRef(task string, symbol string) string {
	return strings.Join([]string{task, "symbol=" + safeAuditSymbol(symbol)}, ":")
}

func safeAuditSymbol(symbol string) string {
	symbol = strings.TrimSpace(symbol)
	if symbol == "" {
		return ""
	}
	if len(symbol) > 32 {
		return "redacted"
	}
	for _, r := range symbol {
		if (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_' || r == '-' || r == '.' {
			continue
		}
		return "redacted"
	}
	return symbol
}

func safeAuditFingerprint(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) > 96 {
		return "redacted"
	}
	for _, r := range value {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') {
			continue
		}
		return "redacted"
	}
	return value
}

func publicEvidenceAuditInputRef(task string, symbol string, startDate, endDate time.Time) string {
	parts := []string{task, "symbol=" + strings.TrimSpace(symbol)}
	if !startDate.IsZero() {
		parts = append(parts, "start="+startDate.Format(time.DateOnly))
	}
	if !endDate.IsZero() {
		parts = append(parts, "end="+endDate.Format(time.DateOnly))
	}
	return strings.Join(parts, ":")
}

func p34ExpandedAuditInputRef(task string, source string, symbol string, startDate, endDate time.Time) string {
	parts := []string{task, "source=" + strings.TrimSpace(source), "symbol=" + strings.TrimSpace(symbol)}
	if !startDate.IsZero() {
		parts = append(parts, "start="+startDate.Format(time.DateOnly))
	}
	if !endDate.IsZero() {
		parts = append(parts, "end="+endDate.Format(time.DateOnly))
	}
	return strings.Join(parts, ":")
}

func errorCode(err error) string {
	if appErr, ok := apperr.AsAppError(err); ok {
		return string(appErr.Code)
	}
	return string(apperr.CodeInternalError)
}
