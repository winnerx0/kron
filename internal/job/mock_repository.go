package job

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (mock *MockRepository) FindAll(ctx context.Context) ([]Job, error) {
	args := mock.Called(ctx)
	return args.Get(0).([]Job), args.Error(1)
}

func (mock *MockRepository) FindAllNextRun(ctx context.Context) ([]Job, error) {
	args := mock.Called(ctx)
	return args.Get(0).([]Job), args.Error(1)
}

func (mock *MockRepository) FindByID(ctx context.Context, id string) (Job, error) {
	args := mock.Called(ctx, id)
	return args.Get(0).(Job), args.Error(1)
}

func (mock *MockRepository) Create(ctx context.Context, job Job) (Job, error) {
	args := mock.Called(ctx, job)
	return args.Get(0).(Job), args.Error(1)
}

func (mock *MockRepository) Update(ctx context.Context, job Job) (Job, error) {
	args := mock.Called(ctx, job)
	return args.Get(0).(Job), args.Error(1)
}

func (mock *MockRepository) Delete(ctx context.Context, id string) error {
	args := mock.Called(ctx, id)
	return args.Error(0)
}
