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
		fmt.Println("error while getting books: ", err)
		http.Error(w,"Bad request",http.StatusBadRequest)
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
		errorDecodingBook(w,err)
		return
	}
	fmt.Println(createBook.ISBN)
	book, err := client.FetchBook(createBook.ISBN)
	if err != nil {
		fmt.Println("error while fetching book: ", err)
		http.Error(w,"Error while fetching book: " + err.Error(),http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		errorEncoding(w,err)
		return
	}
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book BookModel
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		errorDecodingBook(w,err)
		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		errorConvertingId(w,err)
		return
	}
	found := false
	for i, b := range books {
		if b.Id == id {
			books[i] = book
			if err := json.NewEncoder(w).Encode(book); err != nil {
				errorEncoding(w,err)
				return
			}
			found = true
			break
		}
	}
	if !found {
		errorFindingBook(w,err)
		return
	}
}

func FindBookById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		errorConvertingId(w, err)
		return
	}

	found := false
	for i, b := range books {
		if b.Id == id {
			err := json.NewEncoder(w).Encode(books[i])
			if err != nil {
				errorEncoding(w, err)
			}
			found = true
			break
		}
	}
	if !found {
		errorFindingBook(w, err)
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
	r.HandleFunc("/book/{id}", FindBookById).Methods("GET")
	http.ListenAndServe(":8080", r)
}

// Handling errors ***************
func errorDecodingBook(w http.ResponseWriter,err error) {
	fmt.Println("error while decoding book from response body: ", err)
	http.Error(w, "Error while decoding from request body", http.StatusBadRequest)
}

func errorEncoding(w http.ResponseWriter,err error){
	fmt.Println("error while encoding book: ", err)
	http.Error(w,"Internal server error:"+err.Error(),http.StatusInternalServerError)
}

func errorConvertingId(w http.ResponseWriter,err error){
	fmt.Println("Error while converting Id to integer ",err)
	http.Error(w,"Error while converting url parameter into integer",http.StatusBadRequest)
}

func errorFindingBook(w http.ResponseWriter,err error) {
	fmt.Println("Cannot find book with given Id ")
	http.Error(w, "Book with given Id can not be found", http.StatusBadRequest)
}
