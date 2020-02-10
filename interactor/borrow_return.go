package interactor

import (
	"errors"

	"github.com/sfqi/library/domain/model"
)

type newLoan interface {
	CreateLoan(*model.Loan) error
	FindBookById(int) (*model.Book, error)
	UpdateBook(*model.Book) error
}

type LoanWriter struct {
	store         Store
	uuidGenerator uuidGenerator
}

type uuidGenerator interface {
	Do() (string, error)
}

func NewBookLoan(store Store, uuidGenerator uuidGenerator) *LoanWriter {
	return &LoanWriter{
		store:         store,
		uuidGenerator: uuidGenerator,
	}
}

func (l *LoanWriter) Borrow(userID int, bookID int) (*model.Loan, error) {
	tx := l.store.Transaction()
  
	uuid, err := l.uuidGenerator.Do()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	book, err := l.store.FindBookByIDForUpdate(bookID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("Error finding book")
	}

	if book.Available != true {
		tx.Rollback()
		return nil, errors.New("Book is not available")
	}
	book.Available = false
	err = tx.UpdateBook(book)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	loan := model.BorrowedLoan(userID, book.Id, uuid)
	err = tx.CreateLoan(loan)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return loan, nil
}

func (l *LoanWriter) Return(userID int, bookID int) (*model.Loan, error) {
	tx := l.store.Transaction()
  
	uuid, err := l.uuidGenerator.Do()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	book, err := l.store.FindBookById(bookID)
	book.Available = true
	err = tx.UpdateBook(book)
	if err != nil {
		return nil, err
	}

	loan := model.ReturnedLoan(userID, book.Id, uuid)
	err = tx.CreateLoan(loan)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return loan, nil
}
