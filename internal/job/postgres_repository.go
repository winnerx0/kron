package job

import (
	"context"

	"gorm.io/gorm"
)

type PostgresRepository struct {
	
	db *gorm.DB
}
func NewRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) FindAll(ctx context.Context) ([]Job, error) {
	jobs, err := gorm.G[Job](r.db).Where("next_run_at <= NOW()").Find(ctx)
	
	if err != nil {
		return []Job{}, err
	}
	
	return jobs, nil
}

func (r *PostgresRepository) Create(ctx context.Context, job Job) (Job, error) {
	if err := r.db.WithContext(ctx).Create(&job).Error; err != nil {
		return Job{}, err
	}
	return job, nil
}

func (r *PostgresRepository) Update(ctx context.Context, job Job) error {
	return r.db.WithContext(ctx).Save(&job).Error
}
