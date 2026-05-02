package execution

import "context"

type Repository interface {
	Save(ctx context.Context, execution Execution) error
	FindByJobID(ctx context.Context, jobID string) ([]Execution, error)
	FindAll(ctx context.Context, limit int, offset int) ([]Execution, int64, error)
	Update(ctx context.Context, execution Execution) error
}
