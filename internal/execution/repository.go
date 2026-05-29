package execution

import (
	"context"
	"errors"

	"github.com/winnerx0/kron/internal/domain"
)

var ErrNotFound = errors.New("execution not found")

type Repository interface {
	Save(ctx context.Context, execution domain.Execution) error
	FindByJobID(ctx context.Context, jobID string) ([]domain.Execution, error)
	FindAll(ctx context.Context, limit int, offset int, jobID string) ([]ExecutionDTO, int64, error)
	FindByID(ctx context.Context, id string) (ExecutionDetailDTO, error)
	Update(ctx context.Context, execution domain.Execution) error
}
