package rest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/utils"
)

type ReadingHandler struct {
	ReadingSvc domain.ReadingService
}

func NewReadingHandler(readingSvc domain.ReadingService) *ReadingHandler {
	return &ReadingHandler{
		ReadingSvc: readingSvc,
	}
}

func (h *ReadingHandler) GetReadings(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	page, limit := getPaginationParams(c)
	readings, hasMore, err := h.ReadingSvc.GetReadings(ctx, userID, page, limit)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"readings": readings,
		"hasMore":  hasMore,
	})
}

func (h *ReadingHandler) CreateReading(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var req domain.Reading
	if err := c.Bind(&req); err != nil {
		utils.Error("failed to bind request", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "invalid request format"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	reading, err := h.ReadingSvc.CreateReading(ctx, userID, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusCreated, reading)
}
