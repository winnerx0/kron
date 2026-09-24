package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/winnerx0/kron/internal/config"
	refreshtoken "github.com/winnerx0/kron/internal/refresh_token"
	"github.com/winnerx0/kron/internal/user"
)

const (
	AccessTokenTTL  = 15 * time.Minute
	RefreshTokenTTL = 30 * 24 * time.Hour
)

type TokenResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Service struct {
	config           *config.Config
	userRepo         user.Repository
	refreshTokenRepo refreshtoken.Repository
}

func NewService(config *config.Config, userRepo user.Repository, refreshTokenRepo refreshtoken.Repository) *Service {
	return &Service{config: config, userRepo: userRepo, refreshTokenRepo: refreshTokenRepo}
}

func (s *Service) IssueTokens(ctx context.Context, u user.User) (*TokenResult, error) {
	accessClaims := struct {
		Role string
		jwt.RegisteredClaims
	}{
		Role: "User",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.ID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS512, accessClaims)
	accessToken, err := accessJWT.SignedString([]byte(s.config.JwtAccessTokenSecret))
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   u.ID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS512, refreshClaims)
	refreshToken, err := refreshJWT.SignedString([]byte(s.config.JwtRefreshTokenSecret))
	if err != nil {
		return nil, err
	}

	if err := s.refreshTokenRepo.Save(ctx, refreshtoken.RefreshToken{
		ID:        uuid.NewString(),
		UserID:    u.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(RefreshTokenTTL),
	}); err != nil {
		return nil, err
	}

	return &TokenResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Refresh(ctx context.Context, refreshToken string) (*TokenResult, error) {
	stored, err := s.refreshTokenRepo.FindByToken(ctx, refreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if time.Now().After(stored.ExpiresAt) {
		_ = s.refreshTokenRepo.DeleteByToken(ctx, refreshToken)
		return nil, ErrRefreshTokenExpired
	}

	parsed, err := jwt.ParseWithClaims(refreshToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.config.JwtRefreshTokenSecret), nil
	})
	if err != nil || !parsed.Valid {
		return nil, ErrInvalidRefreshToken
	}

	u, err := s.userRepo.FindById(ctx, stored.UserID)
	if err != nil || u == nil {
		return nil, ErrUserNotFound
	}

	if err := s.refreshTokenRepo.DeleteByToken(ctx, refreshToken); err != nil {
		return nil, err
	}

	return s.IssueTokens(ctx, *u)
}
