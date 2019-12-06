package middleware

import (
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type bookHandler struct{
	bookFromContext *model.Book
}

func(bh *bookHandler) ServeHTTP(w http.ResponseWriter,r *http.Request){
	book := r.Context().Value("book").(*model.Book)
	bh.bookFromContext = book
	w.WriteHeader(http.StatusOK)
}


func testHandler(w http.ResponseWriter,r *http.Request){
	book := r.Context().Value("book").(*model.Book)
	fmt.Println("Context: ")
	fmt.Println(book)
	w.WriteHeader(http.StatusOK)
}

func TestGetBook(t *testing.T){
	t.Run("Error converting id",func(t *testing.T){
		bookStore := mock.NewStore(nil,nil)
		bookLoader := BookLoader{Db:bookStore}

		req, err := http.NewRequest("GET", "/{id}", nil)
		if err != nil {
			t.Fatal(err)
		}
		params := map[string]string{"id": "rrr"}
		req = mux.SetURLVars(req, params)
		handler := http.HandlerFunc(testHandler)
		rr := httptest.NewRecorder()
		newHandler := bookLoader.GetBook(handler)
		newHandler.ServeHTTP(rr, req)
		expectedError := "Error whilr converting url parameter into integer"
		if rr.Code != http.StatusBadRequest && expectedError != rr.Body.String() {
			t.Errorf("Expected code to be %d, got %d", http.StatusOK, rr.Code)
		}
	})
	t.Run("Error finding book with given ID",func(t *testing.T){
		books := []*model.Book{
			{
				Id:            1,
				Title:         "some title",
				Author:        "some author",
				Isbn:          "some isbn",
				Isbn13:        "some isbon13",
				OpenLibraryId: "again some id",
				CoverId:       "some cover ID",
				Year:          2019,
			},
			{
				Id:            2,
				Title:         "other title",
				Author:        "other author",
				Isbn:          "other isbn",
				Isbn13:        "other isbon13",
				OpenLibraryId: "other some id",
				CoverId:       "other cover ID",
				Year:          2019,
			},
		}
		bookStore := mock.NewStore(books,errors.New("Book with given Id can not be found"+"\n"))
		bookLoader := BookLoader{Db:bookStore}

		req, err := http.NewRequest("GET", "/{id}", nil)
		params := map[string]string{"id": "6"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Fatal(err)
		}
		handler := http.HandlerFunc(testHandler)
		rr := httptest.NewRecorder()
		newHandler := bookLoader.GetBook(handler)
		newHandler.ServeHTTP(rr, req)
		assert.Equal(t,bookStore.Err.Error(),rr.Body.String(),"Response body differs")
	})
	t.Run("Expected response and actual response",func(t *testing.T){
		books := []*model.Book{
			&model.Book{
				Id:            1,
				Title:         "some title",
				Author:        "some author",
				Isbn:          "some isbn",
				Isbn13:        "some isbon13",
				OpenLibraryId: "again some id",
				CoverId:       "some cover ID",
				Year:          2019,
			},
			&model.Book{
				Id:            2,
				Title:         "other title",
				Author:        "other author",
				Isbn:          "other isbn",
				Isbn13:        "other isbon13",
				OpenLibraryId: "other some id",
				CoverId:       "other cover ID",
				Year:          2019,
			},
		}
		bookStore := mock.NewStore(books,nil)
		bookLoader := BookLoader{Db:bookStore}

		req, err := http.NewRequest("GET", "/{id}", nil)
		params := map[string]string{"id": "1"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		bookHandler := &bookHandler{}
		newHandler := bookLoader.GetBook(bookHandler)
		newHandler.ServeHTTP(rr, req)
		expectedResponse:=model.Book{
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
		fmt.Println(book)
		assert.Equal(t,expectedResponse,*book,"Response body differs")
	})


}
