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
	Interactor loanReader
}

type loanReader interface {
	FindByID(ID int) (*model.Loan, error)
	FindAll() ([]*model.Loan, error)
	FindByUserID(id int) ([]*model.Loan, error)
	FindByBookID(id int) ([]*model.Loan, error)
}

type NewLoanHandler struct {
	Interactor loanWriter
}

type loanWriter interface {
	Borrow(userID int, bookID int) error
	Return(userID int, bookID int) error
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

func (n *NewLoanHandler) BorrowBook(w http.ResponseWriter, r *http.Request) *HTTPError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}
	err = n.Interactor.Borrow(10, id)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	w.Write([]byte("Loan successfully createad"))
	return nil
}

func (n *NewLoanHandler) ReturnBook(w http.ResponseWriter, r *http.Request) *HTTPError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}
	err = n.Interactor.Return(10, id)
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
