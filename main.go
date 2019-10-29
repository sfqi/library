package main

import (
	"github.com/library/openlibrary"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/library/handler"
)

func main() {
	openLibraryUrl := os.Getenv("LIBRARY")
	olc := openlibrary.NewClient(openLibraryUrl)
	bookHandler := handler.Book{
		olc,
	}
	//bookHandler.OpenlibraryClient
	r := mux.NewRouter()

	r.HandleFunc("/books", bookHandler.Get).Methods("GET")
	r.HandleFunc("/books", bookHandler.Create).Methods("POST")
	r.HandleFunc("/books/{id}", bookHandler.Update).Methods("PUT")
	r.HandleFunc("/book/{id}", bookHandler.Index).Methods("GET")
	http.ListenAndServe(":8080", r)
}
