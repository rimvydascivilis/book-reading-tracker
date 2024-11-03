package user

import (
	"context"
	"errors"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type UserService struct {
	userRepo domain.UserRepository
}

func NewUserService(ur domain.UserRepository) *UserService {
	return &UserService{
		userRepo: ur,
	}
}

func (a *UserService) GetOrCreateUser(ctx context.Context, email string) (domain.User, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			user = domain.User{
				Email: email,
			}
			return a.userRepo.CreateUser(ctx, user)
		}
		return domain.User{}, err // propagate other errors
	}

	return user, nil
}
