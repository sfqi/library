package mock

import (
	"github.com/sfqi/library/domain/model"
	"github.com/stretchr/testify/mock"
)

type BookLoan struct {
	mock.Mock
}

func (b *BookLoan) Borrow(userID int, book *model.Book) error {
	args := b.Called(userID, book)
	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (b *BookLoan) Return(userID int, book *model.Book) error {
	args := b.Called(userID, book)
	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}
