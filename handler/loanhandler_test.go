package handler_test

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler"
	"github.com/sfqi/library/handler/dto"
	imock "github.com/sfqi/library/interactor/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"

	"testing"
)

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
		loanHandler := handler.LoanHandler{Interactor: interactor}
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

		loanHandler := handler.LoanHandler{Interactor: interactor}
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
		loanHandler := handler.LoanHandler{Interactor: interactor}
		httpError := loanHandler.FindLoansByBookID(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusInternalServerError, httpError.Code())

		assert.NoError(err)

		assert.Equal(expectedError, httpError.Error())

	})
}

func TestBorrowReturnHandler_BorrowBook(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Successfully borrowed book", func(t *testing.T) {
		interactor := &imock.BookLoan{}

		expectedResponse := "Loan successfully createad"

		req, err := http.NewRequest("GET", "/books/5/borrow", nil)
		require.NoError(err)

		params := map[string]string{"id": "5"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		interactor.On("Borrow", 10, 5).Return(nil)
		bookLoanHandler := handler.BorrowReturnHandler{Interactor: interactor}
		httpError := bookLoanHandler.BorrowBook(rr, req)
		assert.Nil(httpError)

		assert.Equal(http.StatusOK, rr.Code)

		assert.Equal(expectedResponse, rr.Body.String())
	})
	t.Run("error converting id to integer", func(t *testing.T) {

		expectedError := "HTTP 404: strconv.Atoi: parsing \"ww\": invalid syntax"

		req, err := http.NewRequest("GET", "/books/ww/borrow", nil)
		require.NoError(err)

		params := map[string]string{"id": "ww"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		bookLoanHandler := handler.BorrowReturnHandler{Interactor: nil}
		httpError := bookLoanHandler.BorrowBook(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusNotFound, httpError.Code())

		assert.Equal(expectedError, httpError.Error())
	})
	t.Run("error borrowing book in database", func(t *testing.T) {
		interactor := &imock.BookLoan{}
		expectedError := "HTTP 500: Error borrowing book"

		req, err := http.NewRequest("GET", "/books/-4/borrow", nil)
		require.NoError(err)

		params := map[string]string{"id": "-4"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		interactor.On("Borrow", 10, -4).Return(errors.New("Error borrowing book"))
		bookLoanHandler := handler.BorrowReturnHandler{Interactor: interactor}
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

		expectedResponse := "Loan successfully createad"

		req, err := http.NewRequest("GET", "/books/5/return", nil)
		require.NoError(err)

		params := map[string]string{"id": "5"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		interactor.On("Return", 10, 5).Return(nil)
		bookLoanHandler := handler.BorrowReturnHandler{Interactor: interactor}
		httpError := bookLoanHandler.ReturnBook(rr, req)
		assert.Nil(httpError)

		assert.Equal(http.StatusOK, rr.Code)

		assert.Equal(expectedResponse, rr.Body.String())
	})

	t.Run("error converting id to integer", func(t *testing.T) {

		expectedError := "HTTP 404: strconv.Atoi: parsing \"ww\": invalid syntax"

		req, err := http.NewRequest("GET", "/books/ww/return", nil)
		require.NoError(err)

		params := map[string]string{"id": "ww"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		bookLoanHandler := handler.BorrowReturnHandler{Interactor: nil}
		httpError := bookLoanHandler.ReturnBook(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusNotFound, httpError.Code())

		assert.Equal(expectedError, httpError.Error())
	})

	t.Run("error returning book in database", func(t *testing.T) {
		interactor := &imock.BookLoan{}
		expectedError := "HTTP 500: Error returning book"

		req, err := http.NewRequest("GET", "/books/-4/return", nil)
		require.NoError(err)

		params := map[string]string{"id": "-4"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		interactor.On("Return", 10, -4).Return(errors.New("Error returning book"))
		bookLoanHandler := handler.BorrowReturnHandler{Interactor: interactor}
		httpError := bookLoanHandler.ReturnBook(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusInternalServerError, httpError.Code())

		assert.Equal(expectedError, httpError.Error())
	})

}
