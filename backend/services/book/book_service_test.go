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

func TestCreateBook(t *testing.T) {
	mockRepo := new(mocks.BookRepository)
	service := NewBookService(mockRepo)

	userID := int64(1)
	var rating float64 = 4
	book := domain.Book{
		Title:  "New Book",
		Rating: &rating,
	}

	mockRepo.On("CreateBook", mock.Anything, userID, book).Return(book, nil)

	createdBook, err := service.CreateBook(context.Background(), userID, book)

	assert.NoError(t, err)
	assert.Equal(t, book.Title, createdBook.Title)
	mockRepo.AssertExpectations(t)
}

func TestCreateBook_EmptyTitle(t *testing.T) {
	mockRepo := new(mocks.BookRepository)
	service := NewBookService(mockRepo)

	userID := int64(1)
	var rating float64 = 4
	book := domain.Book{
		Title:  "", // Invalid title
		Rating: &rating,
	}

	createdBook, err := service.CreateBook(context.Background(), userID, book)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrInvalidBook, err)
	assert.Equal(t, domain.Book{}, createdBook)
}

func TestGetBooks(t *testing.T) {
	mockRepo := new(mocks.BookRepository)
	service := NewBookService(mockRepo)

	userID := int64(1)
	page := int64(1)
	limit := int64(10)
	var rating float64 = 5
	var rating2 float64 = 4
	books := []domain.Book{
		{ID: 1, Title: "Book 1", Rating: &rating},
		{ID: 2, Title: "Book 2", Rating: &rating2},
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
	mockRepo := new(mocks.BookRepository)
	service := NewBookService(mockRepo)

	userID := int64(1)
	var rating float64 = 5
	var rating2 float64 = 4
	book := domain.Book{
		ID:     1,
		Title:  "Updated Book",
		Rating: &rating,
	}
	existingBook := domain.Book{
		ID:     1,
		Title:  "Old Title",
		Rating: &rating2,
	}

	mockRepo.On("GetBookByUserID", mock.Anything, userID, book.ID).Return(existingBook, nil)
	mockRepo.On("UpdateBook", mock.Anything, userID, book).Return(book, nil)

	updatedBook, err := service.UpdateBook(context.Background(), userID, book)

	assert.NoError(t, err)
	assert.Equal(t, book.Title, updatedBook.Title)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBook_InvalidRating(t *testing.T) {
	mockRepo := new(mocks.BookRepository)
	service := NewBookService(mockRepo)

	userID := int64(1)
	var rating float64 = 5
	var invalidRating float64 = -1
	book := domain.Book{
		ID:     1,
		Title:  "Updated Book",
		Rating: &invalidRating,
	}
	existingBook := domain.Book{
		ID:     1,
		Title:  "Old Title",
		Rating: &rating,
	}

	mockRepo.On("GetBookByUserID", mock.Anything, userID, book.ID).Return(existingBook, nil)
	mockRepo.On("UpdateBook", mock.Anything, userID, book).Return(domain.Book{}, fmt.Errorf("constraint violation"))

	updatedBook, err := service.UpdateBook(context.Background(), userID, book)

	assert.Error(t, err)
	assert.Equal(t, domain.Book{}, updatedBook)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBook(t *testing.T) {
	mockRepo := new(mocks.BookRepository)
	service := NewBookService(mockRepo)

	userID := int64(1)
	bookID := int64(1)

	mockRepo.On("GetBookByUserID", mock.Anything, userID, bookID).Return(domain.Book{}, nil)
	mockRepo.On("DeleteBook", mock.Anything, userID, bookID).Return(nil)

	err := service.DeleteBook(context.Background(), userID, bookID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBook_NotFound(t *testing.T) {
	mockRepo := new(mocks.BookRepository)
	service := NewBookService(mockRepo)

	userID := int64(1)
	bookID := int64(1)

	mockRepo.On("GetBookByUserID", mock.Anything, userID, bookID).Return(domain.Book{}, domain.ErrBookNotFound)

	err := service.DeleteBook(context.Background(), userID, bookID)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrBookNotFound, err)
	mockRepo.AssertExpectations(t)
}
