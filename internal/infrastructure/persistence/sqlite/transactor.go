package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"investment-agent/internal/domain/repository"
	"investment-agent/internal/pkg/apperr"
)

// Transactor coordinates SQLite-backed repository writes in one transaction.
type Transactor struct{ db *sql.DB }

// NewTransactor creates a SQLite transaction coordinator.
func NewTransactor(db *sql.DB) *Transactor { return &Transactor{db: db} }

// WithinTx runs fn with repositories backed by one *sql.Tx.
func (t *Transactor) WithinTx(ctx context.Context, fn func(context.Context, repository.Repositories) error) error {
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		return apperr.FromRepositoryError(err)
	}
	repos := repository.Repositories{
		PortfolioRepo:                 &PortfolioRepository{db: txDB{tx}},
		DecisionRepo:                  &DecisionRepository{db: txDB{tx}},
		AuditRepo:                     &AuditRepository{db: txDB{tx}},
		RuleRepo:                      &RuleRepository{db: txDB{tx}},
		MarketRepo:                    &MarketRepository{db: txDB{tx}},
		SettingsRepo:                  &SettingsRepository{db: txDB{tx}},
		IntelligenceRepo:              &IntelligenceRepository{db: txDB{tx}},
		NotificationRepo:              &NotificationRepository{db: txDB{tx}},
		DailyAutoRunRepo:              &DailyAutoRunRepository{db: txDB{tx}},
		DailyDisciplineReportRepo:     &DailyDisciplineReportRepository{db: txDB{tx}},
		RiskAlertRepo:                 &RiskAlertRepository{db: txDB{tx}},
		RuleEffectRepo:                &RuleEffectRepository{db: txDB{tx}},
		DataQualityGateResolutionRepo: &DataQualityGateResolutionRepository{db: txDB{tx}},
	}
	if err := fn(ctx, repos); err != nil {
		return errors.Join(err, tx.Rollback())
	}
	return apperr.FromRepositoryError(tx.Commit())
}
