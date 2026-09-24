package refreshtoken

import "context"

type Repository interface {
	Save(ctx context.Context, token RefreshToken) error
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)
	DeleteByToken(ctx context.Context, token string) error
}
