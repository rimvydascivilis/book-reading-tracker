package mariadb

import (
	"context"
	"database/sql"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type ReadingRepository struct {
	DB *sql.DB
}

func NewReadingRepository(db *sql.DB) *ReadingRepository {
	return &ReadingRepository{
		DB: db,
	}
}

func (r *ReadingRepository) getAll(ctx context.Context, query string, args ...interface{}) (res []domain.Reading, err error) {
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return []domain.Reading{}, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return []domain.Reading{}, err
	}
	defer rows.Close()

	res = []domain.Reading{}
	for rows.Next() {
		b := domain.Reading{}
		err = rows.Scan(&b.ID, &b.UserID, &b.BookID, &b.TotalPages, &b.Link, &b.CreatedAt, &b.UpdatedAt)
		if err != nil {
			return []domain.Reading{}, err
		}
		res = append(res, b)
	}

	return res, nil
}

func (r *ReadingRepository) getOne(ctx context.Context, query string, args ...interface{}) (domain.Reading, error) {
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Reading{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, args...)

	b := domain.Reading{}
	err = row.Scan(&b.ID, &b.UserID, &b.BookID, &b.TotalPages, &b.Link, &b.CreatedAt, &b.UpdatedAt)
	if err != nil {
		return domain.Reading{}, err
	}

	return b, nil
}

func (r *ReadingRepository) GetReadingsByUserID(ctx context.Context, userID, offset, limit int64) ([]domain.Reading, error) {
	query := `
SELECT id, user_id, book_id, total_pages, COALESCE(link, ''), created_at, updated_at
FROM reading WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	return r.getAll(ctx, query, userID, limit, offset)
}

func (r *ReadingRepository) GetReadingByID(ctx context.Context, id int64) (domain.Reading, error) {
	query := `
SELECT id, user_id, book_id, total_pages, COALESCE(link, ''), created_at, updated_at
FROM reading WHERE id = ?`
	return r.getOne(ctx, query, id)
}

func (r *ReadingRepository) CountReadingsByUserID(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(id) FROM reading WHERE user_id = ?`
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userID)

	var count int64
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *ReadingRepository) CreateReading(ctx context.Context, reading domain.Reading) (domain.Reading, error) {
	query := `INSERT INTO reading (user_id, book_id, total_pages, link, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Reading{}, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, reading.UserID, reading.BookID, reading.TotalPages, reading.Link, reading.CreatedAt, reading.UpdatedAt)
	if err != nil {
		return domain.Reading{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.Reading{}, err
	}

	reading.ID = id
	return reading, nil
}

func (r *ReadingRepository) CountReadingsByUserIDAndBookID(ctx context.Context, userID, bookID int64) (int64, error) {
	query := `SELECT COUNT(id) FROM reading WHERE user_id = ? AND book_id = ?`
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userID, bookID)

	var count int64
	err = row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
