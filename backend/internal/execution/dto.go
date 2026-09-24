package execution

import (
	"time"

	"github.com/winnerx0/kron/internal/domain"
)

type PaginatedExecutionsResponse struct {
	Items      []ExecutionDTO `json:"items"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"pageSize"`
	TotalPages int            `json:"totalPages"`
}

type ExecutionDTO struct {
	ID       string                 `json:"id"`
	JobID    string                 `json:"jobID"`
	JobName  string                 `json:"jobName"`
	Status   domain.ExecutionStatus `json:"status"`
	Started  time.Time              `json:"startedAt"`
	Finished time.Time              `json:"finishedAt"`
}

type ExecutionDetailDTO struct {
	ID           string                 `json:"id"`
	JobID        string                 `json:"jobID"`
	JobName      string                 `json:"jobName"`
	Endpoint     string                 `json:"endpoint"`
	Method       string                 `json:"method"`
	Status       domain.ExecutionStatus `json:"status"`
	Started      time.Time              `json:"startedAt"`
	Finished     time.Time              `json:"finishedAt"`
	ResponseBody string                 `json:"responseBody"`
}
