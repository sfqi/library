package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"

	"net/http"

	openlibrarydto "github.com/sfqi/library/openlibrary/dto"
)

var yearRgx = regexp.MustCompile(`[0-9]{4}`)

type store interface {
	FindBookById(int) (*model.Book, error)
	CreateBook(*model.Book) error
	UpdateBook(*model.Book) error
	FindAllBooks() ([]*model.Book, error)
	DeleteBook(*model.Book)error

}

type BookHandler struct {
	Db  store
	Olc openLibraryClient
}

type openLibraryClient interface {
	FetchBook(isbn string) (*openlibrarydto.Book, error)
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
	allBooks, err := b.Db.FindAllBooks()
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
	openlibraryBook, err := b.Olc.FetchBook(createBook.ISBN)

	if err != nil {
		fmt.Println("error while fetching book: ", err)
		http.Error(w, "Error while fetching book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	book := b.toBook(openlibraryBook)

	if err := b.Db.CreateBook(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	bookResponse := *toBookResponse(*book)
	if err := json.NewEncoder(w).Encode(bookResponse); err != nil {
		errorEncoding(w, err)
		return
	}
}

func (b *BookHandler) toBook(book *openlibrarydto.Book) (bm *model.Book) {
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
	author := ""
	if book.Author != nil {
		author = book.Author[0].Name
	}

	year := 0
	var err error
	yearString := yearRgx.FindString(book.Year)
	if(yearString != ""){
		year, err = strconv.Atoi(yearString)
		if err != nil {
			fmt.Println("error while converting year from string to int", err)
			return nil
		}
	}


	bookToAdd := model.Book{
		Title:         book.Title,
		Author:        author,
		Isbn:          isbn10,
		Isbn13:        isbn13,
		OpenLibraryId: libraryId,
		CoverId:       CoverId,
		Year:          year,
	}
	return &bookToAdd
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

	book.Title = updateBookRequest.Title
	book.Year = updateBookRequest.Year
	if err := b.Db.UpdateBook(book); err != nil {
		http.Error(w, "error updating book", http.StatusInternalServerError)
		return
	}
	bookResponse := *toBookResponse(*book)

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

	if err := b.Db.DeleteBook(book); err != nil {
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
