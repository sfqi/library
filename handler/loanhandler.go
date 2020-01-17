package handler

import (
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
	return nil
}
