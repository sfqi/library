package dto

type Book struct {
	Title      string     `json:"title"`
	Identifier Identifier `json:"identifiers"`
	Author     []Author   `json:"authors"`
	Cover      Cover      `json:"cover"`
	Year       string     `json:"year"`
}

type Identifier struct {
	ISBN10      []string `json:"isbn_10"`
	ISBN13      []string `json:"isbn_13"`
	Openlibrary []string `json:"openlibrary"`
}
type Author struct {
	Name string `json:"name"`
}

type Cover struct {
	Url string `json:"small"`
}
