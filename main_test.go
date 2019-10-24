package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)


func TestUpdateBook(t *testing.T) {
	t.Run("assertion of expected response, and actual response", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"Title":"another title","Author":"another author","Isbn":"another isbn","Isbn13":"another isbon13","OpenLibraryId": "another some id","CoverId":"another cover ID","Year":"2019"}`)))
		params := map[string]string{"id":"2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdateBook)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("Status code differs. Expected %d. Got %d", http.StatusOK, status)
		}
		expected := `{"id":0,"title":"another title","author":"another author","isbn_10":"","isbn_13":"","OpenLibraryId":"another some id","cover":"","publish_date":""}` + "\n"
		assert.Equal(t, expected, rr.Body.String(), "Response body differs")
	})
	t.Run("Error decoding Book attributes", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"id":"12","title":zdravo}`)))
		params := map[string]string{"id":"2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdateBook)
		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})

	t.Run("Converting id parameter into integer", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"id":12,"title":"zdravo"}`)))
		params := map[string]string{"id":"2ss"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdateBook)
		handler.ServeHTTP(rr, req)
		expectedError := "Error while converting url parameter into integer"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})
	t.Run("Book with given Id can not be found", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"id":12,"title":"zdravo"}`)))
		params := map[string]string{"id":"12"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(UpdateBook)
		handler.ServeHTTP(rr, req)
		expectedError := "Book with given Id can not be found"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})

}


func TestCreateBook(t *testing.T){
	t.Run("Error decoding Book attributes", func(t *testing.T) {
		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":0140447938}`)))

		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(CreateBook)
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

		handler := http.HandlerFunc(CreateBook)
		handler.ServeHTTP(rr, req)
		expectedError := "Error while fetching book"
		contains := strings.Contains(rr.Body.String(), expectedError)
		if status := rr.Code; status != http.StatusInternalServerError && !contains {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})

}

func TestGetBook(t *testing.T) {
	t.Run("Given Id can not be converted", func(t *testing.T){
		req, err := http.NewRequest("GET", "/book/{id}", nil)
		params := map[string]string{"id":"ee"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetBook)
		handler.ServeHTTP(rr, req)

		expectedError := "Error while converting url parameter into integer"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})
	t.Run("Book with given Id can not be found", func(t *testing.T){
		req, err := http.NewRequest("GET", "/book/{id}", nil)
		params := map[string]string{"id":"44"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(GetBook)
		handler.ServeHTTP(rr, req)

		expectedError := "Book with given Id can not be found"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})
}

