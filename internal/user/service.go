package user

import (
	"context"

	"github.com/google/uuid"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (s *Service) CreateUser(ctx context.Context, createUserRequest CreateUserRequest) error {
	user := User{
		ID:             uuid.NewString(),
		Username:       createUserRequest.Name,
		Email:          createUserRequest.Email,
		ProfilePicture: createUserRequest.ProfilePicture,
	}
	return s.repo.Save(ctx, user)
}

func (s *Service) FindUserById(ctx context.Context, id string) (*User, error) {
	return s.repo.FindById(ctx, id)
}
