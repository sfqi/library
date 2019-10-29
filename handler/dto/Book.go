package dto

import (
	"fmt"
	"github.com/library/domain/model"
	"github.com/library/repository/mock"
	"strings"
)

type Book struct {
	Title      string     `json:"title"`
	Identifier identifier `json:"identifiers"`
	Author     []author   `json:"authors"`
	Cover      cover      `json:"cover"`
	Year       string     `json:"publish_date"`
}

type identifier struct {
	ISBN10      []string `json:"isbn_10"`
	ISBN13      []string `json:"isbn_13"`
	Openlibrary []string `json:"openlibrary"`
}
type author struct {
	Name string `json:"name"`
}

type cover struct {
	Url string `json:"small"`
}


func CreateBookModelFromBook(b Book) (bm model.Book) {

	mock.Shelf.Id = mock.Shelf.Id + 1

	isbn10 := ""
	if b.Identifier.ISBN10 != nil {
		isbn10 = b.Identifier.ISBN10[0]
	}
	isbn13 := ""
	if b.Identifier.ISBN13 != nil {
		isbn13 = b.Identifier.ISBN13[0]
	}
	fmt.Println(mock.Shelf.Id, isbn10, isbn13)
	CoverId := ""
	if b.Cover.Url != "" {
		part1 := strings.Split(b.Cover.Url, "/")[5]
		part2 := strings.Split(part1, ".")[0]
		CoverId = strings.Split(part2, "-")[0]
	}
	libraryId := ""
	if b.Identifier.Openlibrary != nil {
		libraryId = b.Identifier.Openlibrary[0]
	}

	bookToAdd := model.Book{
		Id:            mock.Shelf.Id,
		Title:         b.Title,
		Author:        b.Author[0].Name,
		Isbn:          isbn10,
		Isbn13:        isbn13,
		OpenLibraryId: libraryId,
		CoverId:       CoverId,
		Year:          b.Year,
	}
	return bookToAdd
}