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

		params := map[string]string{"book_id": "2"}
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

		params := map[string]string{"book_id": "ww"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()

		loanHandler := handler.LoanHandler{Interactor: interactor}
		httpError := loanHandler.FindLoansByBookID(rr, req)
		assert.NotNil(httpError)

		assert.Equal(http.StatusNotFound, httpError.Code())

		assert.NoError(err)

		assert.Equal(expectedError, httpError.Error())

	})
	t.Run("Book is not loaned yet", func(t *testing.T) {
		interactor := &imock.Loan{}

		req, err := http.NewRequest("GET", "/books/12/loans", nil)
		require.NoError(err)

		params := map[string]string{"book_id": "12"}
		req = mux.SetURLVars(req, params)

		rr := httptest.NewRecorder()
		expected := []*model.Loan{}

		interactor.On("FindByBookID", 12).Return([]*model.Loan{}, nil)
		loanHandler := handler.LoanHandler{Interactor: interactor}
		httpError := loanHandler.FindLoansByBookID(rr, req)
		assert.Nil(httpError)
		var response []*model.Loan
		err = json.NewDecoder(rr.Body).Decode(&response)

		assert.Equal(expected, response)
	})
	t.Run("error for given id returned from database", func(t *testing.T) {
		interactor := &imock.Loan{}

		expectedError := "HTTP 500: Internal server error"

		req, err := http.NewRequest("GET", "/books/12/loans", nil)
		require.NoError(err)

		params := map[string]string{"book_id": "-2"}
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
