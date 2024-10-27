package mariadb_test

import (
	"book-tracker/domain"
	"book-tracker/internal/repository/mariadb"
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetByEmail_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := mariadb.NewUserRepository(db)

	ctx := context.Background()
	testEmail := "test@example.com"
	testUser := domain.User{
		ID:        1,
		Email:     testEmail,
		CreatedAt: time.Now(),
	}

	rows := sqlmock.NewRows([]string{"id", "email", "created_at"}).
		AddRow(testUser.ID, testUser.Email, testUser.CreatedAt)

	mock.ExpectPrepare("SELECT id, email, created_at FROM user WHERE email = ?").
		ExpectQuery().
		WithArgs(testEmail).
		WillReturnRows(rows)

	user, err := userRepo.GetByEmail(ctx, testEmail)

	assert.NoError(t, err)
	assert.Equal(t, testUser, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_GetByEmail_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := mariadb.NewUserRepository(db)

	ctx := context.Background()
	testEmail := "nonexistent@example.com"

	mock.ExpectPrepare("SELECT id, email, created_at FROM user WHERE email = ?").
		ExpectQuery().
		WithArgs(testEmail).
		WillReturnError(sql.ErrNoRows)

	user, err := userRepo.GetByEmail(ctx, testEmail)

	assert.Error(t, err)
	assert.Equal(t, domain.ErrUserNotFound, err)
	assert.Equal(t, domain.User{}, user)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepository_CreateUser_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	userRepo := mariadb.NewUserRepository(db)

	ctx := context.Background()
	testUser := domain.User{
		Email: "newuser@example.com",
	}

	mock.ExpectPrepare("INSERT INTO user").
		ExpectExec().
		WithArgs(testUser.Email, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	createdUser, err := userRepo.CreateUser(ctx, testUser)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), createdUser.ID)
	assert.Equal(t, testUser.Email, createdUser.Email)
	assert.NotZero(t, createdUser.CreatedAt)

	assert.NoError(t, mock.ExpectationsWereMet())
}
