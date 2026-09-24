package dashboard

import "time"

type JobSummary struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Enabled  bool   `json:"enabled"`
}

type RecentExecution struct {
	ID         string    `json:"id"`
	JobID      string    `json:"jobID"`
	JobName    string    `json:"jobName"`
	Status     string    `json:"status"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt"`
}

type SummaryResponse struct {
	TotalJobs            int64             `json:"totalJobs"`
	EnabledJobs          int64             `json:"enabledJobs"`
	DisabledJobs         int64             `json:"disabledJobs"`
	TotalExecutions      int64             `json:"totalExecutions"`
	SuccessfulExecutions int64             `json:"successfulExecutions"`
	FailedExecutions     int64             `json:"failedExecutions"`
	SuccessRate          int               `json:"successRate"`
	RecentExecutions     []RecentExecution `json:"recentExecutions"`
	Jobs                 []JobSummary      `json:"jobs"`
}
