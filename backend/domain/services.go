package domain

import "context"

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
	UpdateBook(ctx context.Context, userID int64, book Book) (Book, error)
	DeleteBook(ctx context.Context, userID, bookID int64) error
}

type UserService interface {
	GetOrCreateUser(ctx context.Context, email string) (User, error)
}
