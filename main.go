package main

import (
	"fmt"
	"github.com/library/repository/postgres"
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
	config,err := postgres.LoadConfig("config.yml")
	//for now we dont need DbStore, just checking if we configured all fine
	_,err = postgres.Open(*config)
	if err != nil{
		panic(err)
	}
	fmt.Println("Successfully connected to dataabse")


	db := mock.NewDB()
	bookHandler := handler.BookHandler{
		Db:  db,
		Olc: olc,
	}

	r := mux.NewRouter()
	r.HandleFunc("/books", bookHandler.Get).Methods("GET")
	r.HandleFunc("/books", bookHandler.Create).Methods("POST")
	r.HandleFunc("/books/{id}", bookHandler.Update).Methods("PUT")
	r.HandleFunc("/book/{id}", bookHandler.Index).Methods("GET")

	http.ListenAndServe(":8080", r)
}
