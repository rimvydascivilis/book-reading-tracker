package dto

type NoteBook struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type BooksWithNotesResponse struct {
	Books []NoteBook `json:"books"`
}

type NoteRequest struct {
	PageNumber int64  `json:"page_number"`
	Content    string `json:"content"`
}

type NoteResponse struct {
	ID         int64  `json:"id"`
	PageNumber int64  `json:"page_number"`
	Content    string `json:"content"`
}
