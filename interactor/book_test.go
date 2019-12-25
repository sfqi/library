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
	assert := assert.New(t)
	t.Run("Successfully returned books", func(t *testing.T) {
		store := &mock.Store{}

		b := interactor.NewBook(store, nil)

		store.On("FindAllBooks").Return([]*model.Book{
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

		expectedBooks := []*model.Book{
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
			},
		}

		books, err := b.FindAll()
		assert.Equal(books, expectedBooks)
		assert.NoError(err)
	})
	t.Run("Error retrieving books", func(t *testing.T) {
		store := &mock.Store{}
		storeError := errors.New("Error finding books")

		b := interactor.NewBook(store, nil)

		store.On("FindAllBooks").Return(nil, errors.New("Error finding books"))

		book, err := b.FindAll()

		assert.Nil(book)
		assert.Error(err)
		assert.Equal(errors.New("Error finding books"), storeError)
	})

}

func TestBook_FindById(t *testing.T) {

}

func TestBook_Delete(t *testing.T) {

}

func TestBook_Update(t *testing.T) {

}
