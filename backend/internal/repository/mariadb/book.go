package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type BookRepository struct {
	DB *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{
		DB: db,
	}
}

func (m *BookRepository) getOne(ctx context.Context, query string, args ...interface{}) (domain.Book, error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Book{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)
	b := domain.Book{}
	err = row.Scan(&b.ID, &b.UserID, &b.Title, &b.Rating, &b.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.Book{}, fmt.Errorf("%w: %s", domain.ErrRecordNotFound, "book")
	}
	if err != nil {
		return domain.Book{}, err
	}

	return b, nil
}

func (m *BookRepository) getAll(ctx context.Context, query string, args ...interface{}) (res []domain.Book, err error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return []domain.Book{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return []domain.Book{}, err
	}
	defer rows.Close()

	res = []domain.Book{}
	for rows.Next() {
		b := domain.Book{}
		err = rows.Scan(&b.ID, &b.Title, &b.Rating, &b.CreatedAt)
		if err != nil {
			return []domain.Book{}, err
		}
		res = append(res, b)
	}

	return res, nil
}

func (m *BookRepository) GetBooksByUser(ctx context.Context, userID, offset, limit int64) ([]domain.Book, error) {
	query := `SELECT id, title, rating, created_at FROM book WHERE user_id = ? LIMIT ? OFFSET ?`
	return m.getAll(ctx, query, userID, limit, offset)
}

func (m *BookRepository) CountBooksByUser(ctx context.Context, userID int64) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM book WHERE user_id = ?`
	err := m.DB.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (m *BookRepository) CreateBook(ctx context.Context, b domain.Book) (domain.Book, error) {
	query := `INSERT INTO book (user_id, title, rating, created_at) VALUES (?, ?, ?, ?)`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Book{}, err
	}
	defer stmt.Close()

	b.CreatedAt = time.Now()
	res, err := stmt.ExecContext(ctx, b.UserID, b.Title, b.Rating, b.CreatedAt)
	if err != nil {
		return domain.Book{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.Book{}, err
	}

	b.ID = id
	return b, nil
}

func (m *BookRepository) GetBookByUserID(ctx context.Context, userID, bookID int64) (domain.Book, error) {
	query := `SELECT id, user_id, title, rating, created_at FROM book WHERE user_id = ? AND id = ?`
	return m.getOne(ctx, query, userID, bookID)
}

func (m *BookRepository) SearchBooksByTitle(ctx context.Context, userID int64, title string, limit int64) ([]domain.Book, error) {
	query := `SELECT id, title, rating, created_at FROM book WHERE user_id = ? AND title LIKE ? LIMIT ?`
	return m.getAll(ctx, query, userID, "%"+title+"%", limit)
}

func (m *BookRepository) UpdateBook(ctx context.Context, b domain.Book) (domain.Book, error) {
	query := `UPDATE book SET title = ?, rating = ? WHERE user_id = ? AND id = ?`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Book{}, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, b.Title, b.Rating, b.UserID, b.ID)
	if err != nil {
		return domain.Book{}, err
	}

	return b, nil
}

func (m *BookRepository) DeleteBook(ctx context.Context, userID, bookID int64) error {
	query := `DELETE FROM book WHERE user_id = ? AND id = ?`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, userID, bookID)
	if err != nil {
		return err
	}

	return nil
}
