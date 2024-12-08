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
	UserID    int64     `json:"user_id" validate:"required"`
	Title     string    `json:"title" validate:"required,min=1,max=50"`
	Rating    float64   `json:"rating,omitempty" validate:"gte=0,lte=5"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	GoalTypeBooks        = "books"
	GoalTypePages        = "pages"
	GoalFrequencyDaily   = "daily"
	GoalFrequencyMonthly = "monthly"
)

type Goal struct {
	UserID    int64  `json:"user_id" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=books pages"`
	Frequency string `json:"frequency" validate:"required,oneof=daily monthly"`
	Value     int64  `json:"value" validate:"required,min=1"`
}

var (
	ReadingStatusNotStarted = "not started"
	ReadingStatusReading    = "reading"
	ReadingStatusCompleted  = "completed"
)

type Reading struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id" validate:"required"`
	BookID     int64     `json:"book_id" validate:"required"`
	TotalPages int64     `json:"total_pages" validate:"required,min=1"`
	Link       string    `json:"link,omitempty" validate:"omitempty,url"`
	CreatedAt  time.Time `json:"created_at" validate:"required"`
	UpdatedAt  time.Time `json:"updated_at" validate:"required"`
}

func (r *Reading) GetStatus(progress int64) string {
	if progress == 0 {
		return ReadingStatusNotStarted
	} else if progress < r.TotalPages {
		return ReadingStatusReading
	}

	return ReadingStatusCompleted
}

type Progress struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id" validate:"required"`
	ReadingID   int64     `json:"reading_id" validate:"required"`
	Pages       int64     `json:"pages" validate:"required,min=1"`
	ReadingDate time.Time `json:"reading_date" validate:"required"`
}

type List struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id" validate:"required"`
	Title     string    `json:"title" validate:"required,min=1,max=50"`
	CreatedAt time.Time `json:"created_at"`
}

type ListItem struct {
	ID        int64     `json:"id"`
	ListID    int64     `json:"list_id" validate:"required"`
	BookID    int64     `json:"book_id" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
}

type Note struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id" validate:"required"`
	BookID     int64     `json:"book_id" validate:"required"`
	PageNumber int64     `json:"page_number" validate:"required,min=1"`
	Content    string    `json:"content" validate:"required"`
	CreatedAt  time.Time `json:"created_at"`
}
