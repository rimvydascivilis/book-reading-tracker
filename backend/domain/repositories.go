package domain

import (
	"context"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, u User) (User, error)
}

type BookRepository interface {
	CountBooksByUser(ctx context.Context, userID int64) (int64, error)
	GetBooksByUser(ctx context.Context, userID, offset, limit int64) ([]Book, error)
	GetBookByUserID(ctx context.Context, userID, bookID int64) (Book, error)
	SearchBooksByTitle(ctx context.Context, userID int64, title string, limit int64) ([]Book, error)
	UpdateBook(ctx context.Context, book Book) (Book, error)
	CreateBook(ctx context.Context, book Book) (Book, error)
	DeleteBook(ctx context.Context, userID, bookID int64) error
}

type GoalRepository interface {
	GetGoalByUserID(ctx context.Context, userID int64) (Goal, error)
	CreateGoal(ctx context.Context, goal Goal) (Goal, error)
	UpdateGoal(ctx context.Context, goal Goal) (Goal, error)
}

type ReadingRepository interface {
	GetReadingsByUserID(ctx context.Context, userID, offset, limit int64) ([]Reading, error)
	GetReadingByID(ctx context.Context, id int64) (Reading, error)
	CountReadingsByUserID(ctx context.Context, userID int64) (int64, error)
	CountReadingsByUserIDAndBookID(ctx context.Context, userID, bookID int64) (int64, error)
	CreateReading(ctx context.Context, reading Reading) (Reading, error)
}

type ProgressRepository interface {
	GetTotalProgressByReadingID(ctx context.Context, readingID int64) (int64, error)
	GetUserReadingIDsByPeriod(ctx context.Context, userID int64, period string) ([]int64, error)
	CreateProgress(ctx context.Context, progressReq Progress) (Progress, error)
}
