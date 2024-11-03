package book

import (
	"context"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type bookService struct {
	bookRepo domain.BookRepository
}

func NewBookService(repo domain.BookRepository) *bookService {
	return &bookService{
		bookRepo: repo,
	}
}

func (s *bookService) CreateBook(ctx context.Context, userID int64, book domain.Book) (domain.Book, error) {
	if book.Title == "" {
		return domain.Book{}, domain.ErrInvalidBook
	}

	return s.bookRepo.CreateBook(ctx, userID, book)
}

func (s *bookService) GetBooks(ctx context.Context, userID, page, limit int64) ([]domain.Book, bool, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	books, err := s.bookRepo.GetBooksByUser(ctx, userID, offset, limit)
	if err != nil {
		return nil, false, err
	}

	totalCount, err := s.bookRepo.CountBooksByUser(ctx, userID)
	if err != nil {
		return nil, false, err
	}

	hasMore := totalCount > page*limit

	return books, hasMore, nil
}

func (s *bookService) UpdateBook(ctx context.Context, userID int64, book domain.Book) (domain.Book, error) {
	currBook, err := s.bookRepo.GetBookByUserID(ctx, userID, book.ID)
	if err != nil {
		return domain.Book{}, err
	}

	if book.Title == "" {
		book.Title = currBook.Title
	}
	if book.Rating == nil {
		book.Rating = currBook.Rating
	}

	return s.bookRepo.UpdateBook(ctx, userID, book)
}

func (s *bookService) DeleteBook(ctx context.Context, userID, bookID int64) error {
	_, err := s.bookRepo.GetBookByUserID(ctx, userID, bookID)
	if err != nil {
		return err
	}

	return s.bookRepo.DeleteBook(ctx, userID, bookID)
}
