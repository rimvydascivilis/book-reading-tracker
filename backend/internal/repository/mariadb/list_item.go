package mariadb

import (
	"context"
	"database/sql"

	"github.com/rimvydascivilis/book-tracker/backend/domain"
)

type ListItemRepository struct {
	DB *sql.DB
}

func NewListItemRepository(db *sql.DB) *ListItemRepository {
	return &ListItemRepository{
		DB: db,
	}
}

func (r *ListItemRepository) getAll(query string, args ...interface{}) ([]domain.ListItem, error) {
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listItems []domain.ListItem
	for rows.Next() {
		var listItem domain.ListItem
		err := rows.Scan(&listItem.ID, &listItem.ListID, &listItem.BookID)
		if err != nil {
			return nil, err
		}
		listItems = append(listItems, listItem)
	}
	return listItems, nil
}

func (r *ListItemRepository) GetListItemsByListID(ctx context.Context, listID int64) ([]domain.ListItem, error) {
	query := "SELECT id, list_id, book_id FROM list_item WHERE list_id = ?"
	return r.getAll(query, listID)
}

func (r *ListItemRepository) CreateListItem(ctx context.Context, listItem domain.ListItem) (domain.ListItem, error) {
	query := "INSERT INTO list_item (list_id, book_id) VALUES (?, ?)"
	res, err := r.DB.Exec(query, listItem.ListID, listItem.BookID)
	if err != nil {
		return domain.ListItem{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return domain.ListItem{}, err
	}
	listItem.ID = id
	return listItem, nil
}

func (r *ListItemRepository) DeleteListItem(ctx context.Context, id int64) error {
	query := "DELETE FROM list_item WHERE id = ?"
	_, err := r.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
