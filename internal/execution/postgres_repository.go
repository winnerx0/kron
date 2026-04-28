package execution

import (
	"context"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	db *gorm.DB
}
func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(ctx context.Context, execution Execution) error {
	return r.db.WithContext(ctx).Create(&execution).Error
}

func (r *PostgresRepository) FindByJobID(ctx context.Context, jobID string) ([]Execution, error) {
	var executions []Execution
	err := r.db.WithContext(ctx).Where("job_id = ?", jobID).Find(&executions).Error
	return executions, err
}

func (r *PostgresRepository) FindAll(ctx context.Context) ([]Execution, error) {
	var executions []Execution
	err := r.db.WithContext(ctx).Find(&executions).Error
	return executions, err
}

func (r *PostgresRepository) Update(ctx context.Context, execution Execution) error {
	return r.db.WithContext(ctx).Model(&Execution{}).Where("id = ?", execution.ID).Updates(execution).Error
}