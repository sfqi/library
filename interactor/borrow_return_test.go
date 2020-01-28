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

		loan := model.BorrowedLoan(1, 1, "gen123-gen321")

		generator.On("Do").Return("gen123-gen321", nil)

		store.On("CreateLoan", loan).Return(nil)
		l := interactor.NewBookLoan(store, generator)
		err := l.Borrow(loan.UserID, loan.BookID)
		assert.NoError(err)
	})

	t.Run("Error creating borrow loan in database", func(t *testing.T) {
		store := &repomock.Store{}
		generator := &uuid.Generator{}

		loan := model.BorrowedLoan(1, 1, "gen123-gen321")

		l := interactor.NewBookLoan(store, generator)
		storeError := errors.New("Error saving borrow loan in database")

		generator.On("Do").Return("gen123-gen321", nil)
		store.On("CreateLoan", loan).Return(storeError)

		err := l.Borrow(loan.UserID, loan.BookID)
		assert.Equal(err, storeError)
	})
}

func TestReturn(t *testing.T) {
	assert := assert.New(t)
	t.Run("Return loan successfully saved in database", func(t *testing.T) {
		store := &repomock.Store{}
		generator := &uuid.Generator{}

		loan := model.ReturnedLoan(1, 1, "")

		generator.On("Do").Return("", nil)

		store.On("CreateLoan", loan).Return(nil)
		l := interactor.NewBookLoan(store, generator)
		err := l.Return(loan.UserID, loan.BookID)
		assert.NoError(err)

	})
	t.Run("Error creating return loan in database", func(t *testing.T) {
		store := &repomock.Store{}
		generator := &uuid.Generator{}

		loan := model.ReturnedLoan(0, 0, "")

		l := interactor.NewBookLoan(store, generator)
		storeError := errors.New("Error saving return loan in database")

		generator.On("Do").Return("", nil)
		store.On("CreateLoan", loan).Return(storeError)

		err := l.Return(loan.UserID, loan.BookID)
		assert.Equal(err, storeError)
	})
}
