package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/library/openlibrary"
)

type BookModel struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Isbn          string `json:"isbn_10"`
	Isbn13        string `json:"isbn_13"`
	OpenLibraryId string `json:"olid"`
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
	if err := json.NewEncoder(w).Encode(books); err != nil {
		fmt.Println("error while Marshaling from GetBooks: ", err)
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
		fmt.Println("error while decoding in CreateBook: ", err)
		return
	}
	fmt.Println(createBook.ISBN)
	book, err := client.FetchBook(createBook.ISBN)
	if err != nil {
		fmt.Println("error while fething book: ", err)
		return
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		fmt.Println("error while encode from CreateBook: ", err)
		return
	}
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book BookModel
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		fmt.Println("error while decoding from UpdateBook", err)
		return
	}
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println("error converting string to int from UpdateBook", err)
		return
	}

	for i, b := range books {
		if b.Id == id {
			books[i] = book
			if err := json.NewEncoder(w).Encode(book); err != nil {
				fmt.Println("error while Marshaling from UpdateBook: ", err)
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
		fmt.Println("error converting string to int from GetOneBook", err)
		return
	}

	for i, b := range books {
		if b.Id == id {
			book = books[i]
		}
	}
	json.NewEncoder(w).Encode(book)

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
