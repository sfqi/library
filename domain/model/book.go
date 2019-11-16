package model

import "time"

type Book struct {
	Id            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Isbn          string    `json:"isbn_10"`
	Isbn13        string    `json:"isbn_13" gorm:"column:isbn_13"`
	OpenLibraryId string    `json:"olid"`
	CoverId       string    `json:"cover"`
	Year          string    `json:"year"`
	CreatedAt     time.Time `json:created_at`
	UpdatedAt     time.Time `json:updated_at`
}
