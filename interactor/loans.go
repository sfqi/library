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

type loansGenerator interface {
	GenerateUUID() (string, error)
}

type Loan struct {
	loansStore loansStore
	generator  loansGenerator
}

func NewLoan(loansStore loansStore, generator loansGenerator) *Loan {
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

func (l *Loan) CreateLoan(userID int, bookID int, state model.LoanType) error {
	loan, err := model.NewLoan(userID, bookID, state)
	uuid, err := l.generator.GenerateUUID()
	if err != nil {
		return err
	}
	loan.TransactionID = uuid
	return l.loansStore.CreateLoan(loan)
}
