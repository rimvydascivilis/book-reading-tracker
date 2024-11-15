package domain

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email" validate:"required,email,max=255"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"required,min=1,max=50"`
	Rating    float64   `json:"rating,omitempty" validate:"gte=0,lte=5"`
	CreatedAt time.Time `json:"created_at"`
}

type Goal struct {
	UserID    int64  `json:"user_id" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=books pages"`
	Frequency string `json:"frequency" validate:"required,oneof=daily monthly"`
	Value     int64  `json:"value" validate:"required,min=1"`
}
