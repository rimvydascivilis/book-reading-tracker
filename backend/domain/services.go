package domain

import "context"

type AuthService interface {
	Login(ctx context.Context, googleOauthToken string) (string, error)
}

type TokenService interface {
	GenerateToken(ctx context.Context, userID int64) (string, error)
}

type OAuth2Service interface {
	ValidateToken(token string) (string, error)
}
