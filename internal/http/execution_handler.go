package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/winnerx0/kron/internal/execution"
)

type ExecutionHandler struct {
	service execution.Service
}

func NewExecutionHandler(service execution.Service) *ExecutionHandler {
	return &ExecutionHandler{
		service: service,
	}
}

// @Summary Find all executions
// @Description Get a list of all executions
// @Tags executions
// @Produce json
// @Success 200 {array} execution.Execution
// @Failure 500 {object} map[string]string
// @Router /execution/all [get]
func (h *ExecutionHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	page := parsePositiveInt(r.URL.Query().Get("page"), 1)
	pageSize := parsePositiveInt(r.URL.Query().Get("pageSize"), 10)

	executions, err := h.service.FindAll(r.Context(), page, pageSize)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(executions)
}

func parsePositiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed < 1 {
		return fallback
	}

	return parsed
}
