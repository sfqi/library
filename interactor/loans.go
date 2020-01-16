package interactor

import (
	"github.com/sfqi/library/domain/model"
)

type loansStore interface {
	FindLoanByID(int) (*model.Loan, error)
	FindAllLoans() ([]*model.Loan, error)
	FindLoansByBookID(int) ([]*model.Loan, error)
	FindLoansByUserID(int) ([]*model.Loan, error)
	CreateLoan(*model.Loan) error
}

type uuidGenerator interface {
	Do() (string, error)
}

type Loan struct {
	loansStore loansStore
	generator  uuidGenerator
}

func NewLoan(loansStore loansStore, generator uuidGenerator) *Loan {
	return &Loan{
		loansStore: loansStore,
		generator:  generator,
	}
}

func (l *Loan) FindByID(ID int) (*model.Loan, error) {

	return l.loansStore.FindLoanByID(ID)
}

func (l *Loan) FindAll() ([]*model.Loan, error) {

	return l.loansStore.FindAllLoans()
}

func (l *Loan) Borrow(userID int, bookID int) error {
	uuid, err := l.generator.Do()
	if err != nil {
		return err
	}

	loan := model.BorrowedLoan(userID, bookID, uuid)
	return l.loansStore.CreateLoan(loan)
}

func (l *Loan) Return(userID int, bookID int) error {
	uuid, err := l.generator.Do()
	if err != nil {
		return err
	}

	loan := model.ReturnedLoan(userID, bookID, uuid)
	return l.loansStore.CreateLoan(loan)
}
