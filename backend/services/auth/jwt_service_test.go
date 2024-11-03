package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/mocks"
	"github.com/stretchr/testify/assert"
)

func TestJWTService_GenerateToken_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	secret := "my_secret"
	jwtService := NewJWTService(secret, userRepo)

	ctx := context.Background()
	userID := int64(123)

	userRepo.On("GetByID", ctx, userID).Return(domain.User{ID: userID}, nil)

	token, err := jwtService.GenerateToken(ctx, userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims := &jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	assert.NoError(t, err)
	assert.Equal(t, userID, int64((*claims)["id"].(float64)))
}

func TestJWTService_GenerateToken_UserNotFound(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	secret := "my_secret"
	jwtService := NewJWTService(secret, userRepo)

	ctx := context.Background()
	userID := int64(456)

	userRepo.On("GetByID", ctx, userID).Return(domain.User{}, errors.New("user not found"))

	token, err := jwtService.GenerateToken(ctx, userID)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "user not found", err.Error())
}

func TestJWTService_GenerateToken_InvalidUserID(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	secret := "my_secret"
	jwtService := NewJWTService(secret, userRepo)

	ctx := context.Background()
	invalidUserID := int64(999) // An ID that does not exist

	userRepo.On("GetByID", ctx, invalidUserID).Return(domain.User{}, domain.ErrUserNotFound)

	token, err := jwtService.GenerateToken(ctx, invalidUserID)

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, domain.ErrUserNotFound, err)
}
