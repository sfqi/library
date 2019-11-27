package main

import (
	"net/http"
	"os"

	"github.com/sfqi/library/log"

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

	bodyDump := middleware.BodyDump{
		Logger: log.New(),
	}

	bookLoad := middleware.BookLoader{
		Db: db,
	}
	r := mux.NewRouter()
	s := r.PathPrefix("/book").Subrouter()

	r.HandleFunc("/books", bookHandler.Index).Methods("GET")
	r.HandleFunc("/books", bookHandler.Create).Methods("POST")
	s.HandleFunc("/{id}", bookHandler.Update).Methods("PUT")
	s.HandleFunc("/{id}", bookHandler.Get).Methods("GET")
	s.HandleFunc("/{id}", bookHandler.Delete).Methods("DELETE")
	r.Use(bodyDump.Dump)
	s.Use(bookLoad.GetBook)

	http.ListenAndServe(":8080", r)
}
