package execution

import "context"

type Repository interface {
	Save(ctx context.Context, execution Execution) error
	FindByJobID(ctx context.Context, jobID string) ([]Execution, error)
	FindAll(ctx context.Context) ([]Execution, error)
	Update(ctx context.Context, execution Execution) error
}