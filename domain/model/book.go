package model

import "time"

type Book struct {

	Id            int
	Title         string
	Author        string
	Isbn          string
	Isbn13        string	`gorm:"column:isbn_13"`
	OpenLibraryId string
	CoverId       string
	Year          int
	Publisher     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}