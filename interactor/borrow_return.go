package interactor

import "github.com/sfqi/library/domain/model"

type borrowReturn interface {
	CreateLoan(*model.Loan) error
}

type BookLoan struct {
	store     borrowReturn
	generator uuidGenerator
}

type uuidGenerator interface {
	Do() (string, error)
}

func NewBookLoan(borrowReturn borrowReturn, generator uuidGenerator) *BookLoan {
	return &BookLoan{
		store:     borrowReturn,
		generator: generator,
	}
}

func (l *BookLoan) Borrow(userID int, bookID int) error {
	uuid, err := l.generator.Do()
	if err != nil {
		return err
	}

	loan := model.BorrowedLoan(userID, bookID, uuid)
	return l.store.CreateLoan(loan)
}

func (l *BookLoan) Return(userID int, bookID int) error {
	uuid, err := l.generator.Do()
	if err != nil {
		return err
	}

	loan := model.ReturnedLoan(userID, bookID, uuid)
	return l.store.CreateLoan(loan)
}
