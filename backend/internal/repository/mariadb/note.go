package mariadb

import (
	"context"
	"database/sql"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type NoteRepository struct {
	DB *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{
		DB: db,
	}
}

func (r *NoteRepository) getOne(query string, args ...interface{}) (domain.Note, error) {
	row := r.DB.QueryRow(query, args...)
	note := domain.Note{}
	err := row.Scan(&note.ID, &note.UserID, &note.BookID, &note.PageNumber, &note.Content, &note.CreatedAt)
	if err != nil {
		return domain.Note{}, err
	}
	return note, nil
}

func (r *NoteRepository) getAll(query string, args ...interface{}) ([]domain.Note, error) {
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []domain.Note
	for rows.Next() {
		var note domain.Note
		err := rows.Scan(&note.ID, &note.UserID, &note.BookID, &note.PageNumber, &note.Content, &note.CreatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func (r *NoteRepository) GetNoteByUserID(ctx context.Context, noteID, userID int64) (domain.Note, error) {
	query := "SELECT id, user_id, book_id, page_number, content, created_at FROM note WHERE id = ? AND user_id = ?"
	return r.getOne(query, noteID, userID)
}

func (r *NoteRepository) GetBookIDsByUserID(ctx context.Context, userID int64) ([]int64, error) {
	query := "SELECT book_id FROM note WHERE user_id = ? GROUP BY book_id"
	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookIDs []int64
	for rows.Next() {
		var bookID int64
		err := rows.Scan(&bookID)
		if err != nil {
			return nil, err
		}
		bookIDs = append(bookIDs, bookID)
	}
	return bookIDs, nil
}

func (r *NoteRepository) GetNotesByUserIDAndBookID(ctx context.Context, userID, bookID int64) ([]domain.Note, error) {
	query := "SELECT id, user_id, book_id, page_number, content, created_at FROM note WHERE user_id = ? AND book_id = ?"
	return r.getAll(query, userID, bookID)
}

func (r *NoteRepository) CreateNote(ctx context.Context, note domain.Note) (domain.Note, error) {
	query := "INSERT INTO note (user_id, book_id, page_number, content) VALUES (?, ?, ?, ?)"
	res, err := r.DB.Exec(query, note.UserID, note.BookID, note.PageNumber, note.Content)
	if err != nil {
		return domain.Note{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.Note{}, err
	}
	note.ID = id
	return note, nil
}

func (r *NoteRepository) DeleteNote(ctx context.Context, id int64) error {
	query := "DELETE FROM note WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
