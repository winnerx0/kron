package job

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/winnerx0/kron/internal/domain"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) FindAll(ctx context.Context, userID string) ([]domain.Job, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.Job), args.Error(1)
}

func (m *MockRepository) FindAllNextRun(ctx context.Context) ([]domain.Job, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Job), args.Error(1)
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (domain.Job, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Job), args.Error(1)
}

func (m *MockRepository) Create(ctx context.Context, job domain.Job) (domain.Job, error) {
	args := m.Called(ctx, job)
	return args.Get(0).(domain.Job), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, job domain.Job) (domain.Job, error) {
	args := m.Called(ctx, job)
	return args.Get(0).(domain.Job), args.Error(1)
}

func (m *MockRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
