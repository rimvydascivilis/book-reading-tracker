package user

import (
	"context"
	"errors"
	"testing"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserService() (*UserService, *mocks.UserRepository) {
	userRepo := new(mocks.UserRepository)
	userService := NewUserService(userRepo)

	return userService, userRepo
}

func TestUserService_GetOrCreateUser_Success(t *testing.T) {
	userService, userRepo := setupUserService()

	ctx := context.Background()
	testEmail := "existinguser@example.com"
	testUser := domain.User{ID: 1, Email: testEmail}

	userRepo.On("GetByEmail", ctx, testEmail).Return(testUser, nil)

	user, err := userService.GetOrCreateUser(ctx, testEmail)

	assert.NoError(t, err)
	assert.Equal(t, testUser, user)

	userRepo.AssertExpectations(t)
}

func TestUserService_GetOrCreateUser_CreateUser(t *testing.T) {
	userService, userRepo := setupUserService()

	ctx := context.Background()
	testEmail := "newuser@example.com"
	testUser := domain.User{ID: 2, Email: testEmail}

	userRepo.On("GetByEmail", ctx, testEmail).Return(domain.User{}, domain.ErrUserNotFound)

	userRepo.On("CreateUser", ctx, mock.AnythingOfType("domain.User")).Return(testUser, nil)

	user, err := userService.GetOrCreateUser(ctx, testEmail)

	assert.NoError(t, err)
	assert.Equal(t, testUser, user)

	userRepo.AssertExpectations(t)
}

func TestUserService_GetOrCreateUser_Error(t *testing.T) {
	userService, userRepo := setupUserService()

	ctx := context.Background()
	testEmail := "erroruser@example.com"

	userRepo.On("GetByEmail", ctx, testEmail).Return(domain.User{}, errors.New("some error"))

	user, err := userService.GetOrCreateUser(ctx, testEmail)

	assert.Error(t, err)
	assert.EqualError(t, err, "some error")
	assert.Equal(t, domain.User{}, user)

	userRepo.AssertExpectations(t)
}
