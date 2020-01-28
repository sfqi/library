package handler

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
	"net/http"
	"strconv"
)

type ReadLoanHandler struct {
	Interactor loanReader
}

type loanReader interface {
	FindByID(ID int) (*model.Loan, error)
	FindAll() ([]*model.Loan, error)
	FindByUserID(id int) ([]*model.Loan, error)
	FindByBookID(id int) ([]*model.Loan, error)
}

type WriteLoanHandler struct {
	Interactor loanWriter
}

type loanWriter interface {
	Borrow(userID int, bookID int) error
	Return(userID int, bookID int) error
}

func (rl *ReadLoanHandler) FindLoansByBookID(w http.ResponseWriter, r *http.Request) *HTTPError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}

	loans, err := rl.Interactor.FindByBookID(id)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}

	loanResponses := []*dto.LoanResponse{}
	for _, l := range loans {
		loan, err := toLoanResponse(l)
		if err != nil {
			return newHTTPError(http.StatusInternalServerError, err)
		}
		loanResponses = append(loanResponses, loan)
	}
	err = json.NewEncoder(w).Encode(loanResponses)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func (wl *WriteLoanHandler) BorrowBook(w http.ResponseWriter, r *http.Request) *HTTPError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}
	err = wl.Interactor.Borrow(10, id)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	w.Write([]byte("Loan successfully createad"))
	return nil
}

func (wl *WriteLoanHandler) ReturnBook(w http.ResponseWriter, r *http.Request) *HTTPError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}
	err = wl.Interactor.Return(10, id)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	w.Write([]byte("Loan successfully createad"))
	return nil
}

func toLoanResponse(l *model.Loan) (*dto.LoanResponse, error) {
	loanType, err := l.PrintType()
	if err != nil {
		return nil, err
	}

	return &dto.LoanResponse{
		ID:            l.ID,
		TransactionID: l.TransactionID,
		UserID:        l.UserID,
		BookID:        l.BookID,
		Type:          loanType,
	}, nil
}
