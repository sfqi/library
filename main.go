package main

import (
	"net/http"
	"os"

	"github.com/sfqi/library/openlibrary"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/handler"
	middleware "github.com/sfqi/library/middleware"
	"github.com/sfqi/library/repository/inmemory"
)

func main() {
	openLibraryUrl := os.Getenv("LIBRARY")
	olc := openlibrary.NewClient(openLibraryUrl)

	db := inmemory.NewDB()
	bookHandler := handler.BookHandler{
		Db:  db,
		Olc: olc,
	}

	r := mux.NewRouter()

	r.HandleFunc("/books", bookHandler.Index).Methods("GET")
	r.HandleFunc("/books", bookHandler.Create).Methods("POST")
	r.HandleFunc("/books/{id}", bookHandler.Update).Methods("PUT")
	r.HandleFunc("/book/{id}", bookHandler.Get).Methods("GET")
	r.HandleFunc("/book/{id}", bookHandler.Delete).Methods("DELETE")
	r.Use(middleware.LoggingMiddleware)
	http.ListenAndServe(":8080", r)
}
