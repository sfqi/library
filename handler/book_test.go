package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler"
	"github.com/sfqi/library/handler/dto"
	imock "github.com/sfqi/library/interactor/mock"
	"github.com/stretchr/testify/assert"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Successfully returned books", func(t *testing.T) {
		interactor := &imock.Book{}
		bookHandler := handler.BookHandler{}

		booksExpected := []dto.BookResponse{
			{
				ID:            1,
				Title:         "some title",
				Author:        "some author",
				Isbn:          "some isbn",
				Isbn13:        "some isbon13",
				OpenLibraryId: "again some id",
				CoverId:       "some cover ID",
				Year:          2019,
			},
			{
				ID:            2,
				Title:         "other title",
				Author:        "other author",
				Isbn:          "other isbn",
				Isbn13:        "other isbon13",
				OpenLibraryId: "other some id",
				CoverId:       "other cover ID",
				Year:          2019,
			},
		}

		req, err := http.NewRequest("GET", "/books", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		interactor.On("FindAll").Return([]*model.Book{
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

		bookHandler.Interactor = interactor
		httpError := bookHandler.Index(rr, req)
		assert.Nil(httpError)
		assert.Equal(http.StatusOK, rr.Code)

		var response []dto.BookResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(err)

		assert.Equal(booksExpected, response)
	})
	t.Run("Error retrieving books", func(t *testing.T) {
		interactor := &imock.Book{}
		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("GET", "/books", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		interactor.On("FindAll").Return(nil, errors.New("Error finding books"))
		bookHandler.Interactor = interactor

		handler := bookHandler.Index
		httperror := handler(rr, req)

		expectedResponse := "HTTP 500: Error finding books"
		assert.Equal(expectedResponse, httperror.Error())

		assert.Equal(http.StatusInternalServerError, httperror.Code())

	})

}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("error getting book from context", func(t *testing.T) {
		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"title":"test title", "year":2019}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)

		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", 5)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		httpError := bookHandler.Update(rr, req)

		expectedError := "HTTP 500: error retrieving book from context"
		assert.Equal(expectedError, httpError.Error())
	})
	t.Run("error updating book in database", func(t *testing.T) {
		interactor := &imock.Book{}
		bookHandler := handler.BookHandler{}
		book := &model.Book{
			Id:            2,
			Title:         "some title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2000,
		}
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"title":"test title", "year":2019}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)

		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		interactor.On("Update", &model.Book{
			Id:            2,
			Title:         "some title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2000,
		}, dto.UpdateBookRequest{
			Title: "test title",
			Year:  2019,
		},
		).Return(nil, errors.New("error updating book"))
		bookHandler.Interactor = interactor

		httpError := bookHandler.Update(rr, req)

		assert.Equal(http.StatusInternalServerError, httpError.Code())

		expectedError := "HTTP 500: error updating book"
		assert.Equal(expectedError, httpError.Error())
	})
	t.Run("assertion of expected response, and actual response", func(t *testing.T) {
		interactor := &imock.Book{}
		bookHandler := handler.BookHandler{}
		book := &model.Book{
			Id:            2,
			Title:         "some title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2000,
		}
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"title":"test title", "year":2019}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)

		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		interactor.On("Update", &model.Book{
			Id:            2,
			Title:         "some title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2000,
		}, dto.UpdateBookRequest{
			Title: "test title",
			Year:  2019,
		}).Return(&model.Book{
			Id:            2,
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2019,
		}, nil)
		bookHandler.Interactor = interactor

		httpError := bookHandler.Update(rr, req)
		assert.Nil(httpError)

		require.Equal(http.StatusOK, rr.Code)

		bookExpected := dto.BookResponse{

			ID:            2,
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2019,
		}
		var response dto.BookResponse
		err = json.NewDecoder(rr.Body).Decode(&response)

		require.NoError(err)

		assert.Equal(bookExpected, response)
	})
	t.Run("Error decoding Book attributes", func(t *testing.T) {
		bookHandler := handler.BookHandler{}
		book := &model.Book{
			Id:            2,
			Title:         "some title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2000,
		}
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"id":"12","title":zdravo}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)

		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		httpError := bookHandler.Update(rr, req)

		expectedError := "HTTP 400: invalid character 'z' looking for beginning of value"

		assert.Equal(expectedError, httpError.Error())
	})

}

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Invalid request body", func(t *testing.T) {

		bookHandler := handler.BookHandler{}
		interactor := &imock.Book{}
		bookHandler.Interactor = interactor

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{ISBN:"0140447938"}`)))

		require.NoError(err)

		rr := httptest.NewRecorder()

		httpError := bookHandler.Create(rr, req)
		expectedError := "HTTP 400: invalid character 'I' looking for beginning of object key string"

		assert.Equal(expectedError, httpError.Error())
	})
	t.Run("Fetching book error", func(t *testing.T) {
		bookHandler := handler.BookHandler{}
		interactor := &imock.Book{}

		interactor.On("Create", dto.CreateBookRequest{
			ISBN: "0140447938222",
		}).Return(nil, errors.New("Error while fetching book: "))
		bookHandler.Interactor = interactor

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938222"}`)))
		require.NoError(err)

		rr := httptest.NewRecorder()

		httpError := bookHandler.Create(rr, req)

		require.Contains(httpError.Error(), "Error while fetching book: ")
		assert.Equal(http.StatusInternalServerError, httpError.Code())
	})
	t.Run("Creation of book failed in database", func(t *testing.T) {

		interactor := &imock.Book{}
		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938"}`)))

		require.NoError(err)

		rr := httptest.NewRecorder()
		interactor.On("Create", dto.CreateBookRequest{
			ISBN: "0140447938",
		}).Return(nil, errors.New("Error creating book"))

		bookHandler.Interactor = interactor

		httpError := bookHandler.Create(rr, req)

		expectedResponse := "HTTP 500: Error creating book"
		assert.Equal(expectedResponse, httpError.Error())
		assert.Equal(http.StatusInternalServerError, httpError.Code())
	})
	t.Run("Testing book creation", func(t *testing.T) {
		bookHandler := handler.BookHandler{}
		interactor := &imock.Book{}

		bookHandler.Interactor = interactor

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938"}`)))

		rr := httptest.NewRecorder()
		interactor.On("Create", dto.CreateBookRequest{
			ISBN: "0140447938",
		}).Return(&model.Book{
			Id:            0,
			Title:         "War and Peace (Penguin Classics)",
			Author:        "Tolstoy",
			Isbn:          "0140447938",
			Isbn13:        "9780140447934",
			OpenLibraryId: "OL7355422M",
			CoverId:       "5049015",
			Year:          2007,
		}, nil)

		httpError := bookHandler.Create(rr, req)
		assert.Nil(httpError)
		bookExpected := dto.BookResponse{
			0,
			"War and Peace (Penguin Classics)",
			"Tolstoy",
			"0140447938",
			"9780140447934",
			"OL7355422M",
			"5049015",
			2007,
		}
		var response dto.BookResponse
		err = json.NewDecoder(rr.Body).Decode(&response)

		require.NoError(err)

		assert.Equal(bookExpected, response)
	})
}

