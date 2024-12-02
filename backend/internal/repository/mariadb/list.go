package mariadb

import (
	"context"
	"database/sql"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type ListRepository struct {
	DB *sql.DB
}

func NewListRepository(db *sql.DB) *ListRepository {
	return &ListRepository{
		DB: db,
	}
}

func (r *ListRepository) getOne(query string, args ...interface{}) (domain.List, error) {
	var list domain.List
	err := r.DB.QueryRow(query, args...).Scan(&list.ID, &list.UserID, &list.Title)
	if err != nil {
		return domain.List{}, err
	}
	return list, nil
}

func (r *ListRepository) getAll(query string, args ...interface{}) ([]domain.List, error) {
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lists []domain.List
	for rows.Next() {
		var list domain.List
		err := rows.Scan(&list.ID, &list.UserID, &list.Title)
		if err != nil {
			return nil, err
		}
		lists = append(lists, list)
	}
	return lists, nil
}

func (r *ListRepository) GetListByID(ctx context.Context, listID int64) (domain.List, error) {
	query := "SELECT id, user_id, title FROM list WHERE id = ?"
	return r.getOne(query, listID)
}

func (r *ListRepository) GetListsByUserID(ctx context.Context, userID int64) ([]domain.List, error) {
	query := "SELECT id, user_id, title FROM list WHERE user_id = ?"
	return r.getAll(query, userID)
}

func (r *ListRepository) CreateList(ctx context.Context, list domain.List) (domain.List, error) {
	query := "INSERT INTO list (user_id, title) VALUES (?, ?)"
	res, err := r.DB.Exec(query, list.UserID, list.Title)
	if err != nil {
		return domain.List{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.List{}, err
	}
	list.ID = id
	return list, nil
}
