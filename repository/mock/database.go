package mock

import (

	"github.com/library/domain/model"



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
