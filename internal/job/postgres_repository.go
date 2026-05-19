package job

import (
	"context"

	"github.com/winnerx0/kron/internal/domain"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}

}

func (r *PostgresRepository) FindAll(ctx context.Context, userID string) ([]domain.Job, error) {
	jobs, err := gorm.G[domain.Job](r.db).Where("user_id = ?", userID).Find(ctx)

	if err != nil {
		return []domain.Job{}, err
	}

	if len(jobs) == 0 {
		return []domain.Job{}, nil
	}

	return jobs, nil
}

func (r *PostgresRepository) FindAllNextRun(ctx context.Context) ([]domain.Job, error) {

	var jobs []domain.Job

	err := r.db.
		Select("jobs.*").
		Joins("LEFT JOIN executions ON executions.job_id = jobs.id AND executions.status = ?", "running").
		Where("jobs.next_run_at <= NOW() AND jobs.status = ? AND executions.id IS NULL", true).
		Find(&jobs).Error

	if err != nil {
		return []domain.Job{}, err
	}

	if len(jobs) == 0 {
		return []domain.Job{}, nil
	}

	return jobs, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (domain.Job, error) {
	return gorm.G[domain.Job](r.db).Where("id = ?", id).First(ctx)
}

func (r *PostgresRepository) Create(ctx context.Context, job domain.Job) (domain.Job, error) {
	if err := r.db.WithContext(ctx).Create(&job).Error; err != nil {
		return domain.Job{}, err
	}
	return job, nil
}

func (r *PostgresRepository) Update(ctx context.Context, job domain.Job) (domain.Job, error) {
	if err := r.db.WithContext(ctx).Save(&job).Error; err != nil {
		return domain.Job{}, err
	}
	return job, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, jobID string) error {
	return r.db.WithContext(ctx).Delete(&domain.Job{ID: jobID}).Error
}
