package interactor_test

import (
	"errors"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
	"github.com/sfqi/library/interactor"
	"github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
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
	require := require.New(t)
	t.Run("Error updating book in database", func(t *testing.T) {
		var db = &mock.Store{}
		bookToChange := &model.Book{
			Title: "old title",
			Year:  2000,
		}
		updateBook := dto.UpdateBookRequest{}
		db.On("UpdateBook", bookToChange).Return(errors.New("Error updating book"))
		book := interactor.NewBook(db, nil)
		modifiedBook, err := book.Update(bookToChange, updateBook)
		assert.Nil(modifiedBook)
		require.Contains(err.Error(), "Error updating book")
	})
	t.Run("Book successfully updated in database", func(t *testing.T) {
		var db = &mock.Store{}
		bookToChange := &model.Book{
			Title: "old title",
			Year:  2000,
		}
		updateBook := dto.UpdateBookRequest{
			Title: "New Title",
			Year:  2019,
		}
		db.On("UpdateBook", &model.Book{Title: "New Title", Year: 2019}).Return(nil)
		book := interactor.NewBook(db, nil)
		modifiedBook, err := book.Update(bookToChange, updateBook)
		assert.Nil(err)
		assert.NotNil(modifiedBook)
	})
}
