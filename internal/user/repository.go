package user

import "context"

type Repository interface {
	Save(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindById(ctx context.Context, id string) (*User, error)
}