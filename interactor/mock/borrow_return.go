package mock

import "github.com/stretchr/testify/mock"

type BookLoan struct {
	mock.Mock
}

func (b *BookLoan) Borrow(userID int, bookID int) error {
	args := b.Called(userID, bookID)
	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}

func (b *BookLoan) Return(userID int, bookID int) error {
	args := b.Called(userID, bookID)
	if args.Get(0) != nil {
		return args.Error(0)
	}

	return nil
}
