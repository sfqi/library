package dto

type CreateBookRequest struct {
	ISBN string `json:"ISBN"`
}

type UpdateBookRequest struct {
	Title string `json:"title"`
	Year  string `json:"year"`
}

type BookResponse struct {
	ID    int    `json:id`
	Title string `json:"title"`
	Year  string `json:"year"`
}
