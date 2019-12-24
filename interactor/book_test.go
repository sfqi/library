package interactor_test

import (
	"errors"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
	"github.com/sfqi/library/interactor"
	openlibrarydto "github.com/sfqi/library/openlibrary/dto"
	olmock "github.com/sfqi/library/openlibrary/mock"
	"github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBook_Create(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Error while fetching book", func(t *testing.T) {
		var db = &mock.Store{}
		clmock := &olmock.Client{}

		clmock.On("FetchBook", "0140447938222").Return(nil, errors.New("Error while fetching book: "))
		book := interactor.NewBook(db, clmock)
		bookRequest := dto.CreateBookRequest{ISBN: "0140447938222"}
		bookResponse, err := book.Create(bookRequest)
		assert.Nil(bookResponse)
		require.Contains(err.Error(), "Error while fetching book: ")
	})
	t.Run("Error creating book in database", func(t *testing.T) {
		var db = &mock.Store{}
		clmock := &olmock.Client{}
		clmock.On("FetchBook", "0140447938").Return(&openlibrarydto.Book{Title: "War and Peace (Penguin Classics)"}, nil)
		db.On("CreateBook", &model.Book{Title: "War and Peace (Penguin Classics)"}).Return(errors.New("Error saving book in database"))
		book := interactor.NewBook(db, clmock)
		bookRequest := dto.CreateBookRequest{ISBN: "0140447938"}
		bookResponse, err := book.Create(bookRequest)
		assert.Nil(bookResponse)
		require.Contains(err.Error(), "Error saving book in database")
	})
	t.Run("Book successfully saved in database", func(t *testing.T) {
		var db = &mock.Store{}
		clmock := &olmock.Client{}
		clmock.On("FetchBook", "0140447938").Return(&openlibrarydto.Book{Title: "War and Peace (Penguin Classics)"}, nil)
		db.On("CreateBook", &model.Book{Title: "War and Peace (Penguin Classics)"}).Return(nil)
		book := interactor.NewBook(db, clmock)
		bookRequest := dto.CreateBookRequest{ISBN: "0140447938"}
		bookResponse, err := book.Create(bookRequest)
		assert.NotNil(bookResponse)
		assert.Nil(err)
	})

}

func TestBook_FindAll(t *testing.T) {

}

func TestBook_FindById(t *testing.T) {

}

func TestBook_Delete(t *testing.T) {

}

func TestBook_Update(t *testing.T) {

}
