package mariadb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type GoalRepository struct {
	DB *sql.DB
}

func NewGoalRepository(db *sql.DB) *GoalRepository {
	return &GoalRepository{
		DB: db,
	}
}

func (r *GoalRepository) GetGoalByUserID(ctx context.Context, userID int64) (domain.Goal, error) {
	query := `SELECT user_id, type, frequency, value FROM goal WHERE user_id = ?`
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Goal{}, err
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, userID)

	var goal domain.Goal
	err = row.Scan(&goal.UserID, &goal.Type, &goal.Frequency, &goal.Value)
	if err == sql.ErrNoRows {
		return domain.Goal{}, fmt.Errorf("%w: goal for user %d not found", domain.ErrRecordNotFound, userID)
	}
	if err != nil {
		return domain.Goal{}, err
	}

	return goal, nil
}

func (r *GoalRepository) CreateGoal(ctx context.Context, goal domain.Goal) (domain.Goal, error) {
	query := `INSERT INTO goal (user_id, type, frequency, value) VALUES (?, ?, ?, ?)`
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Goal{}, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, goal.UserID, goal.Type, goal.Frequency, goal.Value)
	if err != nil {
		return domain.Goal{}, err
	}

	return goal, nil
}

func (r *GoalRepository) UpdateGoal(ctx context.Context, goal domain.Goal) (domain.Goal, error) {
	query := `UPDATE goal SET type = ?, frequency = ?, value = ? WHERE user_id = ?`
	stmt, err := r.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.Goal{}, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, goal.Type, goal.Frequency, goal.Value, goal.UserID)
	if err != nil {
		return domain.Goal{}, err
	}

	return goal, nil
}
