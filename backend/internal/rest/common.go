package rest

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/utils"
)

type ResponseError struct {
	Message string `json:"message"`
}

func getInt64QueryParam(c echo.Context, key string, defaultValue int64) int64 {
	value := c.QueryParam(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func getUserIDFromToken(c echo.Context) (int64, error) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return 0, errors.New("user token not found")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid token claims")
	}

	userID, ok := claims["id"].(float64)
	if !ok {
		return 0, errors.New("user ID not found in token")
	}

	return int64(userID), nil
}

func getPaginationParams(c echo.Context) (int64, int64) {
	page := getInt64QueryParam(c, "page", 1)
	limit := getInt64QueryParam(c, "limit", 10)
	return page, limit
}

func handleServiceError(c echo.Context, err error) error {
	if errors.Is(err, domain.ErrValidation) {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: err.Error()})
	}
	if errors.Is(err, domain.ErrRecordNotFound) {
		return c.JSON(http.StatusNotFound, ResponseError{Message: err.Error()})
	}
	if errors.Is(err, domain.ErrAuthentication) {
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: err.Error()})
	}
	if errors.Is(err, domain.ErrAlreadyExists) {
		return c.JSON(http.StatusConflict, ResponseError{Message: err.Error()})
	}
	if errors.Is(err, domain.ErrForbidden) {
		return c.JSON(http.StatusForbidden, ResponseError{Message: err.Error()})
	}

	utils.Error("failed to handle request", err)
	return c.JSON(http.StatusInternalServerError, ResponseError{Message: "server error"})
}
