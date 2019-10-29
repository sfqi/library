package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"

	"net/http"
	"strconv"
	"github.com/library/domain/model"
	"github.com/library/handler/dto"
	"github.com/library/repository/mock"

)


//var client = *openlibrary.NewClient(OpenLibraryURL)
type openLibraryClient interface {
	FetchBook(isbn string)(*dto.Book, error)
}
type Book struct{
	 Olc openLibraryClient
}

func(h *Book) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(mock.Shelf.Books)
	if err != nil {
		fmt.Println("error while getting books: ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (h *Book)Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createBook  dto.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&createBook); err != nil {
		errorDecodingBook(w, err)
		return
	}
	fmt.Println(createBook.ISBN)
	book, err := h.Olc.FetchBook(createBook.ISBN)
	if err != nil {
		fmt.Println("error while fetching book: ", err)
		http.Error(w, "Error while fetching book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bookToAdd := mock.CreateBookModelFromBook(*book)
	mock.Shelf.Books = append(mock.Shelf.Books, bookToAdd)

	if err := json.NewEncoder(w).Encode(book); err != nil {
		errorEncoding(w, err)
		return
	}
}

func (h *Book)Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book model.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		errorDecodingBook(w, err)
		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		errorConvertingId(w, err)
		return
	}
	bookWithId, location, found := mock.Shelf.FindBookById(id)
	fmt.Println(bookWithId.Id, bookWithId.Title, bookWithId.Author)
	if !found {
		errorFindingBook(w, err)
		return
	}

	mock.Shelf.Books[location] = book
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		errorEncoding(w, err)
		return
	}
}

func (h *Book)Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		errorConvertingId(w, err)
		return
	}
	book, _, found := mock.Shelf.FindBookById(id)
	if !found {
		errorFindingBook(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		errorEncoding(w, err)
		return
	}
}

func errorDecodingBook(w http.ResponseWriter, err error) {
	fmt.Println("error while decoding book from response body: ", err)
	http.Error(w, "Error while decoding from request body", http.StatusBadRequest)
}

func errorEncoding(w http.ResponseWriter, err error) {
	fmt.Println("error while encoding book: ", err)
	http.Error(w, "Internal server error:"+err.Error(), http.StatusInternalServerError)
}

func errorConvertingId(w http.ResponseWriter, err error) {
	fmt.Println("Error while converting Id to integer ", err)
	http.Error(w, "Error while converting url parameter into integer", http.StatusBadRequest)
}

func errorFindingBook(w http.ResponseWriter, err error) {
	fmt.Println("Cannot find book with given Id ")
	http.Error(w, "Book with given Id can not be found", http.StatusBadRequest)
}
