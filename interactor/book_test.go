package interactor_test

import (
	"errors"
	"testing"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/interactor"
	"github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
)

func TestBook_Create(t *testing.T) {

}

func TestFindAll(t *testing.T) {
	t.Run("Successfully returned books", func(t *testing.T) {
		var db = &mock.Store{}

		b := interactor.NewBook(db, nil)

		db.On("FindAllBooks").Return([]*model.Book{
			{
				Id:            1,
				Title:         "some title",
				Author:        "some author",
				Isbn:          "some isbn",
				Isbn13:        "some isbon13",
				OpenLibraryId: "again some id",
				CoverId:       "some cover ID",
				Year:          2019,
			},
			{
				Id:            2,
				Title:         "other title",
				Author:        "other author",
				Isbn:          "other isbn",
				Isbn13:        "other isbon13",
				OpenLibraryId: "other some id",
				CoverId:       "other cover ID",
				Year:          2019,
			}}, nil)

		books, err := b.FindAll()
		assert.NotNil(t, books)
		assert.NoError(t, err)
	})
	t.Run("Error retrieving books", func(t *testing.T) {
		var db = &mock.Store{}

		b := interactor.NewBook(db, nil)

		db.On("FindAllBooks").Return(nil, errors.New("Error finding books"))

		book, err := b.FindAll()

		assert.Nil(t, book)
		assert.Error(t, err)
	})

}

func TestBook_FindById(t *testing.T) {

}

func TestBook_Delete(t *testing.T) {

}

func TestBook_Update(t *testing.T) {

}
