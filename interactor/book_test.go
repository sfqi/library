package interactor_test

import (
	"errors"
	"testing"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/interactor"
	olmock "github.com/sfqi/library/openlibrary/mock"
	"github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
)

func TestBook_Create(t *testing.T) {

}

func TestBook_FindAll(t *testing.T) {

}

func TestFindById(t *testing.T) {
	t.Run("Successfully retrieved book", func(t *testing.T) {
		var db = &mock.Store{}
		clmock := &olmock.Client{}
		b := interactor.NewBook(db, clmock)

		db.On("FindBookById", 1).Return(&model.Book{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		}, nil)

		book, err := b.FindById(1)

		assert.NotNil(t, book)
		assert.NoError(t, err)

	})
	t.Run("Cannot retrieve book", func(t *testing.T) {
		var db = &mock.Store{}
		clmock := &olmock.Client{}
		b := interactor.NewBook(db, clmock)
		db.On("FindBookById", 12).Return(nil, errors.New("Error finding ID from database"))

		book, err := b.FindById(12)

		assert.Nil(t, book)
		assert.Error(t, err)
	})
}

func TestBook_Delete(t *testing.T) {

}

func TestBook_Update(t *testing.T) {

}
