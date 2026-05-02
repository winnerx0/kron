package execution

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) Save(ctx context.Context, execution Execution) error {
	args := mock.Called(ctx, execution)
	return args.Error(0)
}

func (mock *MockRepository) FindByJobID(ctx context.Context, jobID string) ([]Execution, error) {
	args := mock.Called(ctx, jobID)
	return args.Get(0).([]Execution), args.Error(0)
}

func (mock *MockRepository) FindAll(ctx context.Context, limit int, offset int) ([]Execution, int64, error) {
	args := mock.Called(ctx, limit, offset)
	return args.Get(0).([]Execution), args.Get(1).(int64), args.Error(2)
}

func (mock *MockRepository) Update(ctx context.Context, execution Execution) error {
	args := mock.Called(ctx, execution)
	return args.Error(0)
}
