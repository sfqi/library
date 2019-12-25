package interactor_test

import (
	"errors"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
	"github.com/sfqi/library/interactor"
	openlibrarydto "github.com/sfqi/library/openlibrary/dto"
	openlibmock "github.com/sfqi/library/openlibrary/mock"
	repomock "github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBook_Create(t *testing.T) {
	assert := assert.New(t)
	//require := require.New(t)
	t.Run("Error while fetching book", func(t *testing.T) {
		store := &repomock.Store{}
		openlibClient := &openlibmock.Client{}
		clientError := "Error while fetching book: "

		openlibClient.On("FetchBook", "0140447938222").Return(nil, errors.New("Error while fetching book: "))
		bookInteractor := interactor.NewBook(store, openlibClient)
		request := dto.CreateBookRequest{ISBN: "0140447938222"}
		book, err := bookInteractor.Create(request)

		assert.Nil(book)
		assert.Equal(err.Error(), clientError)
	})
	t.Run("Error creating book in database", func(t *testing.T) {
		store := &repomock.Store{}
		openlibClient := &openlibmock.Client{}
		storeError := "Error saving book in database"

		openlibClient.On("FetchBook", "0140447938").Return(&openlibrarydto.Book{Title: "War and Peace (Penguin Classics)"}, nil)
		store.On("CreateBook", &model.Book{Title: "War and Peace (Penguin Classics)"}).Return(errors.New("Error saving book in database"))
		bookInteractor := interactor.NewBook(store, openlibClient)
		request := dto.CreateBookRequest{ISBN: "0140447938"}
		book, err := bookInteractor.Create(request)

		assert.Nil(book)
		assert.Equal(err.Error(), storeError)
	})
	t.Run("Book successfully saved in database", func(t *testing.T) {
		store := &repomock.Store{}
		openlibClient := &openlibmock.Client{}
		bookExpected := &model.Book{
			Id:            0,
			Title:         "War and Peace (Penguin Classics)",
			Author:        "Tolstoy",
			Isbn:          "0140447938",
			Isbn13:        "9780140447934",
			OpenLibraryId: "OL7355422M",
			CoverId:       "5049015",
			Year:          2007,
		}

		openlibClient.On("FetchBook", "0140447938").Return(
			&openlibrarydto.Book{
				Title: "War and Peace (Penguin Classics)",
				Identifier: openlibrarydto.Identifier{
					ISBN10:      []string{"0140447938"},
					ISBN13:      []string{"9780140447934"},
					Openlibrary: []string{"OL7355422M"},
				},
				Author: []openlibrarydto.Author{
					{Name: "Tolstoy"},
				},
				Cover: openlibrarydto.Cover{"https://covers.openlibrary.org/b/id/5049015-S.jpg"},
				Year:  "2007",
			},
			nil)
		store.On("CreateBook", &model.Book{
			Id:            0,
			Title:         "War and Peace (Penguin Classics)",
			Author:        "Tolstoy",
			Isbn:          "0140447938",
			Isbn13:        "9780140447934",
			OpenLibraryId: "OL7355422M",
			CoverId:       "5049015",
			Year:          2007,
		}).Return(nil)

		bookInteractor := interactor.NewBook(store, openlibClient)
		request := dto.CreateBookRequest{ISBN: "0140447938"}
		book, err := bookInteractor.Create(request)

		assert.Equal(bookExpected, book)
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
