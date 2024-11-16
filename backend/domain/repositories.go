package domain

import "context"

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, u User) (User, error)
}

type BookRepository interface {
	CountBooksByUser(ctx context.Context, userID int64) (int64, error)
	GetBooksByUser(ctx context.Context, userID, offset, limit int64) ([]Book, error)
	GetBookByUserID(ctx context.Context, userID, bookID int64) (Book, error)
	UpdateBook(ctx context.Context, book Book) (Book, error)
	CreateBook(ctx context.Context, book Book) (Book, error)
	DeleteBook(ctx context.Context, userID, bookID int64) error
}

type GoalRepository interface {
	GetGoalByUserID(ctx context.Context, userID int64) (Goal, error)
	CreateGoal(ctx context.Context, goal Goal) (Goal, error)
	UpdateGoal(ctx context.Context, goal Goal) (Goal, error)
}
