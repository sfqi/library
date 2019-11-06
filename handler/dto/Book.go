package dto

type CreateBookRequest struct {
	ISBN string `json:"ISBN"`
}

type UpdateBookRequest struct {
	Title string `json:"title"`
	Year  string `json:"year"`
}

type BookResponse struct {
	ID            int    `json:id`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Isbn          string `json:"isbn_10"`
	Isbn13        string `json:"isbn_13"`
	OpenLibraryId string `json:"olid"`
	CoverId       string `json:"cover"`
	Year          string `json:"year"`
}
