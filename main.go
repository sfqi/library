package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/library/openlibrary"
	"github.com/library/handler"

)
var Client openlibrary.Client
func main() {
	// Setting env var
	openLibraryURL := os.Getenv("LIBRARY")
	Client := *openlibrary.NewClient(openLibraryURL)

	r := mux.NewRouter()

	r.HandleFunc("/books", handler.GetBooks).Methods("GET")
	r.HandleFunc("/books", handler.CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", handler.UpdateBook).Methods("PUT")
	r.HandleFunc("/book/{id}", handler.GetBook).Methods("GET")
	http.ListenAndServe(":8080", r)
}

