package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/mux"

	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/library/domain/model"
	"github.com/library/handler/dto"
	"github.com/library/openlibrary"
	"github.com/library/repository/mock"
)

var openLibraryURL = os.Getenv("LIBRARY")
var client = *openlibrary.NewClient(openLibraryURL)

type BookHandler struct {
	db *mock.DB
}

func NewBookHandler(db *mock.DB) BookHandler {
	return BookHandler{
		db: db,
	}
}

func (b BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allBooks := b.db.GetAllBooks()

	err := json.NewEncoder(w).Encode(allBooks)
	if err != nil {
		fmt.Println("error while getting books: ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

type createBookRequest struct {
	ISBN string `json:"ISBN"`
}

func (b BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createBook createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&createBook); err != nil {
		errorDecodingBook(w, err)
		return
	}
	book, err := client.FetchBook(createBook.ISBN)
	if err != nil {
		fmt.Println("error while fetching book: ", err)
		http.Error(w, "Error while fetching book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bookToAdd := b.createBookModelFromBook(*book)

	if err := b.db.Create(bookToAdd); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		errorEncoding(w, err)
		return
	}
}

func (b BookHandler) createBookModelFromBook(book dto.Book) (bm model.Book) {

	b.db.Id = b.db.Id + 1

	isbn10 := ""
	if book.Identifier.ISBN10 != nil {
		isbn10 = book.Identifier.ISBN10[0]
	}
	isbn13 := ""
	if book.Identifier.ISBN13 != nil {
		isbn13 = book.Identifier.ISBN13[0]
	}

	CoverId := ""
	if book.Cover.Url != "" {
		part1 := strings.Split(book.Cover.Url, "/")[5]
		part2 := strings.Split(part1, ".")[0]
		CoverId = strings.Split(part2, "-")[0]
	}
	libraryId := ""
	if book.Identifier.Openlibrary != nil {
		libraryId = book.Identifier.Openlibrary[0]
	}

	bookToAdd := model.Book{
		Id:            b.db.Id,
		Title:         book.Title,
		Author:        book.Author[0].Name,
		Isbn:          isbn10,
		Isbn13:        isbn13,
		OpenLibraryId: libraryId,
		CoverId:       CoverId,
		Year:          book.Year,
	}
	return bookToAdd
}

func (b BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
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
	bookWithId, location := b.db.FindBookById(id)
	fmt.Println(bookWithId.Id, bookWithId.Title, bookWithId.Author)
	if location < 0 {
		errorFindingBook(w, err)
		return
	}

	b.db.Update(location, book)
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		errorEncoding(w, err)
		return
	}
}

func (b BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		errorConvertingId(w, err)
		return
	}
	book, location := b.db.FindBookById(id)
	if location < 0 {
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
