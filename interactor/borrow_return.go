package interactor

import (
	"errors"
	"github.com/sfqi/library/domain/model"
)

type LoanWriter struct {
	store     Store
	generator uuidGenerator
}

type uuidGenerator interface {
	Do() (string, error)
}

func NewBookLoan(borrowReturn Store, generator uuidGenerator) *LoanWriter {
	return &LoanWriter{
		store:     borrowReturn,
		generator: generator,
	}
}

func (l *LoanWriter) Borrow(userID int, book *model.Book) (*model.Loan, error) {
	tx := l.store.Transaction()
	uuid, err := l.generator.Do()
	if err != nil {
		tx.Rollback()
		return nil, err
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

func (l *LoanWriter) Return(userID int, book *model.Book) (*model.Loan, error) {
	tx := l.store.Transaction()
	uuid, err := l.generator.Do()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

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
