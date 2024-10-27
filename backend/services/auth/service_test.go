package auth

import (
	"book-tracker/domain"
	"book-tracker/mocks"
	"context"
	"errors"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Login_Success(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	oauth2Svc := new(mocks.OAuth2Service)
	jwtSecret := "test_secret"

	authService := NewAuthService(userRepo, oauth2Svc, jwtSecret)

	ctx := context.Background()
	testToken := "valid_oauth_token"
	testEmail := "test@example.com"
	testUserID := int64(123)

	oauth2Svc.On("ValidateToken", testToken).Return(testEmail, nil)

	userRepo.On("GetByEmail", ctx, testEmail).Return(domain.User{ID: testUserID, Email: testEmail}, nil)

	token, err := authService.Login(ctx, testToken)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(testUserID), claims["id"])
}

func TestAuthService_Login_OAuth2ValidationError(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	oauth2Svc := new(mocks.OAuth2Service)
	jwtSecret := "test_secret"

	authService := NewAuthService(userRepo, oauth2Svc, jwtSecret)

	ctx := context.Background()
	testToken := "invalid_oauth_token"

	oauth2Svc.On("ValidateToken", testToken).Return("", errors.New("invalid token"))

	token, err := authService.Login(ctx, testToken)

	assert.Error(t, err)
	assert.Empty(t, token)
}

func TestAuthService_Login_CreateUser(t *testing.T) {
	userRepo := new(mocks.UserRepository)
	oauth2Svc := new(mocks.OAuth2Service)
	jwtSecret := "test_secret"

	authService := NewAuthService(userRepo, oauth2Svc, jwtSecret)

	ctx := context.Background()
	testToken := "valid_oauth_token"
	testEmail := "newuser@example.com"
	testUserID := int64(456)

	oauth2Svc.On("ValidateToken", testToken).Return(testEmail, nil)

	userRepo.On("GetByEmail", ctx, testEmail).Return(domain.User{}, domain.ErrUserNotFound)

	userRepo.On("CreateUser", ctx, mock.AnythingOfType("domain.User")).Return(domain.User{ID: testUserID, Email: testEmail}, nil)

	token, err := authService.Login(ctx, testToken)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(testUserID), claims["id"])
}
