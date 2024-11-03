package domain

import (
	"time"
)

type User struct {
	ID        int64
	Email     string
	CreatedAt time.Time
}

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Rating    *float64  `json:"rating,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
