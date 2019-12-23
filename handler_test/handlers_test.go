package handler_test

import (
	"bytes"

	"encoding/json"
	"errors"

	"net/http"
	"net/http/httptest"
	"strings"
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
		require.Equal(t, http.StatusOK, rr.Code)

		var response []dto.BookResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(t, err)

		assert.Equal(t, booksExpected, response, "Asserting expectation and actual response")
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

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		expectedResponse := "Error finding books\n"
		assert.Equal(t, expectedResponse, rr.Body.String(), "Asserting expectation and actual response")

	})

}

func TestUpdate(t *testing.T) {
	t.Run("error getting book from context", func(t *testing.T) {
		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"title":"test title", "year":2019}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)

		require.NoError(t, err)

		ctx := context.WithValue(req.Context(), "book", 5)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)

		expectedError := "Internal server error:error retrieving book from context\n"
		assert.Equal(t, expectedError, rr.Body.String(), "Response body differs")
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

		require.NoError(t, err)

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

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		expectedError := "error updating book\n"
		assert.Equal(t, expectedError, rr.Body.String(), "Response body differs")
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

		require.NoError(t, err)

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

		require.Equal(t, http.StatusOK, rr.Code)

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

		require.NoError(t, err, "Error decoding")

		assert.Equal(t, bookExpected, response, "Response body differs")
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

		require.NoError(t, err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body\n"

		assert.Equal(t, expectedError, rr.Body.String(), "Response body differs")

	})
	t.Run("Cannot retreive book from context", func(t *testing.T) {

		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("UPDATE", "/book/{id}", nil)
		params := map[string]string{"id": "5"}
		req = mux.SetURLVars(req, params)

		require.NoError(t, err)

		ctx := context.WithValue(req.Context(), "book", nil)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Update)
		handler.ServeHTTP(rr, req)
		expectedResponse := "Internal server error:error retrieving book from context" + "\n"

		assert.Equal(t, expectedResponse, rr.Body.String(), "Response body differs")
	})

}

func TestCreate(t *testing.T) {
	t.Run("Invalid request body", func(t *testing.T) {

		bookHandler := handler.BookHandler{}
		clmock := olmock.Client{
			Book: nil,
			Err:  errors.New("Error while decoding from request body"),
		}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{ISBN:"0140447938"}`)))

		require.NoError(t, err, "Error occured while sending request")

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Create)

		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body\n"

		assert.Equal(t, expectedError, rr.Body.String(), "Response body differs")
	})
	t.Run("Fetching book error", func(t *testing.T) {
		bookHandler := handler.BookHandler{}
		clmock := olmock.Client{nil,
			errors.New("Error while fetching book"),
		}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938222"}`)))

		require.NoError(t, err)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Create)

		handler.ServeHTTP(rr, req)
		contains := strings.Contains(rr.Body.String(), clmock.Err.Error())

		require.NotEqual(t, contains, rr.Code)

	})
	t.Run("Creation of book failed in database", func(t *testing.T) {
		db := mock.Store{}
		bookHandler := handler.BookHandler{}
		clmock := olmock.Client{&openlibrarydto.Book{
			Title: "War and Peace (Penguin Classics)",
		},
			nil,
		}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938"}`)))

		require.NoError(t, err)

		rr := httptest.NewRecorder()
		db.On("CreateBook", &model.Book{
			Title: "War and Peace (Penguin Classics)",
		}).Return(errors.New("Error creating book"))
		bookHandler.Db = &db
		handler := http.HandlerFunc(bookHandler.Create)

		handler.ServeHTTP(rr, req)

		expectedResponse := "Error creating book\n"
		assert.Equal(t, expectedResponse, rr.Body.String(), "Response body differs")
	})
	t.Run("Testing book creation", func(t *testing.T) {
		db := mock.Store{}
		bookHandler := handler.BookHandler{}
		clmock := olmock.Client{&openlibrarydto.Book{
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
		}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938"}`)))

		require.NoError(t, err)

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

		require.NoError(t, err, "Error decoding")

		assert.Equal(t, bookExpected, response, "Response body differs")
	})

}

func TestGet(t *testing.T) {
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

		require.NoError(t, err)

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

		require.NoError(t, err, "Error decoding")

		assert.Equal(t, expectedBook, response, "Response body differs")
	})
	t.Run("Cannot retreive book from context", func(t *testing.T) {
		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("GET", "/book/{id}", nil)
		params := map[string]string{"id": "5"}
		req = mux.SetURLVars(req, params)

		require.NoError(t, err)

		ctx := context.WithValue(req.Context(), "book", nil)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Get)
		handler.ServeHTTP(rr, req)
		expectedResponse := "Internal server error:error retrieving book from context" + "\n"

		assert.Equal(t, expectedResponse, rr.Body.String(), "Response body differs")
	})
}

func TestDelete(t *testing.T) {
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

		require.NoError(t, err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		db.On("DeleteBook", book).Return(nil)

		bookHandler.Db = &db

		handler := http.HandlerFunc(bookHandler.Delete)

		handler.ServeHTTP(rr, req)

		require.Equal(t, http.StatusNoContent, rr.Code)
	})
	t.Run("Error retrieving book from context", func(t *testing.T) {
		bookHandler := handler.BookHandler{}

		req, err := http.NewRequest("DELETE", "/books/{id}", nil)
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)

		require.NoError(t, err)

		ctx := context.WithValue(req.Context(), "book", nil)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Delete)

		handler.ServeHTTP(rr, req)
		expectedResponse := "Internal server error:error retrieving book from context" + "\n"

		require.Equal(t, http.StatusInternalServerError, rr.Code)

		assert.Equal(t, expectedResponse, rr.Body.String(), "Real and expected response differs")
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

		require.NoError(t, err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		db.On("DeleteBook", book).Return(errors.New("Error while deleting book"))
		bookHandler.Db = db
		handler := http.HandlerFunc(bookHandler.Delete)

		handler.ServeHTTP(rr, req)
		expectedResponse := "Error while deleting book" + "\n"
		assert.Equal(t, expectedResponse, rr.Body.String(), "Real and expected response differs")
	})
}
