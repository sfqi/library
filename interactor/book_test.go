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
	assert := assert.New(t)
	t.Run("Book successfully deleted", func(t *testing.T) {
		store := &mock.Store{}
		b := interactor.NewBook(store, nil)

		bookToDelete := &model.Book{}
		store.On("DeleteBook", bookToDelete).Return(nil)

		err := b.Delete(bookToDelete)
		assert.NoError(err)

	})
	t.Run("Error deleting book", func(t *testing.T) {
		store := &mock.Store{}
		b := interactor.NewBook(store, nil)
		storeError := errors.New("Error while deleting book")

		bookToDelete := &model.Book{}
		store.On("DeleteBook", bookToDelete).Return(errors.New("Error while deleting book"))

		err := b.Delete(bookToDelete)
		assert.Error(err)
		assert.Equal(errors.New("Error while deleting book"), storeError)
	})
}

func TestBook_Update(t *testing.T) {

}
