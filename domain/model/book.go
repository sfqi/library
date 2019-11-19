package model

import "time"

type Book struct {
	Id            int
	Title         string
	Author        string
	Isbn          string
	Isbn13        string
	OpenLibraryId string
	CoverId       string
	Year          string    `gorm:"year"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
