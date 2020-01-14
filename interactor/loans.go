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

type Loan struct {
	loansStore loansStore
}

func NewLoan(loansStore loansStore) *Loan {
	return &Loan{
		loansStore: loansStore,
	}
}

func (l *Loan) FindByID(ID int) (*model.Loan, error) {

	return l.loansStore.FindLoanByID(ID)
}

func (l *Loan) FindAll() ([]*model.Loan, error) {

	return l.loansStore.FindAllLoans()
}

func (l *Loan) CreateLoan(userID int, bookID int, state model.LoanType) error {
	loan, err := model.NewLoan(userID, bookID, state)
	if err != nil {
		return err
	}
	return l.loansStore.CreateLoan(loan)
}
