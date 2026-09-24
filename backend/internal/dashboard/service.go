package dashboard

import "context"

const dashboardListLimit = 6

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s Service) Summary(ctx context.Context, userID string) (SummaryResponse, error) {
	counts, err := s.repo.Counts(ctx, userID)
	if err != nil {
		return SummaryResponse{}, err
	}

	recentExecutions, err := s.repo.RecentExecutions(ctx, userID, dashboardListLimit)
	if err != nil {
		return SummaryResponse{}, err
	}

	jobs, err := s.repo.Jobs(ctx, userID, dashboardListLimit)
	if err != nil {
		return SummaryResponse{}, err
	}

	successRate := 0
	if counts.TotalExecutions > 0 {
		successRate = int((counts.SuccessfulExecutions*100 + counts.TotalExecutions/2) / counts.TotalExecutions)
	}

	return SummaryResponse{
		TotalJobs:            counts.TotalJobs,
		EnabledJobs:          counts.EnabledJobs,
		DisabledJobs:         counts.TotalJobs - counts.EnabledJobs,
		TotalExecutions:      counts.TotalExecutions,
		SuccessfulExecutions: counts.SuccessfulExecutions,
		FailedExecutions:     counts.FailedExecutions,
		SuccessRate:          successRate,
		RecentExecutions:     recentExecutions,
		Jobs:                 jobs,
	}, nil
}
