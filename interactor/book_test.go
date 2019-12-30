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
	t.Run("Error while fetching book", func(t *testing.T) {
		store := &repomock.Store{}
		openlibClient := &openlibmock.Client{}
		clientError := "Error while fetching book: "

		openlibClient.On("FetchBook", "0140447938222").Return(nil, errors.New(clientError))
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

		openlibClient.On("FetchBook", "0140447938").Return(&openlibrarydto.Book{Title: "War and Peace (Penguin Classics)", Year: "2007"}, nil)
		store.On("CreateBook", &model.Book{Title: "War and Peace (Penguin Classics)", Year: 2007}).Return(errors.New(storeError))
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
			nil,
		)
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

func TestFindAll(t *testing.T) {
	assert := assert.New(t)
	t.Run("Successfully returned books", func(t *testing.T) {
		store := &repomock.Store{}

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
		store := &repomock.Store{}
		storeError := errors.New("Error finding books")

		b := interactor.NewBook(store, nil)

		store.On("FindAllBooks").Return(nil, storeError)

		book, err := b.FindAll()

		assert.Nil(book)
		assert.Equal(err, storeError)
	})

}

func TestFindById(t *testing.T) {
	assert := assert.New(t)
	t.Run("Successfully retrieved book", func(t *testing.T) {
		store := &repomock.Store{}

		b := interactor.NewBook(store, nil)

		store.On("FindBookById", 1).Return(&model.Book{
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

		expectedBook := &model.Book{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		}

		assert.Equal(book, expectedBook)
		assert.NoError(err)

	})
	t.Run("Cannot retrieve book", func(t *testing.T) {
		store := &repomock.Store{}
		storeError := errors.New("Error finding ID from database")

		b := interactor.NewBook(store, nil)

		store.On("FindBookById", 12).Return(nil, storeError)

		book, err := b.FindById(12)

		assert.Nil(book)
		assert.Equal(err, storeError)
	})
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	t.Run("Book successfully deleted", func(t *testing.T) {
		store := &repomock.Store{}
		b := interactor.NewBook(store, nil)

		bookToDelete := &model.Book{}
		store.On("DeleteBook", bookToDelete).Return(nil)

		err := b.Delete(bookToDelete)
		assert.NoError(err)

	})
	t.Run("Error deleting book", func(t *testing.T) {
		store := &repomock.Store{}
		b := interactor.NewBook(store, nil)
		storeError := errors.New("Error while deleting book")

		bookToDelete := &model.Book{}
		store.On("DeleteBook", bookToDelete).Return(storeError)

		err := b.Delete(bookToDelete)
		assert.Equal(err, storeError)
	})
}

func TestBook_Update(t *testing.T) {

}
