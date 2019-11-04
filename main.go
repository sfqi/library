package main

import (
	"net/http"
	"os"

	"github.com/library/openlibrary"

	"github.com/gorilla/mux"
	"github.com/library/handler"
	"github.com/library/repository/mock"
)

func main() {
	openLibraryUrl := os.Getenv("LIBRARY")
	olc := openlibrary.NewClient(openLibraryUrl)

	r := mux.NewRouter()

	db := mock.NewDB()
	bookHandler := handler.BookHandler{
		Db:  db,
		Olc: olc,
	}


	r.HandleFunc("/books", bookHandler.Get).Methods("GET")
	r.HandleFunc("/books", bookHandler.Create).Methods("POST")
	r.HandleFunc("/books/{id}", bookHandler.Update).Methods("PUT")
	r.HandleFunc("/book/{id}", bookHandler.Index).Methods("GET")

	http.ListenAndServe(":8080", r)
}
