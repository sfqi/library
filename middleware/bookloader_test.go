package middleware_test

import (
	"errors"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	imock "github.com/sfqi/library/interactor/mock"
	"github.com/sfqi/library/log"
	"github.com/sfqi/library/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	require := require.New(t)
	t.Run("Error converting id", func(t *testing.T) {
		bookLoader := middleware.BookLoader{Logger: log.New()}
		req, err := http.NewRequest("GET", "/{id}", nil)
		require.NoError(err)

		params := map[string]string{"id": "rrr"}
		req = mux.SetURLVars(req, params)
		bookHandler := &bookHandler{}
		newHandler := bookLoader.GetBook(bookHandler)

		rr := httptest.NewRecorder()
		newHandler.ServeHTTP(rr, req)
		expectedError := "invalid book id\n"

		assert.Equal(http.StatusBadRequest, rr.Code)
		assert.Equal(expectedError, rr.Body.String())
	})
	t.Run("Error finding book with given ID", func(t *testing.T) {
		interactor := &imock.Book{}
		bookHandler := &bookHandler{}

		req, err := http.NewRequest("GET", "/{id}", nil)
		require.NoError(err)

		params := map[string]string{"id": "6"}
		req = mux.SetURLVars(req, params)

		interactor.On("FindById", 6).Return(nil, errors.New("Book with given Id can not be found"))

		bookLoader := &middleware.BookLoader{Interactor: interactor, Logger: log.New()}

		newHandler := bookLoader.GetBook(bookHandler)
		rr := httptest.NewRecorder()

		newHandler.ServeHTTP(rr, req)
		expectedRespose := "book not found\n"
		assert.Equal(expectedRespose, rr.Body.String())
	})
	t.Run("Expected response and actual response", func(t *testing.T) {

		interactor := &imock.Book{}
		bookLoader := middleware.BookLoader{Logger: log.New()}

		req, err := http.NewRequest("GET", "/{id}", nil)
		require.NoError(err)

		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		bookHandler := &bookHandler{}
		interactor.On("FindById", 1).Return(&model.Book{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		}, nil, errors.New("Error finding book by ID"))

		bookLoader.Interactor = interactor

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

		assert.Equal(expectedResponse, *book)
	})
}
