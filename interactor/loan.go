package interactor

import (
	"github.com/sfqi/library/domain/model"
)

type Loan struct {
	store Store
}

func NewLoan(loanStore Store) *Loan {
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
