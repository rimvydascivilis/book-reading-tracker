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

type ListHandler struct {
	ListSvc domain.ListService
}

func NewListHandler(listSvc domain.ListService) *ListHandler {
	return &ListHandler{
		ListSvc: listSvc,
	}
}

func (h *ListHandler) ListLists(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	lists, err := h.ListSvc.ListLists(ctx, userID)
	utils.Info("lists", map[string]interface{}{"lists": lists})
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, lists)
}

func (h *ListHandler) GetList(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	listIDString := c.QueryParam("list_id")
	if listIDString == "" {
		utils.Error("failed to get list id from param", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid list id"})
	}

	listID, err := strconv.ParseInt(listIDString, 10, 64)
	if err != nil {
		utils.Error("failed to parse list id", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid list id"})
	}

	list, err := h.ListSvc.GetList(ctx, userID, listID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusOK, list)
}

func (h *ListHandler) CreateList(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var req dto.ListRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("failed to bind request", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid request"})
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	list, err := h.ListSvc.CreateList(ctx, userID, req)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.JSON(http.StatusCreated, list)
}

type AddBookToListRequest struct {
	BookID int64 `json:"book_id"`
	ListID int64 `json:"list_id"`
}

func (h *ListHandler) AddBookToList(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	var req AddBookToListRequest
	if err := c.Bind(&req); err != nil {
		utils.Error("failed to bind request", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid request"})
	}

	err = h.ListSvc.AddBookToList(ctx, userID, req.ListID, req.BookID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ListHandler) RemoveBookFromList(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userID, err := getUserIDFromToken(c)
	if err != nil {
		utils.Error("failed to get user id from token", err)
		return c.JSON(http.StatusUnauthorized, ResponseError{Message: "Invalid token"})
	}

	listIDString := c.Param("list_id")
	if listIDString == "" {
		utils.Error("failed to get list id from param", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid list id"})
	}

	listID, err := strconv.ParseInt(listIDString, 10, 64)
	if err != nil {
		utils.Error("failed to parse list id", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid list id"})
	}

	itemIDString := c.Param("item_id")
	if itemIDString == "" {
		utils.Error("failed to get book id from param", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid book id"})
	}

	bookID, err := strconv.ParseInt(itemIDString, 10, 64)
	if err != nil {
		utils.Error("failed to parse book id", err)
		return c.JSON(http.StatusBadRequest, ResponseError{Message: "Invalid book id"})
	}

	err = h.ListSvc.RemoveBookFromList(ctx, userID, listID, bookID)
	if err != nil {
		return handleServiceError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
