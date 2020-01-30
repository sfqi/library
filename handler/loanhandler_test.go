package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler"
	"github.com/sfqi/library/handler/dto"
	imock "github.com/sfqi/library/interactor/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexLoan(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Successfully returned loans", func(t *testing.T) {
		interactor := &imock.Loan{}
		loanHandler := handler.ReadLoanHandler{}

		loansExpected := []model.Loan{
			{
				ID:            1,
				TransactionID: "12",
				UserID:        1,
				BookID:        1,
				Type:          1,
				CreatedAt:     time.Time{},
			},
			{
				ID:            2,
				TransactionID: "13",
				UserID:        2,
				BookID:        2,
				Type:          2,
				CreatedAt:     time.Time{},
			},
		}

		req, err := http.NewRequest("GET", "/loans", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		interactor.On("FindAll").Return([]*model.Loan{
			{
				ID:            1,
				TransactionID: "12",
				UserID:        1,
				BookID:        1,
				Type:          1,
				CreatedAt:     time.Time{},
			},
			{
				ID:            2,
				TransactionID: "13",
				UserID:        2,
				BookID:        2,
				Type:          2,
				CreatedAt:     time.Time{},
			}}, nil)

		loanHandler.Interactor = interactor
		httpError := loanHandler.Index(rr, req)
		assert.Nil(httpError)
		assert.Equal(http.StatusOK, rr.Code)

		var response []model.Loan
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(err)

		assert.Equal(loansExpected, response)
	})
	t.Run("Error retrieving loans", func(t *testing.T) {
		interactor := &imock.Loan{}
		loanHandler := handler.ReadLoanHandler{}

		req, err := http.NewRequest("GET", "/loans", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		interactor.On("FindAll").Return(nil, errors.New("Error finding loans"))
		loanHandler.Interactor = interactor

		handler := loanHandler.Index
		httperror := handler(rr, req)

		expectedResponse := "HTTP 500: Error finding loans"
		assert.Equal(expectedResponse, httperror.Error())

		assert.Equal(http.StatusInternalServerError, httperror.Code())
	})
}

