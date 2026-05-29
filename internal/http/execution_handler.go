package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/winnerx0/kron/internal/execution"
	"github.com/winnerx0/kron/internal/response"
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
// @Router /api/execution/all [get]
func (h *ExecutionHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	page := parsePositiveInt(r.URL.Query().Get("page"), 1)
	pageSize := parsePositiveInt(r.URL.Query().Get("pageSize"), 10)
	jobID := r.URL.Query().Get("jobID")

	executions, err := h.service.FindAll(r.Context(), page, pageSize, jobID)

	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(executions)
}

// @Summary Find an execution by ID
// @Description Get the details of a single execution including its response body
// @Tags executions
// @Produce json
// @Param id path string true "Execution ID"
// @Success 200 {object} execution.ExecutionDetailDTO
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/execution/{id} [get]
func (h *ExecutionHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	detail, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, execution.ErrNotFound) {
			response.WriteError(w, http.StatusNotFound, "Execution not found")
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(detail)
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
