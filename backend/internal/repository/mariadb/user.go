package mariadb

import (
	"context"
	"database/sql"
	"time"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (m *UserRepository) getOne(ctx context.Context, query string, args ...interface{}) (res domain.User, err error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.User{}, err
	}
	defer stmt.Close()
	row := stmt.QueryRowContext(ctx, args...)
	res = domain.User{}
	err = row.Scan(&res.ID, &res.Email, &res.CreatedAt)
	if err == sql.ErrNoRows {
		return domain.User{}, domain.ErrRecordNotFound
	}
	if err != nil {
		return domain.User{}, err
	}
	return res, nil
}

func (m *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT id, email, created_at FROM user WHERE email = ?`
	return m.getOne(ctx, query, email)
}

func (m *UserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	query := `SELECT id, email, created_at FROM user WHERE id = ?`
	return m.getOne(ctx, query, id)
}

func (m *UserRepository) CreateUser(ctx context.Context, u domain.User) (domain.User, error) {
	query := `INSERT INTO user (email, created_at) VALUES (?, ?)`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return domain.User{}, err
	}
	defer stmt.Close()

	u.CreatedAt = time.Now()
	res, err := stmt.ExecContext(ctx, u.Email, u.CreatedAt)
	if err != nil {
		return domain.User{}, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return domain.User{}, err
	}
	u.ID = id

	return u, nil
}
