package execution

import "github.com/winnerx0/kron/internal/domain"

type PaginatedExecutionsResponse struct {
	Items      []domain.Execution `json:"items"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
	TotalPages int                `json:"totalPages"`
}
