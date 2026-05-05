package http

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/winnerx0/kron/internal/oauth"
	"github.com/winnerx0/kron/internal/response"
)

type OauthHandler struct {
	frontendURL string
	service *oauth.Service
}

func NewOauthHandler(frontendURL string, service *oauth.Service) *OauthHandler {
	return &OauthHandler{frontendURL: frontendURL, service: service}
}

func (h *OauthHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := oauth.RandomString(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		MaxAge:   int(5 * time.Minute.Seconds()),
	})

	http.Redirect(w, r, h.service.GoogleAuthURL(state), http.StatusFound)
}

func (h *OauthHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	state, err := r.Cookie("oauth_state")
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid state")
		return
	}

	if r.FormValue("state") != state.Value {
		response.WriteError(w, http.StatusBadRequest, "Invalid state")
		return
	}

	code := r.FormValue("code")
	if code == "" {
		response.WriteError(w, http.StatusBadRequest, "Missing code")
		return
	}

	tokens, err := h.service.ExchangeCode(r.Context(), code)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tokenResult, err := h.service.SaveUser(r.Context(), tokens.AccessToken)
	if err != nil {
		slog.Info("error while saving", "err", err)
		response.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	params := url.Values{}
	params.Set("access", tokenResult.AccessToken)
	params.Set("refresh", tokenResult.RefreshToken)

	http.Redirect(w, r, fmt.Sprintf("%s/callback?%s", h.frontendURL, params.Encode()), http.StatusFound)
}
