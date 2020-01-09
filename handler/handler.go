package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"

	"net/http"
)

var yearRgx = regexp.MustCompile(`[0-9]{4}`)

type HTTPError struct {
	code     int
	internal error
	context  string
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s. Error: %s", h.code, h.context, h.internal)
}

func newHTTPError(code int, err error) *HTTPError {
	return &HTTPError{
		code:     code,
		internal: err,
	}
}

func (e *HTTPError) Wrap(ctx string) *HTTPError {
	e.context = ctx
	return e
}

type BookHandler struct {
	Interactor bookInteractor
}

type bookInteractor interface {
	FindAll() ([]*model.Book, error)
	Create(bookRequest dto.CreateBookRequest) (*model.Book, error)
	Update(book *model.Book, updateBookRequest dto.UpdateBookRequest) (*model.Book, error)
	FindById(id int) (*model.Book, error)
	Delete(book *model.Book) error
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

func (b *BookHandler) Index(w http.ResponseWriter, r *http.Request) error {

	w.Header().Set("Content-Type", "application/json")
	allBooks, err := b.Interactor.FindAll()
	if err != nil {
		http.Error(w, "Error finding books", http.StatusInternalServerError)
		return err
	}
	var bookResponses []*dto.BookResponse

	for _, book := range allBooks {

		bookResponses = append(bookResponses, toBookResponse(*book))
	}

	err = json.NewEncoder(w).Encode(bookResponses)
	if err != nil {
		fmt.Println("error while getting books: ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return err
	}
	return nil
}

func (b *BookHandler) Create(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var createBook dto.CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&createBook); err != nil {
		errorDecodingBook(w, err)
		return err
	}

	fmt.Println(createBook.ISBN)
	book, err := b.Interactor.Create(createBook)

	if err != nil {
		fmt.Println("error while fetching book: ", err)
		http.Error(w, "Error while fetching book: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	bookResponse := *toBookResponse(*book)
	if err := json.NewEncoder(w).Encode(bookResponse); err != nil {
		errorEncoding(w, err)
		return err
	}
	return nil
}

func (b *BookHandler) Update(w http.ResponseWriter, r *http.Request) error {

	w.Header().Set("Content-Type", "application/json")
	book, ok := r.Context().Value("book").(*model.Book)
	if !ok {
		errorContex(w, errors.New("error retrieving book from context"))
		return errors.New("Error retrieving book from context")
	}
	updateBookRequest := dto.UpdateBookRequest{}
	err := json.NewDecoder(r.Body).Decode(&updateBookRequest)
	if err != nil {
		errorDecodingBook(w, err)
		return err
	}

	updatedBook, err := b.Interactor.Update(book, updateBookRequest)
	if err != nil {
		http.Error(w, "error updating book", http.StatusInternalServerError)
		return err
	}
	bookResponse := *toBookResponse(*updatedBook)

	err = json.NewEncoder(w).Encode(bookResponse)
	if err != nil {
		errorEncoding(w, err)
		return err
	}
	return nil
}

func (b *BookHandler) Get(w http.ResponseWriter, r *http.Request) error {

	w.Header().Set("Content-Type", "application/json")

	book, ok := r.Context().Value("book").(*model.Book)
	if !ok {
		errorContex(w, errors.New("error retrieving book from context"))
		return errors.New("error retrieving book from context")
	}

	bookResponse := *toBookResponse(*book)

	if err := json.NewEncoder(w).Encode(bookResponse); err != nil {
		errorEncoding(w, err)
		return err
	}
	return nil
}

func (b *BookHandler) Delete(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	book, ok := r.Context().Value("book").(*model.Book)
	if !ok {
		errorContex(w, errors.New("error retrieving book from context"))
		return errors.New("error retrieving book from context")
	}
	fmt.Println("Book from context: ", book)

	if err := b.Interactor.Delete(book); err != nil {
		http.Error(w, "Error while deleting book", http.StatusInternalServerError)
		return err
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
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
