package mock

import (
	"github.com/sfqi/library/domain/model"
	"github.com/stretchr/testify/mock"
)

type Loan struct {
	mock.Mock
}

func (l *Loan) FindLoansByBookID(id int) ([]*model.Loan, error) {
	args := l.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Loan), nil
	}
	return nil, args.Error(1)
}

func (l *Loan) FindLoansByUserID(id int) ([]*model.Loan, error) {
	args := l.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Loan), nil
	}
	return nil, args.Error(1)
}

func (l *Loan) FindLoanByID(id int) (*model.Loan, error) {
	return nil, nil
}

func (l *Loan) FindAllLoans() ([]*model.Loan, error) {
	return nil, nil
}
