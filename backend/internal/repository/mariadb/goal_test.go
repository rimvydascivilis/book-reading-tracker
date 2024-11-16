package mariadb_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rimvydascivilis/book-tracker/backend/domain"
	"github.com/rimvydascivilis/book-tracker/backend/internal/repository/mariadb"
	"github.com/stretchr/testify/assert"
)

func setupGoalRepository(t *testing.T) (*mariadb.GoalRepository, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	t.Cleanup(func() { db.Close() })
	goalRepo := mariadb.NewGoalRepository(db)
	return goalRepo, mock
}

func TestGoalRepository_GetGoalByUserID_Success(t *testing.T) {
	goalRepo, mock := setupGoalRepository(t)

	ctx := context.Background()
	testUserID := int64(1)
	testGoal := domain.Goal{
		UserID:    testUserID,
		Type:      "books",
		Frequency: "daily",
		Value:     10,
	}

	rows := sqlmock.NewRows([]string{"user_id", "type", "frequency", "value"}).
		AddRow(testGoal.UserID, testGoal.Type, testGoal.Frequency, testGoal.Value)

	mock.ExpectPrepare(`SELECT user_id, type, frequency, value FROM goal WHERE user_id = \?`).
		ExpectQuery().
		WithArgs(testUserID).
		WillReturnRows(rows)

	goal, err := goalRepo.GetGoalByUserID(ctx, testUserID)

	assert.NoError(t, err)
	assert.Equal(t, testGoal, goal)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGoalRepository_GetGoalByUserID_NotFound(t *testing.T) {
	goalRepo, mock := setupGoalRepository(t)

	ctx := context.Background()
	testUserID := int64(999)

	mock.ExpectPrepare(`SELECT user_id, type, frequency, value FROM goal WHERE user_id = \?`).
		ExpectQuery().
		WithArgs(testUserID).
		WillReturnError(sql.ErrNoRows)

	goal, err := goalRepo.GetGoalByUserID(ctx, testUserID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "goal for user 999 not found")
	assert.Equal(t, domain.Goal{}, goal)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGoalRepository_CreateGoal_Success(t *testing.T) {
	goalRepo, mock := setupGoalRepository(t)

	ctx := context.Background()
	testGoal := domain.Goal{
		UserID:    1,
		Type:      "pages",
		Frequency: "monthly",
		Value:     100,
	}

	mock.ExpectPrepare(`INSERT INTO goal`).
		ExpectExec().
		WithArgs(testGoal.UserID, testGoal.Type, testGoal.Frequency, testGoal.Value).
		WillReturnResult(sqlmock.NewResult(0, 1))

	createdGoal, err := goalRepo.CreateGoal(ctx, testGoal)

	assert.NoError(t, err)
	assert.Equal(t, testGoal, createdGoal)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGoalRepository_CreateGoal_Failure(t *testing.T) {
	goalRepo, mock := setupGoalRepository(t)

	ctx := context.Background()
	testGoal := domain.Goal{
		UserID:    1,
		Type:      "pages",
		Frequency: "monthly",
		Value:     100,
	}

	mock.ExpectPrepare(`INSERT INTO goal`).
		ExpectExec().
		WithArgs(testGoal.UserID, testGoal.Type, testGoal.Frequency, testGoal.Value).
		WillReturnError(errors.New("database error"))

	createdGoal, err := goalRepo.CreateGoal(ctx, testGoal)

	assert.Error(t, err)
	assert.Equal(t, domain.Goal{}, createdGoal)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGoalRepository_UpdateGoal_Success(t *testing.T) {
	goalRepo, mock := setupGoalRepository(t)

	ctx := context.Background()
	testGoal := domain.Goal{
		UserID:    1,
		Type:      "books",
		Frequency: "daily",
		Value:     20,
	}

	mock.ExpectPrepare(`UPDATE goal SET type = \?, frequency = \?, value = \? WHERE user_id = \?`).
		ExpectExec().
		WithArgs(testGoal.Type, testGoal.Frequency, testGoal.Value, testGoal.UserID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	updatedGoal, err := goalRepo.UpdateGoal(ctx, testGoal)

	assert.NoError(t, err)
	assert.Equal(t, testGoal, updatedGoal)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGoalRepository_UpdateGoal_Failure(t *testing.T) {
	goalRepo, mock := setupGoalRepository(t)

	ctx := context.Background()
	testGoal := domain.Goal{
		UserID:    1,
		Type:      "books",
		Frequency: "daily",
		Value:     20,
	}

	mock.ExpectPrepare(`UPDATE goal SET type = \?, frequency = \?, value = \? WHERE user_id = \?`).
		ExpectExec().
		WithArgs(testGoal.Type, testGoal.Frequency, testGoal.Value, testGoal.UserID).
		WillReturnError(errors.New("database error"))

	updatedGoal, err := goalRepo.UpdateGoal(ctx, testGoal)

	assert.Error(t, err)
	assert.Equal(t, domain.Goal{}, updatedGoal)
	assert.NoError(t, mock.ExpectationsWereMet())
}
