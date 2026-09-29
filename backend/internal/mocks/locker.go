package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockLocker struct {
	mock.Mock
}

func (m *MockLocker) tickKey(jobId string, scheduleAt time.Time) string {
	args := m.Called(jobId, scheduleAt)

	return args.String(0)
}

func (m *MockLocker) Claim(ctx context.Context, jobID string, scheduleAt time.Time) (bool, error) {

	args := m.Called(ctx, jobID, scheduleAt)

	return args.Bool(0), args.Error(1)
}

func (m *MockLocker) Release(ctx context.Context, jobID string, scheduleAt time.Time) error {

	args := m.Called(ctx, jobID, scheduleAt)

	return args.Error(0)
}
