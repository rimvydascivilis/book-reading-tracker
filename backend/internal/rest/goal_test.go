package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/internal/rest"
	"github.com/rimvydascivilis/book-tracker/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetGoal_Success(t *testing.T) {
	mockSvc := new(mocks.GoalService)
	handler := rest.NewGoalHandler(mockSvc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/goal", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &mockJWTToken)

	expectedGoal := domain.Goal{
		UserID:    1,
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}
	mockSvc.On("GetGoal", mock.Anything, int64(1)).Return(expectedGoal, nil)

	err := handler.GetGoal(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseGoal domain.Goal
	err = json.Unmarshal(rec.Body.Bytes(), &responseGoal)
	assert.NoError(t, err)
	assert.Equal(t, expectedGoal, responseGoal)
}

func TestGetGoal_InvalidToken(t *testing.T) {
	mockSvc := new(mocks.GoalService)
	handler := rest.NewGoalHandler(mockSvc)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/goal", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetGoal(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid token")
}

func TestSetGoal_CreateGoal_Success(t *testing.T) {
	mockSvc := new(mocks.GoalService)
	handler := rest.NewGoalHandler(mockSvc)

	e := echo.New()
	reqGoal := domain.Goal{
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}
	body, _ := json.Marshal(reqGoal)

	req := httptest.NewRequest(http.MethodPost, "/goal", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &mockJWTToken)

	expectedGoal := domain.Goal{
		UserID:    1,
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}
	mockSvc.On("SetGoal", mock.Anything, int64(1), reqGoal).Return(expectedGoal, nil)

	err := handler.SetGoal(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseGoal domain.Goal
	err = json.Unmarshal(rec.Body.Bytes(), &responseGoal)
	assert.NoError(t, err)
	assert.Equal(t, expectedGoal, responseGoal)
}

func TestSetGoal_InvalidToken(t *testing.T) {
	mockSvc := new(mocks.GoalService)
	handler := rest.NewGoalHandler(mockSvc)

	e := echo.New()
	reqGoal := domain.Goal{
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}
	body, _ := json.Marshal(reqGoal)

	req := httptest.NewRequest(http.MethodPost, "/goal", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.SetGoal(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.Contains(t, rec.Body.String(), "Invalid token")
}

func TestSetGoal_ServiceError(t *testing.T) {
	mockSvc := new(mocks.GoalService)
	handler := rest.NewGoalHandler(mockSvc)

	e := echo.New()
	reqGoal := domain.Goal{
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}
	body, _ := json.Marshal(reqGoal)

	req := httptest.NewRequest(http.MethodPost, "/goal", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", &mockJWTToken)

	mockSvc.On("SetGoal", mock.Anything, int64(1), reqGoal).Return(domain.Goal{}, fmt.Errorf("service error"))

	err := handler.SetGoal(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "server error")
}
