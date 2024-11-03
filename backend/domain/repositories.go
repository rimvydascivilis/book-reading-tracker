package domain

import "context"

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	CreateUser(ctx context.Context, u User) (User, error)
}
