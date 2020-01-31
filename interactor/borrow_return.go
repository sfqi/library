package interactor

import (
	"errors"
	"github.com/sfqi/library/domain/model"
)

type newLoan interface {
	CreateLoan(*model.Loan) error
	UpdateBook(*model.Book) error
}

type LoanWriter struct {
	store     newLoan
	generator uuidGenerator
}

type uuidGenerator interface {
	Do() (string, error)
}

func NewBookLoan(borrowReturn newLoan, generator uuidGenerator) *LoanWriter {
	return &LoanWriter{
		store:     borrowReturn,
		generator: generator,
	}
}

func (l *LoanWriter) Borrow(userID int, book *model.Book) (*model.Loan, error) {
	uuid, err := l.generator.Do()
	if err != nil {
		return nil, err
	}

	if book.Available != true {
		return nil, errors.New("Book is not available")
	}
	book.Available = false
	err = l.store.UpdateBook(book)
	if err != nil {
		return nil, err
	}
	loan := model.BorrowedLoan(userID, book.Id, uuid)
	return loan, l.store.CreateLoan(loan)
}

func (l *LoanWriter) Return(userID int, book *model.Book) (*model.Loan, error) {
	uuid, err := l.generator.Do()
	if err != nil {
		return nil, err
	}

	book.Available = true
	err = l.store.UpdateBook(book)
	if err != nil {
		return nil, err
	}

	loan := model.ReturnedLoan(userID, book.Id, uuid)
	return loan, l.store.CreateLoan(loan)
}
