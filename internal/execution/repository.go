package execution

import (
	"context"

	"github.com/winnerx0/kron/internal/domain"
)

type Repository interface {
	Save(ctx context.Context, execution domain.Execution) error
	FindByJobID(ctx context.Context, jobID string) ([]domain.Execution, error)
	FindAll(ctx context.Context, limit int, offset int, jobID string) ([]domain.Execution, int64, error)
	Update(ctx context.Context, execution domain.Execution) error
}
