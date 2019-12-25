package interactor_test

import (
	"errors"
	"testing"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
	"github.com/sfqi/library/interactor"
	repomock "github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
)

func TestBook_Create(t *testing.T) {

}

func TestBook_FindAll(t *testing.T) {

}

func TestBook_FindById(t *testing.T) {

}

func TestBook_Delete(t *testing.T) {

}

func TestBook_Update(t *testing.T) {
	assert := assert.New(t)
	t.Run("Error updating book in database", func(t *testing.T) {
		store := &repomock.Store{}
		book := &model.Book{
			Title: "old title",
			Year:  2000,
		}
		storeError := "Error updating book"
		request := dto.UpdateBookRequest{}

		store.On("UpdateBook", book).Return(errors.New("Error updating book"))
		bookInteractor := interactor.NewBook(store, nil)
		modifiedBook, err := bookInteractor.Update(book, request)

		assert.Nil(modifiedBook)
		assert.Equal(err.Error(), storeError)
	})
	t.Run("Book successfully updated in database", func(t *testing.T) {
		var db = &repomock.Store{}
		book := &model.Book{
			Title: "old title",
			Year:  2000,
		}
		request := dto.UpdateBookRequest{
			Title: "New Title",
			Year:  2019,
		}
		expectedBook := &model.Book{
			Title: "New Title",
			Year:  2019,
		}

		db.On("UpdateBook", &model.Book{Title: "New Title", Year: 2019}).Return(nil)
		bookInteractor := interactor.NewBook(db, nil)
		modifiedBook, err := bookInteractor.Update(book, request)

		assert.Nil(err)
		assert.Equal(expectedBook, modifiedBook)
	})
}
