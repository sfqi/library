package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/library/handler"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/books", handler.GetBooks).Methods("GET")
	r.HandleFunc("/books", handler.CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", handler.UpdateBook).Methods("PUT")
	r.HandleFunc("/book/{id}", handler.GetBook).Methods("GET")
	http.ListenAndServe(":8080", r)
}
