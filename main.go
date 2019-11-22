package main

import (
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
		User:     "bojan",
		Password: "bojan",
		Name:     "library",
	}

	//db := inmemory.NewDB()

	store,err:= postgres.Open(config)
	if err != nil{
		panic(err)
	}

	bookHandler := &handler.BookHandler{
		DataBase: store,
		Olc: olc,
	}

	bodyDump := middleware.BodyDump{
		Logger: log.New(),
	}
	r := mux.NewRouter()

	r.HandleFunc("/books", bookHandler.Index).Methods("GET")
	r.HandleFunc("/books", bookHandler.Create).Methods("POST")
	r.HandleFunc("/books/{id}", bookHandler.Update).Methods("PUT")
	r.HandleFunc("/book/{id}", bookHandler.Get).Methods("GET")
	r.HandleFunc("/book/{id}", bookHandler.Delete).Methods("DELETE")
	r.Use(bodyDump.Dump)
	http.ListenAndServe(":8080", r)
}
