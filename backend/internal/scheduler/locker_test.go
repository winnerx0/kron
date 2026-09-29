package scheduler

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/winnerx0/kron/internal/domain"
	"github.com/winnerx0/kron/internal/job"
	"github.com/winnerx0/kron/internal/mocks"
)

type MockPoller struct {
	mock.Mock
	jobRepo   job.Repository
	publisher JobPublisher
	locker    JobLocker
	interval  time.Duration
}

func TestLockerClaim_CallsSetNX(t *testing.T) {

	jobRepo := new(mocks.MockJobRepository)

	jobRepo.On("FindAllNextRun", mock.Anything).Return(domain.Job{
		ID:     "job-1",
		Name:   "mock-job",
		Status: true,
	}, nil)

	runAt := time.Unix(100, 0)

	db, mock := redismock.NewClientMock()

	mock.ExpectSetNX("tick:job-1100", "worker-1", time.Minute).SetVal(true)

	locker := NewLocker(db, time.Minute, "worker-1")

	claimed, err := locker.Claim(context.Background(), "job-1", runAt)

	require.NoError(t, err)

	require.True(t, claimed)

	require.NoError(t, mock.ExpectationsWereMet())
}
