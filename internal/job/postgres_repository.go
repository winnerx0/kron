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
	jobs, err := gorm.G[Job](r.db).Find(ctx)

	if err != nil {
		return []Job{}, err
	}

	if len(jobs) == 0 {
		return []Job{}, nil
	}

	return jobs, nil
}

func (r *PostgresRepository) FindAllNextRun(ctx context.Context) ([]Job, error) {
	jobs, err := gorm.G[Job](r.db).Where("next_run_at <= NOW()").Find(ctx)

	if err != nil {
		return []Job{}, err
	}

	if len(jobs) == 0 {
		return []Job{}, nil
	}

	return jobs, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id string) (Job, error) {
	return gorm.G[Job](r.db).Where("id = ?", id).First(ctx)
}

func (r *PostgresRepository) Create(ctx context.Context, job Job) (Job, error) {
	if err := r.db.WithContext(ctx).Create(&job).Error; err != nil {
		return Job{}, err
	}
	return job, nil
}

func (r *PostgresRepository) Update(ctx context.Context, job Job) (Job, error) {
	if err := r.db.WithContext(ctx).Save(&job).Error; err != nil {
		return Job{}, err
	}
	return job, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, jobID string) error {
	return r.db.WithContext(ctx).Delete(&Job{ID: jobID}).Error
}
