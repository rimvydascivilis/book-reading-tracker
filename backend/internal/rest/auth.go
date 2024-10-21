package rest

import (
	"context"

	"github.com/labstack/echo/v4"

	"book-tracker/domain"
)

type ResponseError struct {
	Message string `json:"message"`
}

//go:generate mockery --name AuthService
type AuthService interface {
	Login(ctx context.Context, token string) (domain.User, error)
	Logout(ctx context.Context, id string) error
}

type AuthHandler struct {
	AuthSvc AuthService
}

func NewAuthHandler(e *echo.Echo, as AuthService) {
	handler := &AuthHandler{
		AuthSvc: as,
	}
	e.POST("/login", handler.Login)
	e.POST("/logout", handler.Logout)
}

func (a *AuthHandler) Login(c echo.Context) error {
	return nil
}

func (a *AuthHandler) Logout(c echo.Context) error {
	return nil
}
