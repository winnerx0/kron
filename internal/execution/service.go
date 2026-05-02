package execution

import "context"

type Service struct {
	repo Repository
}

func NewExecutionService(repo Repository) Service {
	return Service{repo: repo}
}

func (s *Service) FindAll(ctx context.Context, page int, pageSize int) (PaginatedExecutionsResponse, error) {
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

func (s *Service) FindByJobID(ctx context.Context, jobID string) ([]Execution, error) {
	return s.repo.FindByJobID(ctx, jobID)
}

func (s *Service) Create(ctx context.Context, execution Execution) error {
	return s.repo.Save(ctx, execution)
}
