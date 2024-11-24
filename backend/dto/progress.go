package dto

import "time"

type ProgressRequest struct {
	Pages int64     `json:"pages"`
	Date  time.Time `json:"date"`
}
