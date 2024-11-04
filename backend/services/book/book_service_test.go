package book

import (
	"context"
	"fmt"
	"testing"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupBookService() (domain.BookService, *mocks.BookRepository, *mocks.ValidationService) {
	bookRepo := new(mocks.BookRepository)
	validationSvc := new(mocks.ValidationService)
	userService := NewBookService(bookRepo, validationSvc)

	return userService, bookRepo, validationSvc
}

func TestCreateBook(t *testing.T) {
	service, mockRepo, validationSvc := setupBookService()

	userID := int64(1)
	book := domain.Book{
		Title:  "New Book",
		Rating: 4,
	}

	mockRepo.On("CreateBook", mock.Anything, userID, book).Return(book, nil)
	validationSvc.On("ValidateStruct", book).Return(nil)

	createdBook, err := service.CreateBook(context.Background(), userID, book)

	assert.NoError(t, err)
	assert.Equal(t, book.Title, createdBook.Title)
	mockRepo.AssertExpectations(t)
}

func TestCreateBook_EmptyTitle(t *testing.T) {
	service, _, validationSvc := setupBookService()

	userID := int64(1)
	book := domain.Book{
		Rating: 4,
	}
	validationSvc.On("ValidateStruct", book).Return(fmt.Errorf("Title is required"))

	createdBook, err := service.CreateBook(context.Background(), userID, book)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Title")
	assert.Equal(t, domain.Book{}, createdBook)
}

func TestCreateBook_InvalidRating(t *testing.T) {
	service, _, validationSvc := setupBookService()

	userID := int64(1)
	book := domain.Book{
		Title:  "New Book",
		Rating: 100, // Invalid rating
	}

	validationSvc.On("ValidateStruct", book).Return(fmt.Errorf("Rating must be between 1 and 5"))

	createdBook, err := service.CreateBook(context.Background(), userID, book)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Rating")
	assert.Equal(t, domain.Book{}, createdBook)
}

func TestGetBooks(t *testing.T) {
	service, mockRepo, _ := setupBookService()

	userID := int64(1)
	page := int64(1)
	limit := int64(10)
	books := []domain.Book{
		{ID: 1, Title: "Book 1", Rating: 5},
		{ID: 2, Title: "Book 2", Rating: 4},
	}

	mockRepo.On("GetBooksByUser", mock.Anything, userID, int64(0), limit).Return(books, nil)
	mockRepo.On("CountBooksByUser", mock.Anything, userID).Return(int64(15), nil)

	resultBooks, hasMore, err := service.GetBooks(context.Background(), userID, page, limit)

	assert.NoError(t, err)
	assert.Equal(t, books, resultBooks)
	assert.True(t, hasMore)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBook(t *testing.T) {
	service, mockRepo, validationSvc := setupBookService()

	userID := int64(1)
	book := domain.Book{
		ID:     1,
		Title:  "Updated Book",
		Rating: 5,
	}
	existingBook := domain.Book{
		ID:     1,
		Title:  "Old Title",
		Rating: 4,
	}

	mockRepo.On("GetBookByUserID", mock.Anything, userID, book.ID).Return(existingBook, nil)
	validationSvc.On("ValidateStruct", book).Return(nil)
	mockRepo.On("UpdateBook", mock.Anything, userID, book).Return(book, nil)

	updatedBook, err := service.UpdateBook(context.Background(), userID, book)

	assert.NoError(t, err)
	assert.Equal(t, book.Title, updatedBook.Title)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBook_InvalidRating(t *testing.T) {
	service, mockRepo, validationSvc := setupBookService()

	userID := int64(1)
	book := domain.Book{
		ID:     1,
		Title:  "Updated Book",
		Rating: 100,
	}
	existingBook := domain.Book{
		ID:     1,
		Title:  "Old Title",
		Rating: 4,
	}

	mockRepo.On("GetBookByUserID", mock.Anything, userID, book.ID).Return(existingBook, nil)
	validationSvc.On("ValidateStruct", book).Return(fmt.Errorf("Rating must be between 1 and 5"))

	updatedBook, err := service.UpdateBook(context.Background(), userID, book)

	assert.Error(t, err)
	assert.Equal(t, domain.Book{}, updatedBook)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBook(t *testing.T) {
	service, mockRepo, _ := setupBookService()

	userID := int64(1)
	bookID := int64(1)

	mockRepo.On("GetBookByUserID", mock.Anything, userID, bookID).Return(domain.Book{}, nil)
	mockRepo.On("DeleteBook", mock.Anything, userID, bookID).Return(nil)

	err := service.DeleteBook(context.Background(), userID, bookID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBook_NotFound(t *testing.T) {
	service, mockRepo, _ := setupBookService()

	userID := int64(1)
	bookID := int64(1)

	mockRepo.On("GetBookByUserID", mock.Anything, userID, bookID).Return(domain.Book{}, domain.ErrRecordNotFound)

	err := service.DeleteBook(context.Background(), userID, bookID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrRecordNotFound, err)
	mockRepo.AssertExpectations(t)
}
