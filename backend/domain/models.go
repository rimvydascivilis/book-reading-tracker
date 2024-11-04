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
