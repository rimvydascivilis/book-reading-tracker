package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/utils"
)

type StatHandler struct {
	StatSvc domain.StatService
}

func NewStatHandler(statSvc domain.StatService) *StatHandler {
	return &StatHandler{
		StatSvc: statSvc,
	}
}

func (h *StatHandler) GetProgress(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	yearString := c.QueryParam("year")
	if yearString == "" {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Year is required"})
	}

	year, err := strconv.Atoi(yearString)
	if err != nil {
		utils.Error("failed to parse year", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid year"})
	}

	frequency := c.Param("frequency") // Expected: "monthly" or "daily"
	isMonthly := frequency == "monthly"

	var month int
	if !isMonthly {
		monthString := c.QueryParam("month")
		if monthString == "" {
			return c.JSON(http.StatusBadRequest, ResponseError{Message: "Month is required for daily progress"})
		}

		month, err = strconv.Atoi(monthString)
		if err != nil {
			utils.Error("failed to parse month", err)
			return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid month"})
		}
	}

	res, err := h.StatSvc.GetProgress(ctx, userID, int64(year), int64(month), isMonthly)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}
