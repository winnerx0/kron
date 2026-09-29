package scheduler

import (
	"context"
	"time"

	"github.com/winnerx0/kron/internal/domain"
)

type JobLocker interface {
	Claim(ctx context.Context, jobID string, scheduleAt time.Time) (bool, error)

	Release(ctx context.Context, jobID string, scheduleAt time.Time) error
}

type JobPoller interface {
	Run(ctx context.Context)

	tick(ctx context.Context)
}

type JobPublisher interface {
	Publish(ctx context.Context, job domain.Job) error
}
