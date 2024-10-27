package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"book-tracker/domain"

	"github.com/golang-jwt/jwt"
	"google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
)

//go:generate mockery --name UserRepository
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, u domain.User) (domain.User, error)
}

type AuthService struct {
	userRepo  UserRepository
	oauth2Svc *oauth2.Service
	jwtSecret string
}

func NewAuthService(ur UserRepository, jwtSecret string) (*AuthService, error) {
	ctx := context.Background()
	oauthSvc, err := oauth2.NewService(ctx, option.WithoutAuthentication())
	if err != nil {
		return nil, err
	}

	return &AuthService{
		userRepo:  ur,
		oauth2Svc: oauthSvc,
		jwtSecret: jwtSecret,
	}, nil
}

// Login verifies a Google OAuth token and creates a user if not exists.
// Returns JWT token or an error if the login fails.
func (a *AuthService) Login(ctx context.Context, token string) (string, error) {
	email, err := a.validateToken(token)
	if err != nil {
		return "", err
	}

	user, err := a.getOrCreateUser(ctx, email)
	if err != nil {
		return "", err
	}

	jwtToken, err := a.generateJWTToken(user.ID)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// validateToken verifies the OAuth token and returns the user's email.
func (a *AuthService) validateToken(token string) (string, error) {
	tokenInfoCall := a.oauth2Svc.Tokeninfo()
	tokenInfoCall.IdToken(token)

	tokenInfo, err := tokenInfoCall.Do()
	if err != nil || tokenInfo.Email == "" {
		return "", fmt.Errorf("failed to validate OAuth token: %w", err)
	}

	return tokenInfo.Email, nil
}

func (a *AuthService) getOrCreateUser(ctx context.Context, email string) (domain.User, error) {
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

func (a *AuthService) generateJWTToken(userID int64) (string, error) {
	claims := &jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(72 * time.Hour).Unix(), // Token expiry time: 3 days
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(a.jwtSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
