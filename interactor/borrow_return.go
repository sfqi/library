package interactor

import "github.com/sfqi/library/domain/model"

type newLoan interface {
	CreateLoan(*model.Loan) error
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

func (l *LoanWriter) Borrow(userID int, bookID int) error {
	uuid, err := l.generator.Do()
	if err != nil {
		return err
	}

	loan := model.BorrowedLoan(userID, bookID, uuid)
	return l.store.CreateLoan(loan)
}

func (l *LoanWriter) Return(userID int, bookID int) error {
	uuid, err := l.generator.Do()
	if err != nil {
		return err
	}

	loan := model.ReturnedLoan(userID, bookID, uuid)
	return l.store.CreateLoan(loan)
}
