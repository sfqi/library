package handler

import (
	"encoding/json"
	"net/http"

	"github.com/sfqi/library/domain/model"
)

type LoanHandler struct {
	Interactor loanInteractor
}

type loanInteractor interface {
	FindByID(ID int) ([]*model.Loan, error)
	FindAll() ([]*model.Loan, error)
	FindByUserID(id int) ([]*model.Loan, error)
	FindByBookID(id int) ([]*model.Loan, error)
}

func (l *LoanHandler) Index(w http.ResponseWriter, r *http.Request) *HTTPError {
	w.Header().Set("Content-Type", "application/json")
	allLoans, err := l.Interactor.FindAll()
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}

	err = json.NewEncoder(w).Encode(allLoans)
	if err != nil {
		return newHTTPError(http.StatusBadRequest, err)
	}
	return nil
}
