package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetBooks(t *testing.T) {

	req, err := http.NewRequest("GET", "/books", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetBooks)

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
