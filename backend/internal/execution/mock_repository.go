package execution

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/winnerx0/kron/internal/domain"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Save(ctx context.Context, execution domain.Execution) error {
	args := m.Called(ctx, execution)
	return args.Error(0)
}

func (m *MockRepository) FindByJobID(ctx context.Context, jobID string) ([]domain.Execution, error) {
	args := m.Called(ctx, jobID)
	return args.Get(0).([]domain.Execution), args.Error(1)
}

func (m *MockRepository) FindAll(ctx context.Context, limit int, offset int, jobID string) ([]ExecutionDTO, int64, error) {
	args := m.Called(ctx, limit, offset, jobID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]ExecutionDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (ExecutionDetailDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(ExecutionDetailDTO), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, execution domain.Execution) error {
	args := m.Called(ctx, execution)
	return args.Error(0)
}
