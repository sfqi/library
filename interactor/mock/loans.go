package mock

import (
	"github.com/sfqi/library/domain/model"
	"github.com/stretchr/testify/mock"
)

type Loans struct {
	mock.Mock
}

func (l *Loans) FindByID(ID int) ([]*model.Loan, error) {
	args := l.Called(ID)
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Loan), nil
	}
	return nil, args.Error(1)
}

func (l *Loans) FindAll() ([]*model.Loan, error) {
	args := l.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Loan), nil
	}
	return nil, args.Error(1)
}

func (l *Loans) CreateLoan(userId int, bookId int, state model.LoanType) error {
	args := l.Called(userId, bookId, state)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
