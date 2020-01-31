package interactor_test

import (
	"errors"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/interactor"
	repomock "github.com/sfqi/library/repository/mock"
	uuid "github.com/sfqi/library/service/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBorrow(t *testing.T) {
	assert := assert.New(t)
	t.Run("Borrow loan successfully saved in database", func(t *testing.T) {
		store := &repomock.Store{}
		generator := &uuid.Generator{}

		bookInUse := &model.Book{Id: 1, Available: true}
		store.On("Transaction").Return(store)
		store.On("UpdateBook", bookInUse).Return(nil)
		store.On("Commit").Return(nil)
		loan := model.BorrowedLoan(1, 1, "gen123-gen321")

		generator.On("Do").Return("gen123-gen321", nil)

		store.On("CreateLoan", loan).Return(nil)
		l := interactor.NewBookLoan(store, generator)
		loanB, err := l.Borrow(loan.UserID, bookInUse)
		assert.NoError(err)
		assert.Equal(loan, loanB)
	})
	t.Run("Book is already borrowed", func(t *testing.T) {
		store := &repomock.Store{}
		generator := &uuid.Generator{}
		expectedError := "Book is not available"
		bookInUse := &model.Book{Id: 1, Available: false}
		store.On("FindBookById", 1).Return(bookInUse, nil)

		generator.On("Do").Return("gen123-gen321", nil)

		l := interactor.NewBookLoan(store, generator)
		loan, err := l.Borrow(1, bookInUse)
		assert.Equal(expectedError, err.Error())
		assert.Nil(loan)
	})

	t.Run("Error creating borrow loan in database", func(t *testing.T) {
		store := &repomock.Store{}
		generator := &uuid.Generator{}
		bookInUse := &model.Book{Id: 1, Available: true}
		store.On("FindBookById", 1).Return(bookInUse, nil)
		store.On("UpdateBook", bookInUse).Return(nil)

		loan := model.BorrowedLoan(1, 1, "gen123-gen321")

		l := interactor.NewBookLoan(store, generator)
		storeError := errors.New("Error saving borrow loan in database")

		generator.On("Do").Return("gen123-gen321", nil)
		store.On("CreateLoan", loan).Return(storeError)

		loan, err := l.Borrow(loan.UserID, bookInUse)
		assert.Equal(err, storeError)
	})
}

func TestReturn(t *testing.T) {
	assert := assert.New(t)
	t.Run("Return loan successfully saved in database", func(t *testing.T) {
		store := &repomock.Store{}
		generator := &uuid.Generator{}
		bookInUse := &model.Book{Id: 1, Available: false}
		store.On("FindBookById", 1).Return(bookInUse, nil)
		store.On("UpdateBook", bookInUse).Return(nil)

		loan := model.ReturnedLoan(1, 1, "")

		generator.On("Do").Return("", nil)

		store.On("CreateLoan", loan).Return(nil)
		l := interactor.NewBookLoan(store, generator)
		loanR, err := l.Return(loan.UserID, bookInUse)
		assert.NoError(err)
		assert.Equal(loan, loanR)
	})
	t.Run("Error creating return loan in database", func(t *testing.T) {
		store := &repomock.Store{}
		generator := &uuid.Generator{}
		bookInUse := &model.Book{Id: 1, Available: false}
		store.On("FindBookById", 1).Return(bookInUse, nil)
		store.On("UpdateBook", bookInUse).Return(nil)

		loan := model.ReturnedLoan(0, 1, "")

		l := interactor.NewBookLoan(store, generator)
		storeError := errors.New("Error saving return loan in database")

		generator.On("Do").Return("", nil)
		store.On("CreateLoan", loan).Return(storeError)

		loan, err := l.Return(loan.UserID, bookInUse)
		assert.Equal(err, storeError)

	})
}
