package job

import (
	"context"

	"github.com/winnerx0/kron/internal/domain"
)

type Repository interface {
	FindAll(ctx context.Context, userID string) ([]domain.Job, error)

	FindAllNextRun(ctx context.Context) ([]domain.Job, error)

	FindByID(ctx context.Context, id string) (domain.Job, error)

	Create(ctx context.Context, job domain.Job) (domain.Job, error)

	Update(ctx context.Context, job domain.Job) (domain.Job, error)

	Delete(ctx context.Context, id string) error
}
