package rest

import (
	"context"
	"net/http"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/utils"

	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

type AuthHandler struct {
	AuthSvc domain.AuthService
}

func NewAuthHandler(as domain.AuthService) *AuthHandler {
	handler := &AuthHandler{
		AuthSvc: as,
	}
	return handler
}

type LoginRequest struct {
	Token string `json:"token"`
}

func (a *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("failed to bind request", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "invalid request format"})
	}

	if req.Token == "" {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "missing token"})
	}

	jwtToken, err := a.AuthSvc.Login(ctx, req.Token)
	if err != nil {
		utils.Error("failed to login", err)
		return c.JSON(http.StatusInternalServerError, ResponseError{Message: err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": jwtToken})
}
