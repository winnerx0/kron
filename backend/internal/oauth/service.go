package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/winnerx0/kron/internal/auth"
	"github.com/winnerx0/kron/internal/config"
	"github.com/winnerx0/kron/internal/user"
)

const (
	googleAuthURL     = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL    = "https://oauth2.googleapis.com/token"
	googleUserInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
)

type TokenResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserInfo struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Picture string `json:"picture"`
}

type Service struct {
	config      *config.Config
	userRepo    user.Repository
	authService *auth.Service
}

func NewService(config *config.Config, userRepo user.Repository, authService *auth.Service) *Service {
	return &Service{config: config, userRepo: userRepo, authService: authService}
}

func (s *Service) GoogleAuthURL(state string) string {
	q := url.Values{}
	q.Set("client_id", s.config.GoogleClientID)
	q.Set("redirect_uri", s.config.GoogleRedirectURI)
	q.Set("response_type", "code")
	q.Set("access_type", "online")
	q.Set("scope", "email profile openid")
	q.Set("state", state)
	q.Set("prompt", "select_account")

	return fmt.Sprintf("%s?%s", googleAuthURL, q.Encode())
}

func (s *Service) ExchangeCode(ctx context.Context, code string) (*TokenResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, googleTokenURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("client_id", s.config.GoogleClientID)
	q.Add("client_secret", s.config.GoogleClientSecret)
	q.Add("redirect_uri", s.config.GoogleRedirectURI)
	q.Add("code", code)
	q.Add("grant_type", "authorization_code")
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, TokenExchangeError{StatusCode: resp.StatusCode}
	}

	var result TokenResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *Service) SaveUser(ctx context.Context, accessToken string) (*auth.TokenResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, googleUserInfoURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, UserInfoError{StatusCode: resp.StatusCode}
	}

	var userInfo UserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	u := user.User{
		ID:             uuid.NewString(),
		Username:       userInfo.Name,
		Email:          userInfo.Email,
		ProfilePicture: userInfo.Picture,
	}

	user, err := s.userRepo.FindByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	if user != nil {
		return s.authService.IssueTokens(ctx, *user)
	}

	if err := s.userRepo.Save(ctx, u); err != nil {
		return nil, err
	}

	return s.authService.IssueTokens(ctx, u)
}

func RandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}
