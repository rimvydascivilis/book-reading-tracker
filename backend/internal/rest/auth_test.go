package rest_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rimvydascivilis/book-tracker/backend/internal/rest"
	"github.com/rimvydascivilis/book-tracker/backend/mocks"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_Login_Success(t *testing.T) {
	e := echo.New()

	mockAuthService := new(mocks.AuthService)
	handler := rest.NewAuthHandler(mockAuthService)

	requestBody, _ := json.Marshal(map[string]string{"token": "valid_token"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockAuthService.On("Login", mock.Anything, "valid_token").Return("jwt_token", nil)

	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"token":"jwt_token"}`, rec.Body.String())
	}

	mockAuthService.AssertExpectations(t)
}

func TestAuthHandler_Login_BindError(t *testing.T) {
	e := echo.New()
	handler := rest.NewAuthHandler(new(mocks.AuthService))

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer([]byte("{invalid_json}")))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"message":"invalid request format"}`, rec.Body.String())
	}
}

func TestAuthHandler_Login_MissingToken(t *testing.T) {
	e := echo.New()
	handler := rest.NewAuthHandler(new(mocks.AuthService))

	requestBody, _ := json.Marshal(map[string]string{})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, `{"message":"missing token"}`, rec.Body.String())
	}
}

func TestAuthHandler_Login_Failure(t *testing.T) {
	e := echo.New()

	mockAuthService := new(mocks.AuthService)
	handler := rest.NewAuthHandler(mockAuthService)

	requestBody, _ := json.Marshal(map[string]string{"token": "invalid_token"})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockAuthService.On("Login", mock.Anything, "invalid_token").Return("", assert.AnError)

	if assert.NoError(t, handler.Login(c)) {
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, `{"message":"failed to login"}`, rec.Body.String()) // Updated expected message
	}

	mockAuthService.AssertExpectations(t)
}
