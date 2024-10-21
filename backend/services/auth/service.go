package auth

import (
	"context"

	"book-tracker/domain"
)

//go:generate mockery --name UserRepository
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	Create(ctx context.Context, u *domain.User) error
}

type AuthService struct {
	UserRepo UserRepository
}

func NewAuthService(ur UserRepository) *AuthService {
	return &AuthService{
		UserRepo: ur,
	}
}

func (a *AuthService) Login(ctx context.Context, token string) (domain.User, error) {
	return domain.User{}, nil
}

func (a *AuthService) Logout(ctx context.Context, id string) error {
	return nil
}
