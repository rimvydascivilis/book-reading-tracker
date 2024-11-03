package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/mocks"
	"github.com/stretchr/testify/assert"
)

func setupAuthService() (*AuthService, *mocks.UserService, *mocks.OAuth2Service, *mocks.TokenService) {
	userSvc := new(mocks.UserService)
	oauth2Svc := new(mocks.OAuth2Service)
	tokenService := new(mocks.TokenService)

	authService := NewAuthService(userSvc, oauth2Svc, tokenService)

	return authService, userSvc, oauth2Svc, tokenService
}

func TestAuthService_Login_Success(t *testing.T) {
	authService, userSvc, oauth2Svc, tokenService := setupAuthService()

	ctx := context.Background()
	testToken := "valid_oauth_token"
	testEmail := "test@example.com"
	testUserID := int64(123)
	expectedJWT := "generated_jwt_token"

	oauth2Svc.On("ValidateToken", testToken).Return(testEmail, nil)
	userSvc.On("GetOrCreateUser", ctx, testEmail).Return(domain.User{ID: testUserID, Email: testEmail}, nil)
	tokenService.On("GenerateToken", ctx, testUserID).Return(expectedJWT, nil)

	token, err := authService.Login(ctx, testToken)

	assert.NoError(t, err)
	assert.Equal(t, expectedJWT, token)

	oauth2Svc.AssertExpectations(t)
	userSvc.AssertExpectations(t)
	tokenService.AssertExpectations(t)
}

func TestAuthService_Login_OAuth2ValidationError(t *testing.T) {
	authService, userSvc, oauth2Svc, tokenService := setupAuthService()

	ctx := context.Background()
	testToken := "invalid_oauth_token"

	oauth2Svc.On("ValidateToken", testToken).Return("", errors.New("invalid token"))

	token, err := authService.Login(ctx, testToken)

	assert.Error(t, err)
	assert.Empty(t, token)

	oauth2Svc.AssertExpectations(t)
	userSvc.AssertExpectations(t)
	tokenService.AssertExpectations(t)
}
