package execution

import "gorm.io/gorm"

type PostgresRepository struct {
	db *gorm.DB
}
func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Save(execution Execution) error {
	return r.db.Create(&execution).Error
}

func (r *PostgresRepository) FindByJobID(jobID uint) ([]Execution, error) {
	var executions []Execution
	err := r.db.Where("job_id = ?", jobID).Find(&executions).Error
	return executions, err
}
