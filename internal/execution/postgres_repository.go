package execution

import (
	"context"

	"github.com/winnerx0/kron/internal/domain"
	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}

func NewPostgresRepository(db *gorm.DB) Repository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(ctx context.Context, execution domain.Execution) error {
	return r.db.WithContext(ctx).Create(&execution).Error
}

func (r *PostgresRepository) FindByJobID(ctx context.Context, jobID string) ([]domain.Execution, error) {
	var executions []domain.Execution
	err := r.db.WithContext(ctx).Where("job_id = ?", jobID).Find(&executions).Error
	return executions, err
}

func (r *PostgresRepository) FindAll(ctx context.Context, limit int, offset int, jobID string) ([]domain.Execution, int64, error) {
	var executions []domain.Execution
	var total int64

	userID, ok := ctx.Value("userId").(string)
	if !ok {
		return executions, 0, nil
	}

	query := r.db.WithContext(ctx).Model(&domain.Execution{}).
		Joins("JOIN jobs ON executions.job_id = jobs.id").
		Where("jobs.user_id = ?", userID)

	if jobID != "" {
		query = query.Where("executions.job_id = ?", jobID)
	}

	if err := query.Count(&total).Error; err != nil {
		return executions, 0, err
	}

	selectQuery := r.db.WithContext(ctx).
		Joins("JOIN jobs ON executions.job_id = jobs.id").
		Where("jobs.user_id = ?", userID)

	if jobID != "" {
		selectQuery = selectQuery.Where("executions.job_id = ?", jobID)
	}

	err := selectQuery.Order("started DESC").
		Limit(limit).
		Offset(offset).
		Find(&executions).Error
	return executions, total, err
}

func (r *PostgresRepository) Update(ctx context.Context, execution domain.Execution) error {
	return r.db.WithContext(ctx).Model(&domain.Execution{}).Where("id = ?", execution.ID).Updates(execution).Error
}
