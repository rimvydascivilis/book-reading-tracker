package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type JWTService struct {
	userRepo domain.UserRepository
	secret   string
}

func NewJWTService(secret string, ur domain.UserRepository) *JWTService {
	return &JWTService{
		userRepo: ur,
		secret:   secret,
	}
}

func (j *JWTService) GenerateToken(ctx context.Context, userID int64) (string, error) {
	user, err := j.userRepo.GetByID(ctx, userID)
	if err != nil {
		return "", err
	}

	claims := &jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(72 * time.Hour).Unix(), // Token expiry time: 3 days
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
