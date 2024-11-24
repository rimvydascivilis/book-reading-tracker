package dto

type ReadingResponse struct {
	BookTitle string  `json:"book_title"`
	Status    string  `json:"status"`
	Progress  int64   `json:"progress"`
	Reading   Reading `json:"reading"`
}

type Reading struct {
	ID         int64  `json:"id"`
	TotalPages int64  `json:"total_pages"`
	Link       string `json:"link"`
}
