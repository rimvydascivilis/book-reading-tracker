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

type NoteHandler struct {
	NoteSvc domain.NoteService
}

func NewNoteHandler(noteSvc domain.NoteService) *NoteHandler {
	return &NoteHandler{
		NoteSvc: noteSvc,
	}
}

func (h *NoteHandler) GetNotes(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	bookID, err := strconv.ParseInt(c.Param("book_id"), 10, 64)
	if err != nil {
		utils.Error("failed to get book id from param", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid book id"})
	}

	res, err := h.NoteSvc.GetNotes(ctx, userID, bookID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *NoteHandler) CreateNote(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	bookID, err := strconv.ParseInt(c.Param("book_id"), 10, 64)
	if err != nil {
		utils.Error("failed to get book id from param", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid book id"})
	}

	var note dto.NoteRequest
	if err := c.Bind(&note); err != nil {
		utils.Error("failed to bind note", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid note"})
	}

	res, err := h.NoteSvc.CreateNote(ctx, userID, bookID, note)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *NoteHandler) DeleteNote(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	noteID, err := strconv.ParseInt(c.Param("note_id"), 10, 64)
	if err != nil {
		utils.Error("failed to get note id from param", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid note id"})
	}

	err = h.NoteSvc.DeleteNote(ctx, userID, noteID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
