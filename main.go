package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/library/openlibrary"
)

type BookModel struct {
	Id            int    `json:"id,omitempty"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Isbn          string `json:"isbn_10"`
	Isbn13        string `json:"isbn_13"`
	OpenLibraryId string
	CoverId       string `json:"cover"`
	Year          string `json:"publish_date"`
}

var client openlibrary.Client
var books = []BookModel{
	{
		Id:            1,
		Title:         "some title",
		Author:        "some author",
		Isbn:          "some isbn",
		Isbn13:        "some isbon13",
		OpenLibraryId: "again some id",
		CoverId:       "some cover ID",
		Year:          "2019",
	},
	{
		Id:            2,
		Title:         "other title",
		Author:        "other author",
		Isbn:          "other isbn",
		Isbn13:        "other isbon13",
		OpenLibraryId: "other some id",
		CoverId:       "other cover ID",
		Year:          "2019",
	},
	{
		Id:            3,
		Title:         "another title",
		Author:        "another author",
		Isbn:          "another isbn",
		Isbn13:        "another isbon13",
		OpenLibraryId: "another some id",
		CoverId:       "another cover ID",
		Year:          "2019",
	},
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(books)
	if err != nil {
		http.Error(w,"Not found",http.StatusNotFound)
		return
	}
}

type createBookRequest struct {
	ISBN string `json:"ISBN"`
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createBook createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&createBook); err != nil {
		http.Error(w,"Error while decoding from request body",400)
		return
	}
	fmt.Println(createBook.ISBN)
	book, err := client.FetchBook(createBook.ISBN)
	if err != nil {
		http.Error(w,"Error while fetching book: " + err.Error(),400)
		return
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		http.Error(w,"Internal server error:"+err.Error(),500)
		return
	}
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book BookModel
	err := json.NewDecoder(r.Body).Decode(&book)
	if strconv.Itoa(book.Id) != "0" {
		errors.New("Id of the book can not be changed")
		http.Error(w,"Id cannot be changed",400)
		return
	}
	if err != nil {
		http.Error(w,"Error while decoding from request body",400)
		return
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w,"Error while converting url parameter into integer",400)
		return
	}
	if id > len(books){
		// if we have 3 books, and pass 5 as id, its out of range... should be bad request?
		http.Error(w,"Given index is out of bounds",400)
		return
	}

	for i, b := range books {
		if b.Id == id {
			books[i] = book
			if err := json.NewEncoder(w).Encode(book); err != nil {
				http.Error(w,"Bad request: "+err.Error(),400)
				return
			}
			break
		}
	}
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book BookModel

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w,"Error while converting url parameter into integer",400)
		return
	}
	if id >len(books){
		http.Error(w,"Given index is out of bounds",400)
		return
	}

	for i, b := range books {
		if b.Id == id {
			book = books[i]
		}
	}
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		http.Error(w,"Error encoding response into book", 500)
		return
	}

}

func main() {
	// Setting env var
	client = openlibrary.Client{}
	fmt.Println(os.Getenv("LIBRARY"))
	client.Url = os.Getenv("LIBRARY")

	r := mux.NewRouter()

	r.HandleFunc("/books", GetBooks).Methods("GET")
	r.HandleFunc("/books", CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/book/{id}", GetBook).Methods("GET")
	http.ListenAndServe(":8080", r)
}

