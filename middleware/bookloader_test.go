package middleware

import (
	"errors"
	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type bookHandler struct {
	bookFromContext *model.Book
}

func (bh *bookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	book := r.Context().Value("book").(*model.Book)
	bh.bookFromContext = book
	w.WriteHeader(http.StatusOK)
}

func TestGetBook(t *testing.T) {
	assert := assert.New(t)
	t.Run("Error converting id", func(t *testing.T) {
		bookLoader := BookLoader{}
		req, err := http.NewRequest("GET", "/{id}", nil)
		if err != nil {
			t.Fatal(err)
		}
		params := map[string]string{"id": "rrr"}
		req = mux.SetURLVars(req, params)
		bookHandler := &bookHandler{}
		newHandler := bookLoader.GetBook(bookHandler)

		rr := httptest.NewRecorder()
		newHandler.ServeHTTP(rr, req)
		expectedError := "Error while converting url parameter into integer\n"

		assert.NotEqual(http.StatusOK, rr.Code)

		assert.Equal(expectedError, rr.Body.String(), "Response body differs")
	})
	t.Run("Error finding book with given ID", func(t *testing.T) {
		store := &mock.Store{}

		req, err := http.NewRequest("GET", "/{id}", nil)
		params := map[string]string{"id": "6"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Fatal(err)
		}
		store.On("FindBookById", 6).Return(nil, errors.New("Book with given Id can not be found"))
		bookHandler := &bookHandler{}
		bookLoader := &BookLoader{Db: store}
		newHandler := bookLoader.GetBook(bookHandler)
		rr := httptest.NewRecorder()

		newHandler.ServeHTTP(rr, req)
		expectedRespose := "Book with given Id can not be found" + "\n"
		assert.Equal(expectedRespose, rr.Body.String(), "Response body differs")
	})
	t.Run("Expected response and actual response", func(t *testing.T) {

		store := mock.Store{}
		bookLoader := BookLoader{}

		req, err := http.NewRequest("GET", "/{id}", nil)
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		bookHandler := &bookHandler{}
		store.On("FindBookById", 1).Return(&model.Book{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		}, nil, errors.New("Error finding book by ID"))
		bookLoader.Db = &store
		newHandler := bookLoader.GetBook(bookHandler)
		newHandler.ServeHTTP(rr, req)
		expectedResponse := model.Book{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		}
		book := bookHandler.bookFromContext

		assert.Equal(expectedResponse, *book, "Response body differs")
	})
}
