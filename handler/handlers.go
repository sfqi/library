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
)

var openLibraryURL = os.Getenv("LIBRARY")
var client = *openlibrary.NewClient(openLibraryURL)

type db struct {
	id    int
	books []model.Book
}

func (bm db) FindBookById(id int) (book model.Book, location int, found bool) {
	for i, b := range bm.books {
		if b.Id == id {
			book = b
			location = i
			found = true
			break
		}
	}

	return book, location, found
}

var books = []model.Book{
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

var shelf = &db{
	id:    len(books),
	books: books,
}

type Book struct{

}

func(h *Book) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(shelf.books)
	if err != nil {
		fmt.Println("error while getting books: ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

type createBookRequest struct {
	ISBN string `json:"ISBN"`
}

func (h *Book)Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createBook createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&createBook); err != nil {
		errorDecodingBook(w, err)
		return
	}
	fmt.Println(createBook.ISBN)
	book, err := client.FetchBook(createBook.ISBN)
	if err != nil {
		fmt.Println("error while fetching book: ", err)
		http.Error(w, "Error while fetching book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bookToAdd := CreateBookModelFromBook(*book)
	shelf.books = append(shelf.books, bookToAdd)

	if err := json.NewEncoder(w).Encode(book); err != nil {
		errorEncoding(w, err)
		return
	}
}

func CreateBookModelFromBook(b dto.Book) (bm model.Book) {

	shelf.id = shelf.id + 1

	isbn10 := ""
	if b.Identifier.ISBN10 != nil {
		isbn10 = b.Identifier.ISBN10[0]
	}
	isbn13 := ""
	if b.Identifier.ISBN13 != nil {
		isbn13 = b.Identifier.ISBN13[0]
	}
	fmt.Println(shelf.id, isbn10, isbn13)
	CoverId := ""
	if b.Cover.Url != "" {
		part1 := strings.Split(b.Cover.Url, "/")[5]
		part2 := strings.Split(part1, ".")[0]
		CoverId = strings.Split(part2, "-")[0]
	}
	libraryId := ""
	if b.Identifier.Openlibrary != nil {
		libraryId = b.Identifier.Openlibrary[0]
	}

	bookToAdd := model.Book{
		Id:            shelf.id,
		Title:         b.Title,
		Author:        b.Author[0].Name,
		Isbn:          isbn10,
		Isbn13:        isbn13,
		OpenLibraryId: libraryId,
		CoverId:       CoverId,
		Year:          b.Year,
	}
	return bookToAdd
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
	bookWithId, location, found := shelf.FindBookById(id)
	fmt.Println(bookWithId.Id, bookWithId.Title, bookWithId.Author)
	if !found {
		errorFindingBook(w, err)
		return
	}

	shelf.books[location] = book
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
	book, _, found := shelf.FindBookById(id)
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
