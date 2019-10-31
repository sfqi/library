package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/library/repository/mock"
	"github.com/stretchr/testify/assert"
)

var bookHandler BookHandler

func init() {
	db := mock.NewDB()

	bookHandler = NewBookHandler(db)
}

func TestGetBooks(t *testing.T) {

	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(bookHandler.GetBooks)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := `[{"id":1,"title":"some title","author":"some author","isbn_10":"some isbn","isbn_13":"some isbon13","olid":"again some id","cover":"some cover ID","publish_date":"2019"},{"id":2,"title":"other title","author":"other author","isbn_10":"other isbn","isbn_13":"other isbon13","olid":"other some id","cover":"other cover ID","publish_date":"2019"},{"id":3,"title":"another title","author":"another author","isbn_10":"another isbn","isbn_13":"another isbon13","olid":"another some id","cover":"another cover ID","publish_date":"2019"}]` + "\n"
	fmt.Println(rr.Body.String())
	fmt.Println(expected)
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got%v want %v",
			rr.Body.String(), expected)
	}
}

func TestUpdateBook(t *testing.T) {
	t.Run("assertion of expected response, and actual response", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"Title":"another title","Author":"another author","Isbn":"another isbn","Isbn13":"another isbon13","olid": "another some id","CoverId":"another cover ID","Year":"2019"}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bookHandler.UpdateBook)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Status code differs. Expected %d. Got %d", http.StatusOK, status)
		}
		expected := `{"id":0,"title":"another title","author":"another author","isbn_10":"","isbn_13":"","olid":"another some id","cover":"","publish_date":""}` + "\n"
		assert.Equal(t, expected, rr.Body.String(), "Response body differs")
	})
	t.Run("Error decoding Book attributes", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"id":"12","title":zdravo}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bookHandler.UpdateBook)
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
		handler := http.HandlerFunc(bookHandler.UpdateBook)
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
		handler := http.HandlerFunc(bookHandler.UpdateBook)
		handler.ServeHTTP(rr, req)
		expectedError := "Book with given Id can not be found"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})

}

func TestCreateBook(t *testing.T) {
	t.Run("Error decoding Book attributes", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":0140447938}`)))

		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bookHandler.CreateBook)
		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})
	t.Run("Fetching book error", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938222"}`))) //kada posaljemo nepostojeci ISBN recimo

		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.CreateBook)
		handler.ServeHTTP(rr, req)
		expectedError := "Error while fetching book"
		contains := strings.Contains(rr.Body.String(), expectedError)
		if status := rr.Code; status != http.StatusInternalServerError && !contains {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})

}

func TestGetBook(t *testing.T) {
	t.Run("Given Id can not be converted", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/book/{id}", nil)
		params := map[string]string{"id": "ee"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bookHandler.GetBook)
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
		handler := http.HandlerFunc(bookHandler.GetBook)
		handler.ServeHTTP(rr, req)

		expectedError := "Book with given Id can not be found"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})
}
