package handler_test

import (
	"bytes"
	"context"
	"errors"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/handler"

	"github.com/stretchr/testify/mock"

	"github.com/sfqi/library/domain/model"
	olmock "github.com/sfqi/library/openlibrary/mock"

	rmock "github.com/sfqi/library/repository/mock"

	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type StoreMock struct {
	mock.Mock
	Books []*model.Book
	Err   error
}

func (sm *StoreMock) FindBookById(id int) (*model.Book, error) {
	args := sm.Called(id)

	err := args.Error(1)
	if err != nil {
		return nil, err
	}
	book := args.Get(0).(model.Book)
	return &book, args.Error(1)
}

func (sm *StoreMock) FindAllBooks() ([]*model.Book, error) {
	args := sm.Called()

	err := args.Error(1)
	if err != nil {
		return nil, err
	}
	book := args.Get(0).([]*model.Book)
	return book, args.Error(1)
}

func (sm *StoreMock) CreateBook(book *model.Book) error {
	args := sm.Called(book)

	err := args.Error(1)
	if err != nil {
		return err
	}
	return err
}

func (sm *StoreMock) UpdateBook(book *model.Book) error {
	args := sm.Called(book)

	err := args.Error(1)
	if err != nil {
		return err
	}
	return err
}

func (sm *StoreMock) DeleteBook(book *model.Book) error {
	args := sm.Called(book)

	err := args.Error(1)
	if err != nil {
		return err
	}
	return err
}

func initializeBooks() []*model.Book {
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
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
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
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
		},
	}
	return books
}

var bookHandler = handler.BookHandler{
	Olc: nil,
}

func TestIndex(t *testing.T) {
	mockClient := StoreMock{}
	mockClient.On("FindAllBooks").Return([]*model.Book{
		{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
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
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
		},
	}, nil)
	mockClient.FindAllBooks()
	mockClient.AssertExpectations(t)

}

func TestUpdate(t *testing.T) {
	book := &model.Book{
		Id:            2,
		Title:         "test title",
		Author:        "other author",
		Isbn:          "other isbn",
		Isbn13:        "other isbon13",
		OpenLibraryId: "other some id",
		CoverId:       "other cover ID",
		Year:          2019,
	}
	t.Run("assertion of expected response, and actual response", func(t *testing.T) {
		mockClient := StoreMock{}
		book := model.Book{
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2019,
		}

		mockClient.On("UpdateBook", &book).Return(model.Book{
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2019,
		}, nil)
		mockClient.UpdateBook(&book)
		mockClient.AssertExpectations(t)

	})
	t.Run("Error decoding Book attributes", func(t *testing.T) {
		var db = &rmock.Store{}
		bookHandler.Db = db
		req, err := http.NewRequest("PUT", "/books/{id}", bytes.NewBuffer([]byte(`{"id":"12","title":zdravo}`)))
		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)
		if err != nil {
			t.Errorf("Error occured, %s", err)
		}
		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(bookHandler.Update)

		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())

		}
	})

}

func TestCreate(t *testing.T) {
	t.Run("Invalid request body", func(t *testing.T) {
		var db = &rmock.Store{}
		bookHandler.Db = db
		clmock := olmock.Client{}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{ISBN:"0140447938"}`)))

		if err != nil {
			t.Errorf("Error occured while sending request, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Create)

		handler.ServeHTTP(rr, req)
		expectedError := "Error while decoding from request body"
		if status := rr.Code; status != http.StatusBadRequest && rr.Body.String() != expectedError {
			t.Errorf("Expected status code: %d and error: %s,  got: %d and %s", http.StatusBadRequest, expectedError, status, rr.Body.String())
		}
	})
	t.Run("Fetching book error", func(t *testing.T) {
		var db = &rmock.Store{}
		bookHandler.Db = db
		clmock := olmock.Client{nil,
			errors.New("Error while fetching book"),
		}
		bookHandler.Olc = &clmock

		req, err := http.NewRequest("POST", "/books", bytes.NewBuffer([]byte(`{"ISBN":"0140447938222"}`)))

		if err != nil {
			t.Errorf("Error occured, %s", err)
		}

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(bookHandler.Create)

		handler.ServeHTTP(rr, req)
		contains := strings.Contains(rr.Body.String(), clmock.Err.Error())
		if !contains && rr.Code != http.StatusBadRequest {
			t.Errorf("Expected error to be %s, got error: %s", clmock.Err.Error(), rr.Body.String())
		}
	})
	t.Run("Testing book creation", func(t *testing.T) {
		mockClient := StoreMock{}
		book := model.Book{
			1,
			"War and Peace (Penguin Classics)",
			"Tolstoy",
			"0140447938",
			"9780140447934",
			"OL7355422M",
			"5049015",
			2019,
			"Penguin Books",
			time.Time{},
			time.Time{},
		}
		mockClient.On("CreateBook", &book).Return(model.Book{
			1,
			"War and Peace (Penguin Classics)",
			"Tolstoy",
			"0140447938",
			"9780140447934",
			"OL7355422M",
			"5049015",
			2019,
			"Penguin Books",
			time.Time{},
			time.Time{},
		}, nil)
		mockClient.CreateBook(&book)
		mockClient.AssertExpectations(t)
	})

}

func TestGet(t *testing.T) {
	t.Run("Successfully retrieved book", func(t *testing.T) {
		mockClient := StoreMock{}
		mockClient.On("FindBookById", 1).Return(model.Book{
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          2019,
		}, nil)
		mockClient.FindBookById(1)
		mockClient.AssertExpectations(t)

	})

	t.Run("No such book ID in the database", func(t *testing.T) {
		mockClient := StoreMock{}
		mockClient.On("FindBookById", 5).Return(nil, errors.New("Can't find the book with that ID"))
		mockClient.FindBookById(5)
		mockClient.AssertExpectations(t)

	})
}

func TestDelete(t *testing.T) {
	t.Run("Book succesfully deleted", func(t *testing.T) {
		mockClient := StoreMock{}
		book := model.Book{

			Id:            2,
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2019,
		}

		mockClient.On("DeleteBook", &book).Return(model.Book{
			Id:            2,
			Title:         "test title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          2019,
		}, nil)
		mockClient.DeleteBook(&book)
		mockClient.AssertExpectations(t)

	})
}
