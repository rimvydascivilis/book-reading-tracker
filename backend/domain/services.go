package domain

import (
	"context"

	"github.com/rimvydascivilis/book-tracker/backend/dto"
)

type ValidationService interface {
	ValidateStruct(s interface{}) error
}

type AuthService interface {
	Login(ctx context.Context, googleOauthToken string) (string, error)
}

type TokenService interface {
	GenerateToken(ctx context.Context, userID int64) (string, error)
}

type OAuth2Service interface {
	ValidateToken(token string) (string, error)
}

type BookService interface {
	CreateBook(ctx context.Context, userID int64, book Book) (Book, error)
	GetBooks(ctx context.Context, userID, page, limit int64) ([]Book, bool, error)
	SearchBooks(ctx context.Context, userID int64, title string, limit int64) ([]Book, error)
	UpdateBook(ctx context.Context, userID int64, book Book) (Book, error)
	DeleteBook(ctx context.Context, userID, bookID int64) error
}

type UserService interface {
	GetOrCreateUser(ctx context.Context, email string) (User, error)
}

type GoalService interface {
	GetGoal(ctx context.Context, userID int64) (Goal, error)
	GetGoalProgress(ctx context.Context, userID int64) (dto.GoalProgressResponse, error)
	SetGoal(ctx context.Context, userID int64, goal Goal) (Goal, error)
}

type ReadingService interface {
	GetReadings(ctx context.Context, userID, page, limit int64) ([]dto.ReadingResponse, bool, error)
	CreateReading(ctx context.Context, userID int64, reading Reading) (Reading, error)
}

type ProgressService interface {
	CreateProgress(ctx context.Context, userID, readingID int64, progressReq dto.ProgressRequest) (Progress, error)
}

type ListService interface {
	ListLists(ctx context.Context, userID int64) ([]dto.ListListsResponse, error)
	GetList(ctx context.Context, userID int64, listID int64) (dto.ListResponse, error)
	CreateList(ctx context.Context, userID int64, list dto.ListRequest) (dto.ListResponse, error)
	AddBookToList(ctx context.Context, userID, listID, bookID int64) error
	RemoveBookFromList(ctx context.Context, userID, listID, bookID int64) error
}

type NoteService interface {
	GetNotes(ctx context.Context, userID, bookID int64) ([]dto.NoteResponse, error)
	CreateNote(ctx context.Context, userID, bookID int64, note dto.NoteRequest) (dto.NoteResponse, error)
	DeleteNote(ctx context.Context, userID, noteID int64) error
}
