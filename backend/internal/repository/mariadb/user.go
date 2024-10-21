package mariadb

import (
	"context"
	"database/sql"

	"book-tracker/domain"
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
	row := stmt.QueryRowContext(ctx, args...)
	res = domain.User{}

	err = row.Scan(
		&res.ID,
		&res.Email,
		&res.CreatedAt,
	)
	if err != nil {
		return domain.User{}, err
	}

	return res, nil
}

func (m *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT id, email, created_at FROM author WHERE email = ?`
	return m.getOne(ctx, query, email)
}

func (m *UserRepository) Create(ctx context.Context, u *domain.User) error {
	query := `INSERT INTO author (email, created_at) VALUES (?, ?)`
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	_, err = stmt.ExecContext(ctx, u.Email, u.CreatedAt)
	return err
}
