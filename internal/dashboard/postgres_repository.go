package dashboard

import (
	"context"

	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/job"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Counts(ctx context.Context, userID string) (Counts, error) {
	var counts Counts

	if err := r.db.WithContext(ctx).
		Model(&job.Job{}).
		Where("user_id = ?", userID).
		Count(&counts.TotalJobs).Error; err != nil {
		return counts, err
	}

	if err := r.db.WithContext(ctx).
		Model(&job.Job{}).
		Where("user_id = ? AND status = ?", userID, true).
		Count(&counts.EnabledJobs).Error; err != nil {
		return counts, err
	}

	if err := executionQuery(r.db.WithContext(ctx), userID).
		Count(&counts.TotalExecutions).Error; err != nil {
		return counts, err
	}

	if err := executionQuery(r.db.WithContext(ctx), userID).
		Where("executions.status = ?", execution.SUCCESS).
		Count(&counts.SuccessfulExecutions).Error; err != nil {
		return counts, err
	}

	if err := executionQuery(r.db.WithContext(ctx), userID).
		Where("executions.status = ?", execution.FAILED).
		Count(&counts.FailedExecutions).Error; err != nil {
		return counts, err
	}

	return counts, nil
}

func (r *PostgresRepository) RecentExecutions(ctx context.Context, userID string, limit int) ([]RecentExecution, error) {
	var executions []RecentExecution

	err := executionQuery(r.db.WithContext(ctx), userID).
		Select(`
			executions.id,
			executions.job_id,
			jobs.name AS job_name,
			executions.status,
			executions.started AS started_at,
			executions.finished AS finished_at
		`).
		Order("executions.started DESC").
		Limit(limit).
		Scan(&executions).Error

	return executions, err
}

func (r *PostgresRepository) Jobs(ctx context.Context, userID string, limit int) ([]JobSummary, error) {
	var jobs []JobSummary

	err := r.db.WithContext(ctx).
		Model(&job.Job{}).
		Select("id, name, schedule, status AS enabled").
		Where("user_id = ?", userID).
		Order("next_run_at ASC").
		Limit(limit).
		Scan(&jobs).Error

	return jobs, err
}

func executionQuery(db *gorm.DB, userID string) *gorm.DB {
	return db.Model(&execution.Execution{}).
		Joins("JOIN jobs ON executions.job_id = jobs.id").
		Where("jobs.user_id = ?", userID)
}
