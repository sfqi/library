package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler"
	imock "github.com/sfqi/library/interactor/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIndexLoan(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	t.Run("Successfully returned loans", func(t *testing.T) {
		interactor := &imock.Loan{}
		loanHandler := handler.LoanHandler{}

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
		loanHandler := handler.LoanHandler{}

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
