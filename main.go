package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/library/handler"
)

func main() {
	bookHandler := handler.Book{}
	r := mux.NewRouter()

	r.HandleFunc("/books", bookHandler.Get).Methods("GET")
	r.HandleFunc("/books", bookHandler.Create).Methods("POST")
	r.HandleFunc("/books/{id}", bookHandler.Update).Methods("PUT")
	r.HandleFunc("/book/{id}", bookHandler.Index).Methods("GET")
	http.ListenAndServe(":8080", r)
}