func TestLoanHandler_FindLoansByBookID(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Successfully returned loans for given book id", func(t *testing.T) {
		interactor := &imock.Loan{}

		expectedLoans := []*dto.LoanResponse{
			{
				ID:            1,
				TransactionID: "asddsa12",
				UserID:        1,
				BookID:        2,
				Type:          "returned",
			},
			{
				ID:            2,
				TransactionID: "dddccf12dc13",
				UserID:        2,
				BookID:        2,
				Type:          "borrowed",
			},
		}

		req, err := http.NewRequest("GET", "/books/2/loans", nil)
		require.NoError(err)

		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		interactor.On("FindByBookID", 2).Return([]*model.Loan{
			{
				ID:            1,
				TransactionID: "asddsa12",
				UserID:        1,
				BookID:        2,
				Type:          1,
			},
			{
				ID:            2,
				TransactionID: "dddccf12dc13",
				UserID:        2,
				BookID:        2,
				Type:          0,
			},
		}, nil)
		loanHandler := handler.ReadLoanHandler{Interactor: interactor}
		httpError := loanHandler.FindLoansByBookID(rr, req)
		assert.Nil(httpError)

		assert.Equal(http.StatusOK, rr.Code)

		var loanResponses []*dto.LoanResponse
		err = json.NewDecoder(rr.Body).Decode(&loanResponses)
		assert.NoError(err)

		assert.Equal(loanResponses, expectedLoans)

	})
	t.Run("error converting book id to integer", func(t *testing.T) {
		interactor := &imock.Loan{}

		expectedError := "HTTP 404: strconv.Atoi: parsing \"ww\": invalid syntax"

		req, err := http.NewRequest("GET", "/books/ww/loans", nil)
		require.NoError(err)

		params := map[string]string{"id": "ww"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		loanHandler := handler.ReadLoanHandler{Interactor: interactor}
		httpError := loanHandler.FindLoansByBookID(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusNotFound, httpError.Code())

		assert.NoError(err)

		assert.Equal(expectedError, httpError.Error())

	})
	t.Run("error for given id returned from database", func(t *testing.T) {
		interactor := &imock.Loan{}

		expectedError := "HTTP 500: Internal server error"

		req, err := http.NewRequest("GET", "/books/12/loans", nil)
		require.NoError(err)

		params := map[string]string{"id": "-2"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		interactor.On("FindByBookID", -2).Return(nil, errors.New("Internal server error"))
		loanHandler := handler.ReadLoanHandler{Interactor: interactor}
		httpError := loanHandler.FindLoansByBookID(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusInternalServerError, httpError.Code())

		assert.NoError(err)

		assert.Equal(expectedError, httpError.Error())

	})
}

func TestLoanHandler_FindLoansByUserID(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Successfully returned loans for given user id", func(t *testing.T) {
		interactor := &imock.Loan{}

		expectedLoans := []*dto.LoanResponse{
			{
				ID:            1,
				TransactionID: "gen123",
				UserID:        1,
				BookID:        2,
				Type:          "returned",
			},
			{
				ID:            2,
				TransactionID: "gen345",
				UserID:        2,
				BookID:        2,
				Type:          "borrowed",
			},
		}

		req, err := http.NewRequest("GET", "/users/1/loans", nil)
		require.NoError(err)

		params := map[string]string{"id": "2"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		interactor.On("FindByUserID", 2).Return([]*model.Loan{
			{
				ID:            1,
				TransactionID: "gen123",
				UserID:        1,
				BookID:        2,
				Type:          1,
			},
			{
				ID:            2,
				TransactionID: "gen345",
				UserID:        2,
				BookID:        2,
				Type:          0,
			},
		}, nil)

		loanHandler := handler.ReadLoanHandler{Interactor: interactor}
		httpError := loanHandler.FindLoansByUserID(rr, req)
		assert.Nil(httpError)

		assert.Equal(http.StatusOK, rr.Code)

		var loanResponses []*dto.LoanResponse
		err = json.NewDecoder(rr.Body).Decode(&loanResponses)
		require.NoError(err)

		assert.Equal(loanResponses, expectedLoans)

	})
	t.Run("error converting user id to integer", func(t *testing.T) {
		interactor := &imock.Loan{}

		expectedError := "HTTP 404: strconv.Atoi: parsing \"ww\": invalid syntax"

		req, err := http.NewRequest("GET", "/users/ww/loans", nil)
		require.NoError(err)

		params := map[string]string{"id": "ww"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		loanHandler := handler.ReadLoanHandler{Interactor: interactor}
		httpError := loanHandler.FindLoansByUserID(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusNotFound, httpError.Code())

		assert.Equal(expectedError, httpError.Error())
	})
	t.Run("user has no loans yet", func(t *testing.T) {
		interactor := &imock.Loan{}

		req, err := http.NewRequest("GET", "/users/12/loans", nil)
		require.NoError(err)

		params := map[string]string{"id": "12"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()
		expected := []*model.Loan{}

		interactor.On("FindByUserID", 12).Return([]*model.Loan{}, nil)
		loanHandler := handler.ReadLoanHandler{Interactor: interactor}
		httpError := loanHandler.FindLoansByUserID(rr, req)

		assert.Nil(httpError)
		var response []*model.Loan
		err = json.NewDecoder(rr.Body).Decode(&response)
		require.NoError(err)

		assert.Equal(expected, response)
	})
	t.Run("error for given id returned from database", func(t *testing.T) {
		interactor := &imock.Loan{}

		expectedError := "HTTP 500: Internal server error"

		req, err := http.NewRequest("GET", "/users/12/loans", nil)
		require.NoError(err)

		params := map[string]string{"id": "-2"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		interactor.On("FindByUserID", -2).Return(nil, errors.New("Internal server error"))
		loanHandler := handler.ReadLoanHandler{Interactor: interactor}
		httpError := loanHandler.FindLoansByUserID(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusInternalServerError, httpError.Code())

		assert.Equal(expectedError, httpError.Error())
	})
}

func TestBorrowReturnHandler_BorrowBook(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Successfully borrowed book", func(t *testing.T) {
		interactor := &imock.BookLoan{}
		book := &model.Book{Id: 5, Available: true}
		expectedResponse := "Loan successfully createad"

		req, err := http.NewRequest("GET", "/books/5/borrow", nil)
		require.NoError(err)

		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)
		params := map[string]string{"id": "5"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		interactor.On("Borrow", 10, book).Return(nil)
		bookLoanHandler := handler.WriteLoanHandler{Interactor: interactor}
		httpError := bookLoanHandler.BorrowBook(rr, req)
		assert.Nil(httpError)

		assert.Equal(http.StatusOK, rr.Code)

		assert.Equal(expectedResponse, rr.Body.String())
	})
	t.Run("Book from context in nil", func(t *testing.T) {

		expectedError := "HTTP 404: Book is not found: "

		req, err := http.NewRequest("GET", "/books/7/borrow", nil)
		require.NoError(err)
		ctx := context.WithValue(req.Context(), "book", nil)
		req = req.WithContext(ctx)

		params := map[string]string{"id": "7"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		bookLoanHandler := handler.WriteLoanHandler{Interactor: nil}
		httpError := bookLoanHandler.BorrowBook(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusNotFound, httpError.Code())

		assert.Equal(expectedError, httpError.Error())
	})
	t.Run("error borrowing book in database", func(t *testing.T) {
		interactor := &imock.BookLoan{}
		expectedError := "HTTP 500: Error borrowing book: "
		book := &model.Book{Id: 5, Available: true}
		req, err := http.NewRequest("GET", "/books/5/borrow", nil)
		require.NoError(err)

		params := map[string]string{"id": "5"}
		req = mux.SetURLVars(req, params)
		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		interactor.On("Borrow", 10, book).Return(errors.New("Error borrowing book"))
		bookLoanHandler := handler.WriteLoanHandler{Interactor: interactor}
		httpError := bookLoanHandler.BorrowBook(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusInternalServerError, httpError.Code())

		assert.Equal(expectedError, httpError.Error())
	})
}

func TestBorrowReturnHandler_ReturnBook(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	t.Run("Successfully returned book", func(t *testing.T) {
		interactor := &imock.BookLoan{}
		book := &model.Book{Id: 5, Available: true}

		expectedResponse := "Loan successfully createad"

		req, err := http.NewRequest("GET", "/books/5/return", nil)
		require.NoError(err)

		params := map[string]string{"id": "5"}
		req = mux.SetURLVars(req, params)
		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		interactor.On("Return", 10, book).Return(nil)
		bookLoanHandler := handler.WriteLoanHandler{Interactor: interactor}
		httpError := bookLoanHandler.ReturnBook(rr, req)
		assert.Nil(httpError)

		assert.Equal(http.StatusOK, rr.Code)

		assert.Equal(expectedResponse, rr.Body.String())
	})

	t.Run("Book from context is nil", func(t *testing.T) {

		expectedError := "HTTP 404: Book is not found: "

		req, err := http.NewRequest("GET", "/books/10/return", nil)
		require.NoError(err)

		params := map[string]string{"id": "10"}
		req = mux.SetURLVars(req, params)
		ctx := context.WithValue(req.Context(), "book", nil)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()

		bookLoanHandler := handler.WriteLoanHandler{Interactor: nil}
		httpError := bookLoanHandler.ReturnBook(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusNotFound, httpError.Code())

		assert.Equal(expectedError, httpError.Error())
	})

	t.Run("error returning book in database", func(t *testing.T) {
		interactor := &imock.BookLoan{}
		expectedError := "HTTP 500: Error returning book: "
		book := &model.Book{Id: 5, Available: true}
		req, err := http.NewRequest("GET", "/books/-4/return", nil)
		require.NoError(err)

		params := map[string]string{"id": "-4"}
		req = mux.SetURLVars(req, params)
		ctx := context.WithValue(req.Context(), "book", book)
		req = req.WithContext(ctx)
		rr := httptest.NewRecorder()

		interactor.On("Return", 10, book).Return(errors.New("Error returning book"))
		bookLoanHandler := handler.WriteLoanHandler{Interactor: interactor}
		httpError := bookLoanHandler.ReturnBook(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusInternalServerError, httpError.Code())

		assert.Equal(expectedError, httpError.Error())
	})

}
