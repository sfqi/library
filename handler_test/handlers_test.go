package handler_test

import (
	"bytes"

	"encoding/json"
	"errors"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler"
	"github.com/sfqi/library/handler/dto"
	openlibrarydto "github.com/sfqi/library/openlibrary/dto"
	olmock "github.com/sfqi/library/openlibrary/mock"
	"github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	t.Run("Successfully returned books", func(t *testing.T) {
		var db = &mock.Store{}
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
		db.On("FindAllBooks").Return([]*model.Book{
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

		bookHandler.Db = db
		handler := http.HandlerFunc(bookHandler.Index)
		handler.ServeHTTP(rr, req)
		require.Equal(http.StatusOK, rr.Code)

		var response []dto.BookResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(err)

		assert.Equal(booksExpected, response)
	})
	t.Run("Error retrieving books", func(t *testing.T) {
		var db = &mock.Store{}
		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("GET", "/books", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		db.On("FindAllBooks").Return(nil, errors.New("Error finding books"))

		bookHandler.Db = db

		handler := http.HandlerFunc(bookHandler.Index)

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code)

		expectedResponse := "Error finding books\n"
		assert.Equal(expectedResponse, rr.Body.String())

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

		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)

		expectedError := "Internal server error:error retrieving book from context\n"
		assert.Equal(expectedError, rr.Body.String())
	})
	t.Run("error updating book in database", func(t *testing.T) {
		var db = &mock.Store{}
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
		db.On("UpdateBook", &model.Book{
			Id:            2,
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2019,
		}).Return(errors.New("error updating book"))
		bookHandler.Db = db
		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)

		assert.Equal(http.StatusInternalServerError, rr.Code)

		expectedError := "error updating book\n"
		assert.Equal(expectedError, rr.Body.String())
	})
	t.Run("assertion of expected response, and actual response", func(t *testing.T) {
		var db = &mock.Store{}
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
		db.On("UpdateBook", &model.Book{
			Id:            2,
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2019,
		}).Return(nil)
		bookHandler.Db = db
		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)

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
		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body\n"

		assert.Equal(expectedError, rr.Body.String())

	})
	t.Run("Cannot retreive book from context", func(t *testing.T) {

		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("UPDATE", "/book/{id}", nil)
		params := map[string]string{"id": "5"}
		req = mux.SetURLVars(req, params)
		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", nil)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Update)
		handler.ServeHTTP(rr, req)
		expectedResponse := "Internal server error:error retrieving book from context" + "\n"

		assert.Equal(expectedResponse, rr.Body.String())
	})

}

func TestCreate(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Invalid request body", func(t *testing.T) {

		bookHandler := handler.BookHandler{}
		clmock := olmock.Client{}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{ISBN:"0140447938"}`)))

		require.NoError(err)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Create)

		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body\n"

		assert.Equal(expectedError, rr.Body.String())
	})
	t.Run("Fetching book error", func(t *testing.T) {
		bookHandler := handler.BookHandler{}
		clmock := olmock.Client{}
		clmock.On("FetchBook", "0140447938222").Return(nil, errors.New("Error while fetching book: "))
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938222"}`)))
		require.NoError(err)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Create)
		handler.ServeHTTP(rr, req)

		require.Contains(rr.Body.String(), "Error while fetching book: ")
	})
	t.Run("Creation of book failed in database", func(t *testing.T) {
		db := mock.Store{}
		bookHandler := handler.BookHandler{}

		clmock := olmock.Client{}
		clmock.On("FetchBook", "0140447938").Return(&openlibrarydto.Book{Title: "War and Peace (Penguin Classics)"}, nil)
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938"}`)))

		require.NoError(err)

		rr := httptest.NewRecorder()
		db.On("CreateBook", &model.Book{
			Title: "War and Peace (Penguin Classics)",
		}).Return(errors.New("Error creating book"))
		bookHandler.Db = &db
		handler := http.HandlerFunc(bookHandler.Create)

		handler.ServeHTTP(rr, req)

		expectedResponse := "Error creating book\n"
		assert.Equal(expectedResponse, rr.Body.String())
	})
	t.Run("Testing book creation", func(t *testing.T) {
		db := mock.Store{}
		bookHandler := handler.BookHandler{}

		clmock := olmock.Client{}
		clmock.On("FetchBook", "0140447938").Return(&openlibrarydto.Book{
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
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938"}`)))

		require.NoError(err)

		rr := httptest.NewRecorder()
		db.On("CreateBook", &model.Book{
			Id:            0,
			Title:         "War and Peace (Penguin Classics)",
			Author:        "Tolstoy",
			Isbn:          "0140447938",
			Isbn13:        "9780140447934",
			OpenLibraryId: "OL7355422M",
			CoverId:       "5049015",
			Year:          2007,
		}).Return(nil)
		bookHandler.Db = &db
		handler := http.HandlerFunc(bookHandler.Create)

		handler.ServeHTTP(rr, req)

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
		var db = mock.Store{}
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

		bookHandler.Db = &db
		handler := http.HandlerFunc(bookHandler.Get)
		handler.ServeHTTP(rr, req)
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

		handler := http.HandlerFunc(bookHandler.Get)
		handler.ServeHTTP(rr, req)
		expectedResponse := "Internal server error:error retrieving book from context" + "\n"

		assert.Equal(expectedResponse, rr.Body.String())
	})
}

func TestDelete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Book successfully deleted", func(t *testing.T) {
		var db = mock.Store{}
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

		db.On("DeleteBook", book).Return(nil)

		bookHandler.Db = &db

		handler := http.HandlerFunc(bookHandler.Delete)

		handler.ServeHTTP(rr, req)

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

		handler := http.HandlerFunc(bookHandler.Delete)

		handler.ServeHTTP(rr, req)
		expectedResponse := "Internal server error:error retrieving book from context" + "\n"

		require.Equal(http.StatusInternalServerError, rr.Code)

		assert.Equal(expectedResponse, rr.Body.String())
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
		db := &mock.Store{}
		req, err := http.NewRequest("DELETE", "/books/{id}", nil)
		params := map[string]string{"id": "-5"}
		req = mux.SetURLVars(req, params)

		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		db.On("DeleteBook", book).Return(errors.New("Error while deleting book"))
		bookHandler.Db = db
		handler := http.HandlerFunc(bookHandler.Delete)

		handler.ServeHTTP(rr, req)
		expectedResponse := "Error while deleting book" + "\n"
		assert.Equal(expectedResponse, rr.Body.String())
	})
}
