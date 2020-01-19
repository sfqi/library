package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
	"net/http"
	"strconv"
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

func (l *LoanHandler) FindLoansByBookID(w http.ResponseWriter, r *http.Request) *HTTPError {
	id, err := strconv.Atoi(mux.Vars(r)["book_id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}

	loans, err := l.Interactor.FindByBookID(id)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	if len(loans) == 0 {
		err = json.NewEncoder(w).Encode(loans)
	}
	var loanResponses []*dto.LoanResponse
	for _, loan := range loans {
		loanResponses = append(loanResponses, toLoanResponse(*loan))
	}
	err = json.NewEncoder(w).Encode(loanResponses)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func toLoanResponse(b model.Loan) *dto.LoanResponse {
	return dto.CreateLoanResponse(b.ID, b.TransactionID, b.UserID, b.BookID, b.PrintType())
}
