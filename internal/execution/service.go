package execution

import (
	"context"

	"github.com/winnerx0/kron/internal/domain"
)

type Service interface {
	FindAll(ctx context.Context, page int, pageSize int) (PaginatedExecutionsResponse, error)
	FindByJobID(ctx context.Context, jobID string) ([]domain.Execution, error)
	Create(ctx context.Context, execution domain.Execution) error
}

type ExecutionService struct {
	repo Repository
}

func NewExecutionService(repo Repository) *ExecutionService {
	return &ExecutionService{repo: repo}
}

func (s *ExecutionService) FindAll(ctx context.Context, page int, pageSize int) (PaginatedExecutionsResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	executions, total, err := s.repo.FindAll(ctx, pageSize, offset)
	if err != nil {
		return PaginatedExecutionsResponse{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(pageSize) - 1) / int64(pageSize))
	}

	return PaginatedExecutionsResponse{
		Items:      executions,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *ExecutionService) FindByJobID(ctx context.Context, jobID string) ([]domain.Execution, error) {
	return s.repo.FindByJobID(ctx, jobID)
}

func (s *ExecutionService) Create(ctx context.Context, execution domain.Execution) error {
	return s.repo.Save(ctx, execution)
}
