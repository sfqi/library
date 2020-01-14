package mock

import (
	"github.com/sfqi/library/domain/model"
	"github.com/stretchr/testify/mock"
)

type Loan struct {
	mock.Mock
}

func (l *Loan) FindByUserID(id int) ([]*model.Loan, error) {
	args := l.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Loan), nil
	}
	return nil, args.Error(1)
}

func (l *Loan) FindByBookID(id int) ([]*model.Loan, error) {
	args := l.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Loan), nil
	}
	return nil, args.Error(1)
}
