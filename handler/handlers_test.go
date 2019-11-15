package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/sfqi/library/handler/dto"
	olmock "github.com/sfqi/library/openlibrary/mock"
	"github.com/sfqi/library/repository/inmemory"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var bookHandler BookHandler = BookHandler{
	Db:  inmemory.NewDB(),
	Olc: nil,
}

func TestIndex(t *testing.T) {
	booksExpected := []*dto.BookResponse{
		&dto.BookResponse{
			ID:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          "2019",
		},
		&dto.BookResponse{
			ID:            2,
			Title:         "other title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          "2019",
		},
		&dto.BookResponse{
			ID:            3,
			Title:         "another title",
			Author:        "another author",
			Isbn:          "another isbn",
			Isbn13:        "another isbon13",
			OpenLibraryId: "another some id",
			CoverId:       "another cover ID",
			Year:          "2019",
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
	t.Run("assertion of expected response, and actual response", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"title":"test title", "year":"2019"}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

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
			Year:          "2019",
		}

		var response dto.BookResponse
		err = json.NewDecoder(rr.Body).Decode(&response)

		assert.Equal(t, bookExpected, response, "Response body differs")
	})
	t.Run("Error decoding Book attributes", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"id":"12","title":zdravo}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})

	t.Run("Converting id parameter into integer", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"id":12,"title":"zdravo"}`)))
		params := map[string]string{"id": "2ss"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)
		expectedError := "Error while converting url parameter into integer"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})
	t.Run("Book with given Id can not be found", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"id":12,"title":"zdravo"}`)))
		params := map[string]string{"id": "12"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)
		expectedError := "Book with given Id can not be found"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})

}

func TestCreate(t *testing.T) {
	t.Run("Error decoding Book attributes", func(t *testing.T) {

		clmock := olmock.Client{nil,
			errors.New("Error while decoding from request body"),
		}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":0140447938}`)))

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

}

func TestGet(t *testing.T) {
	t.Run("Given Id can not be converted", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/book/{id}", nil)
		params := map[string]string{"id": "ee"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Get)
		handler.ServeHTTP(rr, req)

		expectedError := "Error while converting url parameter into integer"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})
	t.Run("Book with given Id can not be found", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/book/{id}", nil)
		params := map[string]string{"id": "44"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Get)

		handler.ServeHTTP(rr, req)

		expectedError := "Book with given Id can not be found"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})
	t.Run("Successfully retrieved book", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/book/{id}", nil)
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Get)
		handler.ServeHTTP(rr, req)
		expectedBook := dto.BookResponse{
			ID:            2,
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          "2019",
		}
		var response dto.BookResponse
		err = json.NewDecoder(rr.Body).Decode(&response)

		assert.Equal(t, expectedBook, response, "Response body differs")
	})
}

func TestDelete(t *testing.T) {
	t.Run("Error converting Id to integer", func(t *testing.T) {
		req, err := http.NewRequest("DELETE", "/books/{id}", nil)
		params := map[string]string{"id": "e"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Delete)

		handler.ServeHTTP(rr, req)

		expectedError :="Error while converting url parameter into integer"+"\n"
		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Status code differs. Expected %d. Got %d", http.StatusBadRequest, status)
		}
		assert.Equal(t, expectedError, rr.Body.String(), "Response body differs")
	})
	t.Run("Error finding book with given Id",func(t *testing.T){
		req, err := http.NewRequest("DELETE", "/books/{id}", nil)
		params := map[string]string{"id": "7"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Delete)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("Status code differs. Expected %d. Got %d", http.StatusBadRequest, status)
		}
		expectedError :="Book with given Id can not be found"+"\n"
		assert.Equal(t, expectedError, rr.Body.String(), "Response body differs")
	})
	t.Run("Book succesfully deleted",func(t *testing.T){
		req, err := http.NewRequest("DELETE", "/books/{id}", nil)
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Delete)

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("Status code differs. Expected %d. Got %d", http.StatusNoContent, status)
		}
	})
}