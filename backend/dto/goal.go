package dto

type GoalProgressResponse struct {
	Percentage float64 `json:"percentage"`
	Left       int64   `json:"left"`
}
