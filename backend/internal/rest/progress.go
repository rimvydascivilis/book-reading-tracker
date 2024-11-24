package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/dto"
	"github.com/rimvydascivilis/book-tracker/backend/utils"
)

type ProgressHandler struct {
	ProgressSvc domain.ProgressService
}

func NewProgressHandler(progressSvc domain.ProgressService) *ProgressHandler {
	return &ProgressHandler{
		ProgressSvc: progressSvc,
	}
}

func (h *ProgressHandler) CreateProgress(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	var req dto.ProgressRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("failed to bind request", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "invalid request format"})
	}

	readingID, err := strconv.ParseInt(c.Param("readingId"), 10, 64)
	if err != nil {
		utils.Error("failed to parse reading ID", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "invalid reading ID"})
	}

	progress, err := h.ProgressSvc.CreateProgress(ctx, userID, readingID, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusCreated, progress)
}
