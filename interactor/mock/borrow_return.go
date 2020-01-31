package mock

import (
	"github.com/sfqi/library/domain/model"
	"github.com/stretchr/testify/mock"
)

type BookLoan struct {
	mock.Mock
}

func (b *BookLoan) Borrow(userID int, book *model.Book) (*model.Loan, error) {
	args := b.Called(userID, book)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Loan), nil
	}

	return nil, args.Error(1)
}

func (b *BookLoan) Return(userID int, book *model.Book) (*model.Loan, error) {
	args := b.Called(userID, book)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Loan), nil
	}

	return nil, args.Error(1)
}
