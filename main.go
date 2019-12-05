package main

import (
	"fmt"
	"github.com/sfqi/library/repository/postgres"
	"net/http"
	"os"

	"github.com/sfqi/library/log"

	"github.com/sfqi/library/openlibrary"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/handler"
	"github.com/sfqi/library/middleware"
)

func main() {
	openLibraryUrl := os.Getenv("LIBRARY")
	olc := openlibrary.NewClient(openLibraryUrl)
	config := postgres.PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "library",
	}

	store, err:= postgres.Open(config)
	if err != nil{
		panic(err)
	}
	fmt.Println("Successfully connected")
	bookHandler := &handler.BookHandler{
		Db: store,
		Olc: olc,
	}

	bodyDump := middleware.BodyDump{
		Logger: log.New(),
	}

	bookLoad := middleware.BookLoader{
		Db: store,
	}
	r := mux.NewRouter()
	s := r.PathPrefix("/books").Subrouter()

	r.HandleFunc("/books", bookHandler.Index).Methods("GET")
	r.HandleFunc("/books", bookHandler.Create).Methods("POST")
	s.HandleFunc("/{id}", bookHandler.Update).Methods("PUT")
	s.HandleFunc("/{id}", bookHandler.Get).Methods("GET")
	s.HandleFunc("/{id}", bookHandler.Delete).Methods("DELETE")
	r.Use(bodyDump.Dump)
	s.Use(bookLoad.GetBook)

	http.ListenAndServe(":8080", r)
}
