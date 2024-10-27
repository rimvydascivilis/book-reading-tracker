package auth

import (
	"context"
	"errors"
	"time"

	"book-tracker/domain"

	"github.com/golang-jwt/jwt"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	CreateUser(ctx context.Context, u domain.User) (domain.User, error)
}

type OAuth2Service interface {
	ValidateToken(token string) (string, error)
}

type AuthService struct {
	userRepo  UserRepository
	oauth2Svc OAuth2Service
	jwtSecret string
}

func NewAuthService(ur UserRepository, oauthSvc OAuth2Service, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  ur,
		oauth2Svc: oauthSvc,
		jwtSecret: jwtSecret,
	}
}

// Login verifies a Google OAuth token and creates a user if not exists.
// Returns JWT token or an error if the login fails.
func (a *AuthService) Login(ctx context.Context, token string) (string, error) {
	email, err := a.oauth2Svc.ValidateToken(token)
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
