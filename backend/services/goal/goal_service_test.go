package goal

import (
	"context"
	"fmt"
	"testing"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupGoalService() (domain.GoalService, *mocks.GoalRepository, *mocks.ValidationService) {
	goalRepo := new(mocks.GoalRepository)
	validationSvc := new(mocks.ValidationService)
	goalService := NewGoalService(goalRepo, nil, nil, validationSvc)

	return goalService, goalRepo, validationSvc
}

func TestGetGoal_Success(t *testing.T) {
	service, mockRepo, _ := setupGoalService()

	userID := int64(1)
	expectedGoal := domain.Goal{
		UserID:    userID,
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}

	mockRepo.On("GetGoalByUserID", mock.Anything, userID).Return(expectedGoal, nil)

	result, err := service.GetGoal(context.Background(), userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedGoal, result)
	mockRepo.AssertExpectations(t)
}

func TestGetGoal_NotFound(t *testing.T) {
	service, mockRepo, _ := setupGoalService()

	userID := int64(1)

	mockRepo.On("GetGoalByUserID", mock.Anything, userID).Return(domain.Goal{}, domain.ErrRecordNotFound)

	result, err := service.GetGoal(context.Background(), userID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrRecordNotFound, err)
	assert.Equal(t, domain.Goal{}, result)
	mockRepo.AssertExpectations(t)
}

func TestSetGoal_UpdateExistingGoal(t *testing.T) {
	service, mockRepo, validationSvc := setupGoalService()

	userID := int64(1)
	currentGoal := domain.Goal{
		UserID:    userID,
		Type:      "pages",
		Frequency: "weekly",
		Value:     5,
	}
	updatedGoal := domain.Goal{
		Frequency: "daily",
		Value:     20,
	}
	expectedGoal := domain.Goal{
		UserID:    userID,
		Type:      "pages",
		Frequency: "daily",
		Value:     20,
	}

	mockRepo.On("GetGoalByUserID", mock.Anything, userID).Return(currentGoal, nil)
	validationSvc.On("ValidateStruct", expectedGoal).Return(nil)
	mockRepo.On("UpdateGoal", mock.Anything, expectedGoal).Return(expectedGoal, nil)

	result, err := service.SetGoal(context.Background(), userID, updatedGoal)

	assert.NoError(t, err)
	assert.Equal(t, expectedGoal, result)
	mockRepo.AssertExpectations(t)
	validationSvc.AssertExpectations(t)
}

func TestSetGoal_CreateNewGoal(t *testing.T) {
	service, mockRepo, validationSvc := setupGoalService()

	userID := int64(1)
	newGoal := domain.Goal{
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}
	expectedGoal := domain.Goal{
		UserID:    userID,
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}

	mockRepo.On("GetGoalByUserID", mock.Anything, userID).Return(domain.Goal{}, domain.ErrRecordNotFound)
	validationSvc.On("ValidateStruct", expectedGoal).Return(nil)
	mockRepo.On("CreateGoal", mock.Anything, expectedGoal).Return(expectedGoal, nil)

	result, err := service.SetGoal(context.Background(), userID, newGoal)

	assert.NoError(t, err)
	assert.Equal(t, expectedGoal, result)
	mockRepo.AssertExpectations(t)
	validationSvc.AssertExpectations(t)
}

func TestSetGoal_ValidationError(t *testing.T) {
	service, mockRepo, validationSvc := setupGoalService()

	userID := int64(1)
	newGoal := domain.Goal{
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}
	expectedGoal := domain.Goal{
		UserID:    userID,
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}

	mockRepo.On("GetGoalByUserID", mock.Anything, userID).Return(domain.Goal{}, domain.ErrRecordNotFound)
	validationSvc.On("ValidateStruct", expectedGoal).Return(fmt.Errorf("validation error"))

	result, err := service.SetGoal(context.Background(), userID, newGoal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation error")
	assert.Equal(t, domain.Goal{}, result)
	mockRepo.AssertExpectations(t)
	validationSvc.AssertExpectations(t)
}

func TestSetGoal_UpdateError(t *testing.T) {
	service, mockRepo, validationSvc := setupGoalService()

	userID := int64(1)
	currentGoal := domain.Goal{
		UserID:    userID,
		Type:      "pages",
		Frequency: "weekly",
		Value:     5,
	}
	updatedGoal := domain.Goal{
		Type:      "books",
		Frequency: "daily",
		Value:     20,
	}
	expectedGoal := domain.Goal{
		UserID:    userID,
		Type:      "books",
		Frequency: "daily",
		Value:     20,
	}

	mockRepo.On("GetGoalByUserID", mock.Anything, userID).Return(currentGoal, nil)
	validationSvc.On("ValidateStruct", expectedGoal).Return(nil)
	mockRepo.On("UpdateGoal", mock.Anything, expectedGoal).Return(domain.Goal{}, fmt.Errorf("update error"))

	result, err := service.SetGoal(context.Background(), userID, updatedGoal)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "update error")
	assert.Equal(t, domain.Goal{}, result)
	mockRepo.AssertExpectations(t)
	validationSvc.AssertExpectations(t)
}
