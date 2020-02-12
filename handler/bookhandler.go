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

type BookHandler struct {
	Interactor bookInteractor
}

func NewBookHandler(bookInteractor bookInteractor) *BookHandler {
	return &BookHandler{bookInteractor}
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

func (b *BookHandler) Index(w http.ResponseWriter, r *http.Request) *HTTPError {

	w.Header().Set("Content-Type", "application/json")
	allBooks, err := b.Interactor.FindAll()
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	var bookResponses []*dto.BookResponse

	for _, book := range allBooks {

		bookResponses = append(bookResponses, toBookResponse(*book))
	}

	err = json.NewEncoder(w).Encode(bookResponses)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}

func (b *BookHandler) Create(w http.ResponseWriter, r *http.Request) *HTTPError {
	w.Header().Set("Content-Type", "application/json")

	var createBook dto.CreateBookRequest

	if err := json.NewDecoder(r.Body).Decode(&createBook); err != nil {
		return newHTTPError(http.StatusBadRequest, err)
	}

	fmt.Println(createBook.ISBN)
	book, err := b.Interactor.Create(createBook)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}

	bookResponse := *toBookResponse(*book)
	if err := json.NewEncoder(w).Encode(bookResponse); err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}

func (b *BookHandler) Update(w http.ResponseWriter, r *http.Request) *HTTPError {

	w.Header().Set("Content-Type", "application/json")
	book, ok := r.Context().Value("book").(*model.Book)
	if !ok {
		err := errors.New("error retrieving book from context")
		return newHTTPError(http.StatusInternalServerError, err)
	}
	updateBookRequest := dto.UpdateBookRequest{}
	err := json.NewDecoder(r.Body).Decode(&updateBookRequest)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, err)
	}

	updatedBook, err := b.Interactor.Update(book, updateBookRequest)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	bookResponse := *toBookResponse(*updatedBook)

	err = json.NewEncoder(w).Encode(bookResponse)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}

func (b *BookHandler) Get(w http.ResponseWriter, r *http.Request) *HTTPError {

	w.Header().Set("Content-Type", "application/json")

	book, ok := r.Context().Value("book").(*model.Book)
	if !ok {
		err := errors.New("error retrieving book from context")
		return newHTTPError(http.StatusInternalServerError, err)
	}

	bookResponse := *toBookResponse(*book)

	if err := json.NewEncoder(w).Encode(bookResponse); err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}

func (b *BookHandler) Delete(w http.ResponseWriter, r *http.Request) *HTTPError {
	w.Header().Set("Content-Type", "application/json")
	book, ok := r.Context().Value("book").(*model.Book)
	if !ok {
		err := errors.New("error retrieving book from context")
		return newHTTPError(http.StatusInternalServerError, err)
	}
	fmt.Println("Book from context: ", book)

	if err := b.Interactor.Delete(book); err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	w.WriteHeader(http.StatusNoContent)
	return nil
}
