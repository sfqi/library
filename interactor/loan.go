package interactor

import (
	"github.com/sfqi/library/domain/model"
)

type loanStore interface {
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
	store     loanStore
	generator uuidGenerator
}

func NewLoan(loanStore loanStore, generator uuidGenerator) *Loan {
	return &Loan{
		store:     loanStore,
		generator: generator,
	}
}

func (l *Loan) FindByID(ID int) (*model.Loan, error) {
	return l.store.FindLoanByID(ID)
}

func (l *Loan) FindAll() ([]*model.Loan, error) {
	return l.store.FindAllLoans()
}

func (l *Loan) FindByUserID(id int) ([]*model.Loan, error) {
	return l.store.FindLoansByUserID(id)
}

func (l *Loan) FindByBookID(id int) ([]*model.Loan, error) {
	return l.store.FindLoansByBookID(id)
}

func (l *Loan) Borrow(userID int, bookID int) error {
	uuid, err := l.generator.Do()
	if err != nil {
		return err
	}

	loan := model.BorrowedLoan(userID, bookID, uuid)
	return l.store.CreateLoan(loan)
}

func (l *Loan) Return(userID int, bookID int) error {
	uuid, err := l.generator.Do()
	if err != nil {
		return err
	}

	loan := model.ReturnedLoan(userID, bookID, uuid)
	return l.store.CreateLoan(loan)
}
