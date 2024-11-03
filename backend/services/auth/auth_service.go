package auth

import (
	"context"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type AuthService struct {
	userSvc   domain.UserService
	oauth2Svc domain.OAuth2Service
	tokenSvc  domain.TokenService
}

func NewAuthService(userSvc domain.UserService, oauthSvc domain.OAuth2Service, tokenSvc domain.TokenService) *AuthService {
	return &AuthService{
		userSvc:   userSvc,
		oauth2Svc: oauthSvc,
		tokenSvc:  tokenSvc,
	}
}

// Login verifies a Google OAuth token and creates a user if not exists.
// Returns JWT token or an error if the login fails.
func (a *AuthService) Login(ctx context.Context, token string) (string, error) {
	email, err := a.oauth2Svc.ValidateToken(token)
	if err != nil {
		return "", err
	}

	user, err := a.userSvc.GetOrCreateUser(ctx, email)
	if err != nil {
		return "", err
	}

	jwtToken, err := a.tokenSvc.GenerateToken(ctx, user.ID)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}
