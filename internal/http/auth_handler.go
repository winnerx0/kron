package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/winnerx0/kron/internal/auth"
	"github.com/winnerx0/kron/internal/response"
)

type AuthHandler struct {
	service *auth.Service
}

func NewAuthHandler(service *auth.Service) *AuthHandler {
	return &AuthHandler{service: service}
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.RefreshToken == "" {
		response.WriteError(w, http.StatusBadRequest, "missing refresh_token")
		return
	}

	tokens, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidRefreshToken) ||
			errors.Is(err, auth.ErrRefreshTokenExpired) ||
			errors.Is(err, auth.ErrUserNotFound) {
			response.WriteError(w, http.StatusUnauthorized, err.Error())
			return
		}
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tokens); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}
}
