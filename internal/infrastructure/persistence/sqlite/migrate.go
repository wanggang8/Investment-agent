package sqlite

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"embed"
	"fmt"
	"sort"
	"strings"

	moderncsqlite "modernc.org/sqlite"
)

// migrationFiles 内嵌 migration 目录下的 SQL 文件，保证二进制可独立执行迁移。
//
//go:embed migration/*.sql
var migrationFiles embed.FS

// Store 封装 SQLite 连接，后续 Repository 从 DB 派生。
type Store struct {
	DB *sql.DB
}

func init() {
	moderncsqlite.RegisterConnectionHook(configureSQLiteConnection)
}

// Open 打开 SQLite 数据库文件；path 可使用 :memory: 进行测试。
func Open(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := verifyConnection(db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return &Store{DB: db}, nil
}

func verifyConnection(db *sql.DB) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("ping sqlite: %w", err)
	}
	return nil
}

func configureSQLiteConnection(conn moderncsqlite.ExecQuerierContext, dsn string) error {
	pragmas := []string{
		`PRAGMA foreign_keys = ON`,
		`PRAGMA busy_timeout = 5000`,
	}
	if !isMemorySQLitePath(dsn) {
		pragmas = append(pragmas, `PRAGMA journal_mode = WAL`)
	}
	for _, statement := range pragmas {
		if _, err := conn.ExecContext(context.Background(), statement, []driver.NamedValue{}); err != nil {
			return fmt.Errorf("configure sqlite %s: %w", statement, err)
		}
	}
	return nil
}

func isMemorySQLitePath(path string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(path))
	return trimmed == ":memory:" || strings.Contains(trimmed, "mode=memory")
}

// Close 关闭底层数据库连接。
func (s *Store) Close() error {
	return s.DB.Close()
}

// Migrate 按文件名顺序执行所有 SQL migration。
// 迁移文件使用 IF NOT EXISTS / INSERT OR IGNORE，重复执行保持幂等。
func Migrate(ctx context.Context, db *sql.DB) error {
	files, err := migrationFiles.ReadDir("migration")
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	names := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			names = append(names, file.Name())
		}
	}
	sort.Strings(names)

	for _, name := range names {
		content, err := migrationFiles.ReadFile("migration/" + name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}
		if _, err := db.ExecContext(ctx, string(content)); err != nil {
			return fmt.Errorf("apply migration %s: %w", name, err)
		}
	}
	if err := ensureEvidenceRefQualityColumns(ctx, db); err != nil {
		return err
	}
	return nil
}

func ensureEvidenceRefQualityColumns(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, `PRAGMA table_info(evidence_refs)`)
	if err != nil {
		return fmt.Errorf("inspect evidence_refs: %w", err)
	}
	defer rows.Close()
	columns := map[string]bool{}
	for rows.Next() {
		var cid int
		var name, typ string
		var notNull int
		var defaultValue any
		var pk int
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &pk); err != nil {
			return fmt.Errorf("scan evidence_refs column: %w", err)
		}
		columns[name] = true
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("inspect evidence_refs columns: %w", err)
	}
	if !columns["independent_source_count"] {
		if _, err := db.ExecContext(ctx, `ALTER TABLE evidence_refs ADD COLUMN independent_source_count INTEGER NOT NULL DEFAULT 0`); err != nil {
			return fmt.Errorf("add evidence_refs.independent_source_count: %w", err)
		}
	}
	return nil
}
