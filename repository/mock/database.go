package mock

import (
	"fmt"
	"github.com/library/domain/model"
	"github.com/library/handler/dto"
	"strings"

)

var Books = []model.Book{
	{
		Id:            1,
		Title:         "some title",
		Author:        "some author",
		Isbn:          "some isbn",
		Isbn13:        "some isbon13",
		OpenLibraryId: "again some id",
		CoverId:       "some cover ID",
		Year:          "2019",
	},
	{
		Id:            2,
		Title:         "other title",
		Author:        "other author",
		Isbn:          "other isbn",
		Isbn13:        "other isbon13",
		OpenLibraryId: "other some id",
		CoverId:       "other cover ID",
		Year:          "2019",
	},
	{
		Id:            3,
		Title:         "another title",
		Author:        "another author",
		Isbn:          "another isbn",
		Isbn13:        "another isbon13",
		OpenLibraryId: "another some id",
		CoverId:       "another cover ID",
		Year:          "2019",
	},
}

type Db struct {
	Id    int
	Books []model.Book
}

var Shelf = &Db{
	Id:    len(Books),
	Books: Books,
}

func (bm Db) FindBookById(id int) (book model.Book, location int, found bool) {
	for i, b := range bm.Books {
		if b.Id == id {
			book = b
			location = i
			found = true
			break
		}
	}

	return book, location, found
}


func CreateBookModelFromBook(b dto.Book) (bm model.Book) {

	Shelf.Id = Shelf.Id + 1

	isbn10 := ""
	if b.Identifier.ISBN10 != nil {
		isbn10 = b.Identifier.ISBN10[0]
	}
	isbn13 := ""
	if b.Identifier.ISBN13 != nil {
		isbn13 = b.Identifier.ISBN13[0]
	}
	fmt.Println(Shelf.Id, isbn10, isbn13)
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
		Id:            Shelf.Id,
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