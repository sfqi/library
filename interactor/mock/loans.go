package mock

import (
	"github.com/sfqi/library/domain/model"
	"github.com/stretchr/testify/mock"
)

type Loan struct {
	mock.Mock
}

func (l *Loan) FindByUserID(ID int) ([]*model.Loan, error) {
	args := l.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Loan), nil
	}
	return nil, args.Error(1)
}

func (l *Loan) FindByID(ID int) ([]*model.Loan, error) {
	args := l.Called(ID)
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Loan), nil
	}

	return nil, args.Error(1)
}

func (l *Loan) FindByBookID(ID int) ([]*model.Loan, error) {
	args := l.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Loan), nil
	}
	return nil, args.Error(1)
}

func (l *Loan) FindAll() ([]*model.Loan, error) {
	args := l.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Loan), nil
	}

	return nil, args.Error(1)
}

func (l *Loan) CreateLoan(userID int, bookID int) error {
	args := l.Called(userID, bookID)
	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}
