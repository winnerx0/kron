package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/winnerx0/kron/internal/domain"
	"github.com/winnerx0/kron/internal/execution"
)

type MockExecutionRepository struct {
	mock.Mock
}

func (m *MockExecutionRepository) Save(ctx context.Context, value domain.Execution) error {
	args := m.Called(ctx, value)
	return args.Error(0)
}

func (m *MockExecutionRepository) FindByJobID(ctx context.Context, jobID string) ([]domain.Execution, error) {
	args := m.Called(ctx, jobID)
	return args.Get(0).([]domain.Execution), args.Error(1)
}

func (m *MockExecutionRepository) FindAll(ctx context.Context, limit int, offset int, jobID string) ([]execution.ExecutionDTO, int64, error) {
	args := m.Called(ctx, limit, offset, jobID)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]execution.ExecutionDTO), args.Get(1).(int64), args.Error(2)
}

func (m *MockExecutionRepository) FindByID(ctx context.Context, id string) (execution.ExecutionDetailDTO, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(execution.ExecutionDetailDTO), args.Error(1)
}

func (m *MockExecutionRepository) Update(ctx context.Context, value domain.Execution) error {
	args := m.Called(ctx, value)
	return args.Error(0)
}
