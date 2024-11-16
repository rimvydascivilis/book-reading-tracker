package book

import (
	"context"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type bookService struct {
	bookRepo      domain.BookRepository
	validationSvc domain.ValidationService
}

func NewBookService(repo domain.BookRepository, validator domain.ValidationService) domain.BookService {
	return &bookService{
		bookRepo:      repo,
		validationSvc: validator,
	}
}

func (s *bookService) CreateBook(ctx context.Context, userID int64, book domain.Book) (domain.Book, error) {
	book.UserID = userID
	if err := s.validationSvc.ValidateStruct(book); err != nil {
		return domain.Book{}, err
	}

	return s.bookRepo.CreateBook(ctx, book)
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

	totalCount, err := s.bookRepo.CountBooksByUser(ctx, userID)
	if err != nil {
		return nil, false, err
	}

	hasMore := totalCount > page*limit
	if totalCount == 0 {
		return []domain.Book{}, false, nil
	}

	offset := (page - 1) * limit
	books, err := s.bookRepo.GetBooksByUser(ctx, userID, offset, limit)
	if err != nil {
		return nil, false, err
	}

	return books, hasMore, nil
}

func (s *bookService) UpdateBook(ctx context.Context, userID int64, book domain.Book) (domain.Book, error) {
	currBook, err := s.bookRepo.GetBookByUserID(ctx, userID, book.ID)
	if err != nil {
		return domain.Book{}, err
	}

	if book.Title != "" {
		currBook.Title = book.Title
	}
	if book.Rating != 0 {
		currBook.Rating = book.Rating
	}
	currBook.UserID = userID

	if err := s.validationSvc.ValidateStruct(currBook); err != nil {
		return domain.Book{}, err
	}

	return s.bookRepo.UpdateBook(ctx, currBook)
}

func (s *bookService) DeleteBook(ctx context.Context, userID, bookID int64) error {
	_, err := s.bookRepo.GetBookByUserID(ctx, userID, bookID)
	if err != nil {
		return err
	}

	return s.bookRepo.DeleteBook(ctx, userID, bookID)
}