func TestGet(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Successfully retrieved book", func(t *testing.T) {
		interactor := &imock.Book{}
		bookHandler := handler.BookHandler{}
		book := &model.Book{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		}
		req, err := http.NewRequest("GET", "/book/{id}", nil)
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)

		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		interactor.On("FindBookById", 1).Return(&model.Book{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		}, nil)

		bookHandler.Interactor = interactor

		httpError := bookHandler.Get(rr, req)
		assert.Nil(httpError)
		expectedBook := dto.BookResponse{
			ID:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		}
		var response dto.BookResponse
		err = json.NewDecoder(rr.Body).Decode(&response)

		require.NoError(err)

		assert.Equal(expectedBook, response)
	})
	t.Run("Cannot retreive book from context", func(t *testing.T) {
		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("GET", "/book/{id}", nil)
		params := map[string]string{"id": "5"}
		req = mux.SetURLVars(req, params)

		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", nil)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		httpError := bookHandler.Get(rr, req)
		expectedResponse := "HTTP 500: error retrieving book from context"

		assert.Equal(expectedResponse, httpError.Error())
	})
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Book successfully deleted", func(t *testing.T) {
		interactor := &imock.Book{}
		book := &model.Book{
			Id:            2,
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2019,
		}

		bookHandler := handler.BookHandler{}
		req, err := http.NewRequest("DELETE", "/books/{id}", nil)
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)

		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		interactor.On("Delete", book).Return(nil)

		bookHandler.Interactor = interactor

		httpError := bookHandler.Delete(rr, req)
		assert.Nil(httpError)

		require.Equal(http.StatusNoContent, rr.Code)
	})
	t.Run("Error retrieving book from context", func(t *testing.T) {
		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("DELETE", "/books/{id}", nil)
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)

		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", nil)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := bookHandler.Delete
		httpError := handler(rr, req)

		expectedResponse := "HTTP 500: error retrieving book from context"

		require.Equal(http.StatusInternalServerError, httpError.Code())

		assert.Equal(expectedResponse, httpError.Error())
	})
	t.Run("Error deleting book", func(t *testing.T) {
		bookHandler := handler.BookHandler{}
		book := &model.Book{
			Id:            2,
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2019,
		}
		interactor := &imock.Book{}
		req, err := http.NewRequest("DELETE", "/books/{id}", nil)
		params := map[string]string{"id": "-5"}
		req = mux.SetURLVars(req, params)

		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		interactor.On("Delete", book).Return(errors.New("Error while deleting book"))
		bookHandler.Interactor = interactor

		httpError := bookHandler.Delete(rr, req)

		expectedResponse := "HTTP 500: Error while deleting book"
		assert.Equal(expectedResponse, httpError.Error())
		assert.Equal(http.StatusInternalServerError, httpError.Code())
	})
}
