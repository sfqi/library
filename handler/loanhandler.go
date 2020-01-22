package handler

import (
	"encoding/json"
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
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}

	loans, err := l.Interactor.FindByUserID(id)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}

	loanResponses := []*dto.LoanResponse{}

	for _, loan := range loans {
		loans, err := toLoanResponse(loan)
		if err != nil {
			return newHTTPError(http.StatusNotFound, err)
		}
		loanResponses = append(loanResponses, loans)
	}

	err = json.NewEncoder(w).Encode(loanResponses)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}

func (l *LoanHandler) FindLoansByBookID(w http.ResponseWriter, r *http.Request) *HTTPError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}

	loans, err := l.Interactor.FindByBookID(id)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}

	loanResponses := []*dto.LoanResponse{}
	for _, loan := range loans {
		loans, err := toLoanResponse(loan)
		if err != nil {
			return newHTTPError(http.StatusNotFound, err)
		}
		loanResponses = append(loanResponses, loans)
	}

	err = json.NewEncoder(w).Encode(loanResponses)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func toLoanResponse(l *model.Loan) (*dto.LoanResponse, *HTTPError) {
	loanType, err := l.PrintType()
	if err != nil {
		return nil, newHTTPError(http.StatusNotFound, err)
	}

	loanResponse, err := dto.CreateLoanResponse(l.ID, l.TransactionID, l.UserID, l.BookID, loanType)
	if err != nil {
		return nil, newHTTPError(http.StatusNotFound, err)
	}
	return loanResponse, nil
}
