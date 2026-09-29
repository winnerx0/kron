package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/winnerx0/kron/internal/domain"
)

type MockPublisher struct {
	mock.Mock
}

func (m *MockPublisher) Publish(ctx context.Context, job domain.Job) error {

	args := m.Called(ctx, job)

	return args.Error(0)
}
