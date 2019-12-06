package handler_test

import (
	"github.com/sfqi/library/handler"
	"context"
	"bytes"
	"errors"
	"github.com/gorilla/mux"

	"github.com/stretchr/testify/assert"

	"encoding/json"


	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
	openlibrarydto "github.com/sfqi/library/openlibrary/dto"
	olmock "github.com/sfqi/library/openlibrary/mock"
	"github.com/sfqi/library/repository/mock"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func initializeBooks() []*model.Book {
	books := []*model.Book{
		{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
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
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
		},
	}
	return books
}

var bookHandler handler.BookHandler = handler.BookHandler{
	Olc: nil,
}

func TestIndex(t *testing.T) {
	var books = initializeBooks()
	var db = mock.NewStore(books, nil)
	bookHandler.Db = db
	booksExpected := []*dto.BookResponse{
		&dto.BookResponse{
			ID:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		},
		&dto.BookResponse{
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
	handler := http.HandlerFunc(bookHandler.Index)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []dto.BookResponse
	err = json.NewDecoder(rr.Body).Decode(&response)

	for i, _ := range booksExpected {
		if *booksExpected[i] != response[i] {
			t.Errorf("we did not get the same response: got%v want %v",
				response[i], *booksExpected[i])
		}
	}
}

func TestUpdate(t *testing.T) {
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
	t.Run("assertion of expected response, and actual response", func(t *testing.T) {
		var books = initializeBooks()
		var db = mock.NewStore(books, nil)
		bookHandler.Db = db
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"title":"test title", "year":2019}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		ctx := context.WithValue(req.Context(),"book", book) 
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Status code differs. Expected %d. Got %d", http.StatusOK, status)
		}
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

		assert.Equal(t, bookExpected, response, "Response body differs")
	})
	t.Run("Error decoding Book attributes", func(t *testing.T) {
		var db = &mock.Store{}
		bookHandler.Db = db
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"id":"12","title":zdravo}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}
		ctx := context.WithValue(req.Context(),"book", book) 
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())

		}
	})

}

func TestCreate(t *testing.T) {
	t.Run("Invalid request body", func(t *testing.T) {
		var db = &mock.Store{}
		bookHandler.Db = db
		clmock := olmock.Client{}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{ISBN:"0140447938"}`)))

		if err != nil {
			t.Errorf("Error occured while sending request, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Create)

		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})
	t.Run("Fetching book error", func(t *testing.T) {
		var db = &mock.Store{}
		bookHandler.Db = db
		clmock := olmock.Client{nil,
			errors.New("Error while fetching book"),
		}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938222"}`)))

		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Create)

		handler.ServeHTTP(rr, req)
		contains := strings.Contains(rr.Body.String(), clmock.Err.Error())
		if !contains && rr.Code != http.StatusBadRequest {
			t.Errorf("Expected error to be %s, got error: %s", clmock.Err.Error(), rr.Body.String())
		}
	})
	t.Run("Testing book creation", func(t *testing.T) {
		var books = initializeBooks()
		var db = mock.NewStore(books, nil)
		bookHandler.Db = db
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
			Year:  2007,
		},
			nil,
		}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938"}`)))

		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

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
			2019,
		}
		var response dto.BookResponse
		err = json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Errorf("Error decoding %s", err.Error())
		}
		assert.Equal(t, bookExpected, response, "Response body differs")
	})

}

func TestGet(t *testing.T) {
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
	t.Run("Successfully retrieved book", func(t *testing.T) {
		var books = initializeBooks()
		var db = mock.NewStore(books, nil)
		bookHandler.Db = db
		req, err := http.NewRequest("GET", "/book/{id}", nil)
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		ctx := context.WithValue(req.Context(),"book", book) 
		req = req.WithContext(ctx)
		
		rr := httptest.NewRecorder()

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

		assert.Equal(t, expectedBook, response, "Response body differs")
	})
}

func TestDelete(t *testing.T) {
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
	
	t.Run("Book succesfully deleted", func(t *testing.T) {
		var books = initializeBooks()
		var db = mock.NewStore(books, nil)
		bookHandler.Db = db
		req, err := http.NewRequest("DELETE", "/books/{id}", nil)
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		ctx := context.WithValue(req.Context(),"book", book) 
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Delete)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("Status code differs. Expected %d. Got %d", http.StatusNoContent, status)
		}
	})
}
