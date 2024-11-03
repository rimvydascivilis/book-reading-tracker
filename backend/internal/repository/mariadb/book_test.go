package mariadb_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/internal/repository/mariadb"
	"github.com/stretchr/testify/assert"
)

func setupBookRepository(t *testing.T) (*mariadb.BookRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	bookRepo := mariadb.NewBookRepository(db)
	return bookRepo, mock
}

func TestBookRepository_GetBooksByUser_Success(t *testing.T) {
	bookRepo, mock := setupBookRepository(t)

	var rating float64 = 5
	var rating2 float64 = 4
	ctx := context.Background()
	testBooks := []domain.Book{
		{ID: 1, Title: "Book 1", Rating: &rating, CreatedAt: time.Now()},
		{ID: 2, Title: "Book 2", Rating: &rating2, CreatedAt: time.Now()},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "rating", "created_at"}).
		AddRow(testBooks[0].ID, testBooks[0].Title, testBooks[0].Rating, testBooks[0].CreatedAt).
		AddRow(testBooks[1].ID, testBooks[1].Title, testBooks[1].Rating, testBooks[1].CreatedAt)

	mock.ExpectPrepare(`SELECT id, title, rating, created_at FROM book WHERE user_id = \? LIMIT \? OFFSET \?`).
		ExpectQuery().
		WithArgs(1, int64(10), int64(0)).
		WillReturnRows(rows)

	books, err := bookRepo.GetBooksByUser(ctx, 1, 0, 10)

	assert.NoError(t, err)
	assert.Equal(t, testBooks, books)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_CountBooksByUser_Success(t *testing.T) {
	bookRepo, mock := setupBookRepository(t)

	ctx := context.Background()
	userID := int64(1)
	expectedCount := int64(2)

	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM book WHERE user_id = \?`).
		WithArgs(userID).
		WillReturnRows(sqlmock.NewRows([]string{`COUNT(*)`}).AddRow(expectedCount))

	count, err := bookRepo.CountBooksByUser(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_CreateBook_Success(t *testing.T) {
	bookRepo, mock := setupBookRepository(t)

	var rating float64 = 5
	ctx := context.Background()
	newBook := domain.Book{
		Title:     "New Book",
		Rating:    &rating,
		CreatedAt: time.Now(),
	}

	mock.ExpectPrepare(`INSERT INTO book \(user_id, title, rating, created_at\) VALUES \(\?, \?, \?, \?\)`).
		ExpectExec().
		WithArgs(1, newBook.Title, newBook.Rating, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	book, err := bookRepo.CreateBook(ctx, 1, newBook)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), book.ID)
	assert.Equal(t, newBook.Title, book.Title)
	assert.Equal(t, newBook.Rating, book.Rating)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_GetBookByUserID_Success(t *testing.T) {
	bookRepo, mock := setupBookRepository(t)

	ctx := context.Background()
	var rating float64 = 5
	testBook := domain.Book{ID: 1, Title: "Book 1", Rating: &rating, CreatedAt: time.Now()}

	rows := sqlmock.NewRows([]string{"id", "title", "rating", "created_at"}).
		AddRow(testBook.ID, testBook.Title, testBook.Rating, testBook.CreatedAt)

	mock.ExpectPrepare(`SELECT id, title, rating, created_at FROM book WHERE user_id = \? AND id = \?`).
		ExpectQuery().
		WithArgs(1, 1).
		WillReturnRows(rows)

	book, err := bookRepo.GetBookByUserID(ctx, 1, 1)

	assert.NoError(t, err)
	assert.Equal(t, testBook, book)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_GetBookByUserID_NotFound(t *testing.T) {
	bookRepo, mock := setupBookRepository(t)

	ctx := context.Background()

	mock.ExpectPrepare(`SELECT id, title, rating, created_at FROM book WHERE user_id = \? AND id = \?`).
		ExpectQuery().
		WithArgs(1, 99).
		WillReturnError(sql.ErrNoRows)

	book, err := bookRepo.GetBookByUserID(ctx, 1, 99)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrBookNotFound, err)
	assert.Equal(t, domain.Book{}, book)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_UpdateBook_Success(t *testing.T) {
	bookRepo, mock := setupBookRepository(t)

	ctx := context.Background()
	var rating float64 = 5
	updateBook := domain.Book{ID: 1, Title: "Updated Book", Rating: &rating}

	mock.ExpectPrepare(`UPDATE book SET title = \?, rating = \? WHERE user_id = \? AND id = \?`).
		ExpectExec().
		WithArgs(updateBook.Title, updateBook.Rating, 1, updateBook.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	book, err := bookRepo.UpdateBook(ctx, 1, updateBook)

	assert.NoError(t, err)
	assert.Equal(t, updateBook.Title, book.Title)
	assert.Equal(t, updateBook.Rating, book.Rating)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestBookRepository_DeleteBook_Success(t *testing.T) {
	bookRepo, mock := setupBookRepository(t)

	ctx := context.Background()
	bookID := int64(1)

	mock.ExpectPrepare(`DELETE FROM book WHERE user_id = \? AND id = \?`).
		ExpectExec().
		WithArgs(1, bookID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := bookRepo.DeleteBook(ctx, 1, bookID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
