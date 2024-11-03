package rest

import (
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
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
		utils.Info("user ID not found in token", map[string]interface{}{"claims": userID})
		return 0, errors.New("user ID not found in token")
	}

	return int64(userID), nil
}
