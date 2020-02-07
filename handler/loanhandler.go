package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
)

type ReadLoanHandler struct {
	interactor loanReader
}

type loanReader interface {
	FindByID(ID int) (*model.Loan, error)
	FindAll() ([]*model.Loan, error)
	FindByUserID(id int) ([]*model.Loan, error)
	FindByBookID(id int) ([]*model.Loan, error)
}

type WriteLoanHandler struct {
	interactor loanWriter
}

func NewWriteLoanHandler(interactor loanWriter) *WriteLoanHandler {
	return &WriteLoanHandler{interactor}
}
func NewReadLoanHandler(interactor loanReader) *ReadLoanHandler {
	return &ReadLoanHandler{interactor}
}

type loanWriter interface {
	Borrow(userID int, bookID int) (*model.Loan, error)
	Return(userID int, bookID int) (*model.Loan, error)
}

func (rl *ReadLoanHandler) Index(w http.ResponseWriter, r *http.Request) *HTTPError {
	w.Header().Set("Content-Type", "application/json")
	loans, err := rl.interactor.FindAll()
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}

	err = json.NewEncoder(w).Encode(loans)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	return nil
}

func (rl *ReadLoanHandler) FindLoansByBookID(w http.ResponseWriter, r *http.Request) *HTTPError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}

	loans, err := rl.interactor.FindByBookID(id)
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

func (rl *ReadLoanHandler) FindLoansByUserID(w http.ResponseWriter, r *http.Request) *HTTPError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}

	loans, err := rl.interactor.FindByUserID(id)
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
	loan, err := wl.interactor.Borrow(10, id)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	loanResponse, err := toLoanResponse(loan)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(loanResponse)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

func (wl *WriteLoanHandler) ReturnBook(w http.ResponseWriter, r *http.Request) *HTTPError {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return newHTTPError(http.StatusNotFound, err)
	}
	loan, err := wl.interactor.Return(10, id)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
	loanResponse, err := toLoanResponse(loan)
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(loanResponse)
	if err != nil {
		return newHTTPError(http.StatusInternalServerError, err)
	}
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
