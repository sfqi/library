package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gorilla/mux"

	"net/http"
	"strconv"

	"github.com/library/domain/model"
	"github.com/library/handler/dto"
	openlibrarydto "github.com/library/openlibrary/dto"

	"github.com/library/repository/mock"
)

type BookHandler struct {
	Db  *mock.DB
	Olc openLibraryClient
}

type openLibraryClient interface {
	FetchBook(isbn string) (*openlibrarydto.Book, error)
}

func NewBookHandler(db *mock.DB) *BookHandler {
	return &BookHandler{
		Db: db,
	}
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

func (b *BookHandler) Get(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	allBooks := b.Db.GetAllBooks()

	var bookResponses []dto.BookResponse

	for _, book := range allBooks {
		response := dto.BookResponse{
			ID:            book.Id,
			Title:         book.Title,
			Author:        book.Author,
			Isbn:          book.Isbn,
			Isbn13:        book.Isbn13,
			OpenLibraryId: book.OpenLibraryId,
			CoverId:       book.CoverId,
			Year:          book.Year,
		}

		bookResponses = append(bookResponses, response)
	}

	err := json.NewEncoder(w).Encode(bookResponses)
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

	if err := b.Db.Create(book); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	bookResponse := toBookResponse(*book)

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

	bookToAdd := model.Book{
		Id:            b.Db.Id,
		Title:         book.Title,
		Author:        book.Author[0].Name,
		Isbn:          isbn10,
		Isbn13:        isbn13,
		OpenLibraryId: libraryId,
		CoverId:       CoverId,
		Year:          book.Year,
	}
	return &bookToAdd
}

func (b *BookHandler) Update(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	updateBookRequest := dto.UpdateBookRequest{}
	err := json.NewDecoder(r.Body).Decode(&updateBookRequest)
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

	book, err := b.Db.FindBookByID(id)
	if err != nil {
		errorFindingBook(w, err)
		return
	}

	book.Id = id
	book.Title = updateBookRequest.Title
	book.Year = updateBookRequest.Year

	if err := b.Db.Update(book); err != nil {
		errorFindingBook(w, err)
		return
	}
	bookResponse := toBookResponse(*book)

	err = json.NewEncoder(w).Encode(bookResponse)
	if err != nil {
		errorEncoding(w, err)
		return
	}
}

func (b *BookHandler) Index(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		errorConvertingId(w, err)
		return
	}
	book, err := b.Db.FindBookByID(id)
	if err != nil {
		errorFindingBook(w, err)
		return
	}
	bookResponse := toBookResponse(*book)

	err = json.NewEncoder(w).Encode(bookResponse)
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
