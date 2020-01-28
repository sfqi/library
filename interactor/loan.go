package interactor

import (
	"github.com/sfqi/library/domain/model"
)

type loanStore interface {
	FindLoanByID(int) (*model.Loan, error)
	FindAllLoans() ([]*model.Loan, error)
	FindLoansByBookID(int) ([]*model.Loan, error)
	FindLoansByUserID(int) ([]*model.Loan, error)
}

type Loan struct {
	store loanStore
}

func NewLoan(loanStore loanStore) *Loan {
	return &Loan{
		store: loanStore,
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
