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

func TestBook_FindAll(t *testing.T) {

}

func TestBook_FindById(t *testing.T) {

}

func TestDelete(t *testing.T) {
	t.Run("Book successfully deleted", func(t *testing.T) {
		var db = &mock.Store{}

		b := interactor.NewBook(db, nil)

		bookToDelete := &model.Book{}

		db.On("DeleteBook", bookToDelete).Return(nil)

		err := b.Delete(bookToDelete)

		assert.NoError(t, err)

	})
	t.Run("Error deleting book", func(t *testing.T) {
		var db = &mock.Store{}

		b := interactor.NewBook(db, nil)

		bookToDelete := &model.Book{}

		db.On("DeleteBook", bookToDelete).Return(errors.New("Error while deleting book"))

		err := b.Delete(bookToDelete)

		assert.Error(t, err)
	})
}

func TestBook_Update(t *testing.T) {

}
