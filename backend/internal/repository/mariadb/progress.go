package mariadb

import (
	"context"
	"database/sql"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type ProgressRepository struct {
	DB *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{
		DB: db,
	}
}

func (m *ProgressRepository) GetTotalProgressByReadingID(ctx context.Context, readingID int64) (int64, error) {
	query := `SELECT COALESCE(SUM(pages), 0) FROM progress WHERE reading_id = ?`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var totalProgress int64
	err = stmt.QueryRow(readingID).Scan(&totalProgress)
	if err != nil {
		return 0, err
	}

	return totalProgress, nil
}

func (m *ProgressRepository) GetUserReadingIDsByPeriod(ctx context.Context, userID int64, period string) ([]int64, error) {
	query := `
SELECT UNIQUE reading_id
FROM progress
WHERE user_id = ?
	AND (
		(CHAR_LENGTH(?) = 7 AND DATE_FORMAT(reading_date, '%Y-%m') = ?) OR
		(CHAR_LENGTH(?) = 10 AND DATE(reading_date) = ?)
	)
`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID, period, period, period, period)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readingIDs []int64
	for rows.Next() {
		var readingID int64
		err = rows.Scan(&readingID)
		if err != nil {
			return nil, err
		}
		readingIDs = append(readingIDs, readingID)
	}

	return readingIDs, nil
}

func (m *ProgressRepository) CreateProgress(ctx context.Context, progress domain.Progress) (domain.Progress, error) {
	query := `INSERT INTO progress (reading_id, user_id, pages) VALUES (?, ?, ?)`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return domain.Progress{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(progress.ReadingID, progress.UserID, progress.Pages)
	if err != nil {
		return domain.Progress{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.Progress{}, err
	}

	progress.ID = id
	return progress, nil
}
