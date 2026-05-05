package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/winnerx0/kron/internal/config"
	"github.com/winnerx0/kron/internal/response"
	"github.com/winnerx0/kron/internal/user"
)

type AuthMiddleware struct {
	config   *config.Config
	userRepo user.Repository
}

func NewAuthMiddleware(config *config.Config, userRepo user.Repository) *AuthMiddleware {
	return &AuthMiddleware{
		config:   config,
		userRepo: userRepo,
	}
}

func (m *AuthMiddleware) Use(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		token := strings.SplitAfter(authHeader, "Bearer ")[1]

		jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			if exp, err := t.Claims.GetExpirationTime(); err != nil || exp.Before(time.Now()) {
				return nil, jwt.ErrTokenExpired
			}

			return []byte(m.config.JwtAccessTokenSecret), nil
		})

		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		subject, err := jwtToken.Claims.GetSubject()

		if err != nil {
			response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := m.userRepo.FindById(r.Context(), subject)

		if user == nil || err != nil {
			response.WriteError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), "userId", user.ID)

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
