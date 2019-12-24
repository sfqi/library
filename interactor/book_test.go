package interactor_test

import (
	"testing"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/repository/mock"
)

func TestBook_Create(t *testing.T) {

}

func TestBook_FindAll(t *testing.T) {

}

func TestFindById(t *testing.T) {
	t.Run("Successfully retrieved book", func(t *testing.T) {
		var db = &mock.Store{}
		db.On("FindById", 1).Return(&model.Book{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		}, nil)
	})
}

func TestBook_Delete(t *testing.T) {

}

func TestBook_Update(t *testing.T) {

}
