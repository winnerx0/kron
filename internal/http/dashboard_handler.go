package http

import (
	"encoding/json"
	"net/http"

	"github.com/winnerx0/kron/internal/dashboard"
	"github.com/winnerx0/kron/internal/response"
)

type DashboardHandler struct {
	service dashboard.Service
}

func NewDashboardHandler(service dashboard.Service) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) Summary(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	data, err := h.service.Summary(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}
