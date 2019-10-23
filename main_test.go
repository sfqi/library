package main

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
func TestGetBooks(t *testing.T) {

	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetBooks)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
*/
func TestUpdateBook(t *testing.T) {
	req, err := http.NewRequest("PUT", "/books/2", bytes.NewBuffer([]byte(`{"Id":3,"Title":"another title","Author":"another author","Isbn":"another isbn","Isbn13":"another isbon13","OpenLibraryId": "another some id","CoverId":"another cover ID","Year":"2019"}`)))

	if err != nil {
		t.Errorf("Error occured, %s",err)
	}

	rr := httptest.NewRecorder()
	r:=mux.NewRouter()

	r.HandleFunc("/books/{id}",UpdateBook).Methods("PUT")
	r.ServeHTTP(rr,req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusOK, status)
	}
	expected:= `{"id":3,"title":"another title","author":"another author","isbn_10":"","isbn_13":"","OpenLibraryId":"another some id","cover":"","publish_date":""}`+"\n"
	assert.Equal(t, expected, rr.Body.String(), "Response body differs")

}