package execution

import "context"

type Service struct {
	repo Repository
}

func NewExecutionService(repo Repository) Service {
	return Service{repo: repo}
}

func (s *Service) FindAll(ctx context.Context) ([]Execution, error) {
	return s.repo.FindAll(ctx)
}

func (s *Service) FindByJobID(ctx context.Context, jobID string) ([]Execution, error) {
	return s.repo.FindByJobID(ctx, jobID)
}

func (s *Service) Create(ctx context.Context, execution Execution) error {
	return s.repo.Save(ctx, execution)
}