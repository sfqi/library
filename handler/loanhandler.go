package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
)

type LoanHandler struct {
	Interactor loanInteractor
}

type loanInteractor interface {
	FindByID(ID int) (*model.Loan, error)
	FindAll() ([]*model.Loan, error)
	FindByUserID(id int) ([]*model.Loan, error)
	FindByBookID(id int) ([]*model.Loan, error)
}

func (l *LoanHandler) FindLoansByUserID(w http.ResponseWriter, r *http.Request) *HTTPError {

	id, err := strconv.Atoi(mux.Vars(r)["user_id"])
	if err != nil {
		return newHTTPError(http.StatusBadRequest, err)
	}

	loans, err := l.Interactor.FindByUserID(id)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	if len(loans) == 0 {
		return newHTTPError(http.StatusNotFound, errors.New("No loans for given user id can be found"))
	}

	var loanResponses []*dto.LoanResponse
	for _, loan := range loans {
		loanResponses = append(loanResponses, toLoanResponse(*loan))
	}

	err = json.NewEncoder(w).Encode(loanResponses)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, err)
	}
	return nil
}

func toLoanResponse(b model.Loan) *dto.LoanResponse {
	return dto.CreateLoanResponse(b.ID, b.TransactionID, b.UserID, b.BookID, b.PrintType())
}
