package model

import "time"

type Book struct {
	Id            int
	Title         string
	Author        string
	Isbn          string    `gorm:"colum:isbn_10"`
	Isbn13        string    `gorm:"column:isbn_13"`
	OpenLibraryId string    `gorm:"column:open_library_id"`
	CoverId       string    `gorm:"column:cover_id"`
	Year          string    `gorm:"column:publish_date"`
	CreatedAt     time.Time `gorm:"column:created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at"`
}
