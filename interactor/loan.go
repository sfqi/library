package interactor

import (
	"github.com/sfqi/library/domain/model"
)

type loansStore interface {
	FindLoanByID(int) (*model.Loan, error)
	FindAllLoans() ([]*model.Loan, error)
	FindLoansByBookID(int) ([]*model.Loan, error)
	FindLoansByUserID(int) ([]*model.Loan, error)
}

type Loan struct {
	loanStore loansStore
}

func NewLoan(loansStore loansStore) *Loan {
	return &Loan{
		loanStore: loansStore,
	}
}

func (l *Loan) FindByID(ID int) (*model.Loan, error) {
	return l.loanStore.FindLoanByID(ID)
}

func (l *Loan) FindAll() ([]*model.Loan, error) {
	return l.loanStore.FindAllLoans()
}
