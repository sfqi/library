package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/library/handler"
)

func main() {
	bookHandler := handler.Book{}
	r := mux.NewRouter()

	r.HandleFunc("/books", bookHandler.GetBooks).Methods("GET")
	r.HandleFunc("/books", bookHandler.CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", bookHandler.UpdateBook).Methods("PUT")
	r.HandleFunc("/book/{id}", bookHandler.GetBook).Methods("GET")
	http.ListenAndServe(":8080", r)
}
