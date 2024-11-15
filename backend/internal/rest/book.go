package rest

import (
	"context"
	"net/http"
	"strconv"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/utils"

	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	BookSvc domain.BookService
}

func NewBookHandler(bookSvc domain.BookService) *BookHandler {
	handler := &BookHandler{
		BookSvc: bookSvc,
	}
	return handler
}

func (a *BookHandler) GetBooks(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	page := getInt64QueryParam(c, "page", 1)
	limit := getInt64QueryParam(c, "limit", 10)

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	books, hasMore, err := a.BookSvc.GetBooks(ctx, userID, page, limit)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"books":   books,
		"hasMore": hasMore,
	})
}

func (a *BookHandler) CreateBook(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var req domain.Book
	if err := c.Bind(&req); err != nil {
		utils.Error("failed to bind request", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "invalid request format"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	book, err := a.BookSvc.CreateBook(ctx, userID, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusCreated, book)
}

func (a *BookHandler) UpdateBook(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var req domain.Book
	if err := c.Bind(&req); err != nil {
		utils.Error("failed to bind request", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "invalid request format"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	bookID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error("failed to parse book id", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "invalid book id"})
	}

	req.ID = bookID

	book, err := a.BookSvc.UpdateBook(ctx, userID, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, book)
}

func (a *BookHandler) DeleteBook(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	bookID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error("failed to parse book id", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "invalid book id"})
	}

	err = a.BookSvc.DeleteBook(ctx, userID, bookID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
