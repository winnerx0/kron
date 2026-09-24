package http

import (
	"encoding/json"
	"net/http"

	"github.com/winnerx0/kron/internal/response"
	"github.com/winnerx0/kron/internal/user"
)

type UserHandler struct {
	service user.Service
}

func NewUserHandler(service user.Service) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := userIDFromContext(r)
	if !ok {
		response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	currentUser, err := h.service.FindUserById(r.Context(), userID)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if currentUser == nil {
		response.WriteError(w, http.StatusNotFound, "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(user.UserResponse{
		ID:             currentUser.ID,
		Name:           currentUser.Username,
		Email:          currentUser.Email,
		ProfilePicture: currentUser.ProfilePicture,
		JoinedAt:       currentUser.JoinedAt,
	})
}
