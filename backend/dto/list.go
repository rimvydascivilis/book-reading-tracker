package dto

type ListRequest struct {
	Title string `json:"title"`
}

type ListListsResponse struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
}

type ListResponse struct {
	ID        int64               `json:"id"`
	Title     string              `json:"title"`
	ListItems []ListItemsResponse `json:"list_items"`
}

type ListItemsResponse struct {
	ID       int64  `json:"id"`
	ListID   int64  `json:"list_id"`
	BookName string `json:"book_name"`
}
