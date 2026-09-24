package dashboard

import "context"

type Counts struct {
	TotalJobs            int64
	EnabledJobs          int64
	TotalExecutions      int64
	SuccessfulExecutions int64
	FailedExecutions     int64
}

type Repository interface {
	Counts(ctx context.Context, userID string) (Counts, error)
	RecentExecutions(ctx context.Context, userID string, limit int) ([]RecentExecution, error)
	Jobs(ctx context.Context, userID string, limit int) ([]JobSummary, error)
}
