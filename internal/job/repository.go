package job

import (
	"context"
)

type Repository interface {
	
	FindAll(ctx context.Context) ([]Job, error)
	
	FindAllNextRun(ctx context.Context) ([]Job, error)
	
	Create(ctx context.Context, job Job) (Job, error)
	
	Update(ctx context.Context, job Job) (Job, error)
	
	Delete(ctx context.Context, id string) error
}
