package mariadb

import (
	"context"
	"database/sql"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/dto"
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

func (m *ProgressRepository) GetProgressByReadingAndDate(ctx context.Context, readingID int64, date string) (int64, error) {
	query := `
SELECT COALESCE(SUM(pages), 0) FROM progress
WHERE reading_id = ?
	AND (
		(CHAR_LENGTH(?) = 7 AND DATE_FORMAT(reading_date, '%Y-%m') = ?) OR
		(CHAR_LENGTH(?) = 10 AND DATE(reading_date) = ?)
	)
`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var progress int64
	err = stmt.QueryRow(readingID, date, date, date, date).Scan(&progress)
	if err != nil {
		return 0, err
	}

	return progress, nil
}

func (m *ProgressRepository) GetMonthlyProgress(ctx context.Context, userID, year int64) ([]dto.Progress, error) {
	query := `
SELECT
	MONTH(reading_date) AS date,
	COALESCE(SUM(pages), 0) AS pages
FROM progress
WHERE user_id = ?
	AND YEAR(reading_date) = ?
GROUP BY date
`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monthlyProgress []dto.Progress
	for rows.Next() {
		var progress dto.Progress
		err = rows.Scan(&progress.Date, &progress.Pages)
		if err != nil {
			return nil, err
		}
		monthlyProgress = append(monthlyProgress, progress)
	}

	return monthlyProgress, nil
}

func (m *ProgressRepository) GetDailyProgress(ctx context.Context, userID, year, month int64) ([]dto.Progress, error) {
	query := `
SELECT
	DAY(reading_date) AS date,
	COALESCE(SUM(pages), 0) AS pages
FROM progress
WHERE user_id = ?
	AND YEAR(reading_date) = ?
	AND MONTH(reading_date) = ?
GROUP BY date
`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID, year, month)
	if err != nil {
		return nil, err
	}

	var dailyProgress []dto.Progress
	for rows.Next() {
		var progress dto.Progress
		err = rows.Scan(&progress.Date, &progress.Pages)
		if err != nil {
			return nil, err
		}
		dailyProgress = append(dailyProgress, progress)
	}

	return dailyProgress, nil
}

func (m *ProgressRepository) CreateProgress(ctx context.Context, progress domain.Progress) (domain.Progress, error) {
	query := `INSERT INTO progress (reading_id, user_id, pages, reading_date) VALUES (?, ?, ?, ?)`
	stmt, err := m.DB.Prepare(query)
	if err != nil {
		return domain.Progress{}, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(progress.ReadingID, progress.UserID, progress.Pages, progress.ReadingDate)
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
