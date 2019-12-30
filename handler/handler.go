package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
	"github.com/sfqi/library/interactor"
	"regexp"

	"net/http"
)

var yearRgx = regexp.MustCompile(`[0-9]{4}`)

type BookHandler struct {
	Book *interactor.Book
}

func toBookResponse(b model.Book) *dto.BookResponse {
	bookResponse := dto.BookResponse{}

	bookResponse.ID = b.Id
	bookResponse.Title = b.Title
	bookResponse.Author = b.Author
	bookResponse.Isbn = b.Isbn
	bookResponse.Isbn13 = b.Isbn13
	bookResponse.OpenLibraryId = b.OpenLibraryId
	bookResponse.CoverId = b.CoverId
	bookResponse.Year = b.Year

	return &bookResponse
}

func (b *BookHandler) Index(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	allBooks, err := b.Book.FindAll()
	if err != nil {
		http.Error(w, "Error finding books", http.StatusInternalServerError)
		return
	}
	var bookResponses []*dto.BookResponse

	for _, book := range allBooks {

		bookResponses = append(bookResponses, toBookResponse(*book))
	}

	err = json.NewEncoder(w).Encode(bookResponses)
	if err != nil {
		fmt.Println("error while getting books: ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

func (b *BookHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createBook dto.CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&createBook); err != nil {
		errorDecodingBook(w, err)
		return
	}

	fmt.Println(createBook.ISBN)
	book, err := b.Book.Create(createBook)
	if err != nil {
		http.Error(w, "Internal server error:"+err.Error(), http.StatusInternalServerError)
		return
	}

	bookResponse := *toBookResponse(*book)
	if err := json.NewEncoder(w).Encode(bookResponse); err != nil {
		errorEncoding(w, err)
		return
	}
}

func (b *BookHandler) Update(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	book, ok := r.Context().Value("book").(*model.Book)
	if !ok {
		errorContex(w, errors.New("error retrieving book from context"))
		return
	}
	updateBookRequest := dto.UpdateBookRequest{}
	err := json.NewDecoder(r.Body).Decode(&updateBookRequest)
	if err != nil {
		errorDecodingBook(w, err)
		return
	}
	updatedBook, err := b.Book.Update(book, updateBookRequest)
	if err != nil {
		http.Error(w, "error updating book", http.StatusInternalServerError)
		return
	}
	bookResponse := *toBookResponse(*updatedBook)

	err = json.NewEncoder(w).Encode(bookResponse)
	if err != nil {
		errorEncoding(w, err)
		return
	}
}

func (b *BookHandler) Get(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	book, ok := r.Context().Value("book").(*model.Book)
	if !ok {
		errorContex(w, errors.New("error retrieving book from context"))
		return
	}

	bookResponse := *toBookResponse(*book)

	if err := json.NewEncoder(w).Encode(bookResponse); err != nil {
		errorEncoding(w, err)
		return
	}
}

func (b *BookHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	book, ok := r.Context().Value("book").(*model.Book)
	if !ok {
		errorContex(w, errors.New("error retrieving book from context"))
		return
	}
	fmt.Println("Book from context: ", book)

	if err := b.Book.Delete(book); err != nil {
		http.Error(w, "Error while deleting book", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func errorDecodingBook(w http.ResponseWriter, err error) {
	fmt.Println("error while decoding book from response body: ", err)
	http.Error(w, "Error while decoding from request body", http.StatusBadRequest)
}

func errorEncoding(w http.ResponseWriter, err error) {
	fmt.Println("error while encoding book: ", err)
	http.Error(w, "Internal server error:"+err.Error(), http.StatusInternalServerError)
}

func errorContex(w http.ResponseWriter, err error) {
	fmt.Println("error from context: ", err)
	http.Error(w, "Internal server error:"+err.Error(), http.StatusInternalServerError)
}
