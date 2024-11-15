package rest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/utils"
)

type GoalHandler struct {
	GoalSvc domain.GoalService
}

func NewGoalHandler(goalSvc domain.GoalService) *GoalHandler {
	return &GoalHandler{
		GoalSvc: goalSvc,
	}
}

func (h *GoalHandler) GetGoal(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	goal, err := h.GoalSvc.GetGoal(ctx, userID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, goal)
}

func (h *GoalHandler) SetGoal(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var req domain.Goal
	if err := c.Bind(&req); err != nil {
		utils.Error("failed to bind request", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "invalid request format"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	goal, err := h.GoalSvc.SetGoal(ctx, userID, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, goal)
}
