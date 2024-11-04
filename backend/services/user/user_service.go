package user

import (
	"context"
	"errors"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type UserService struct {
	userRepo      domain.UserRepository
	validationSvc domain.ValidationService
}

func NewUserService(ur domain.UserRepository, validator domain.ValidationService) domain.UserService {
	return &UserService{
		userRepo:      ur,
		validationSvc: validator,
	}
}

func (a *UserService) GetOrCreateUser(ctx context.Context, email string) (domain.User, error) {
	user, err := a.userRepo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrRecordNotFound) {
			user = domain.User{Email: email}
			if err := a.validationSvc.ValidateStruct(user); err != nil {
				return domain.User{}, err
			}

			return a.userRepo.CreateUser(ctx, user)
		}
		return domain.User{}, err // propagate other errors
	}

	return user, nil
}
