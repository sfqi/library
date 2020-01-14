package interactor

import (
	"errors"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/interactor/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindByID(t *testing.T) {

}

func TestFindAll(t *testing.T) {

}

func TestFindByUserID(t *testing.T) {
	assert := assert.New(t)
	t.Run("successfully retrieved loans by user id", func(t *testing.T) {
		loanstore := &mock.Loan{}
		loanInteractor := NewLoan(loanstore)
		loanstore.On("FindLoansByUserID", 1).Return([]*model.Loan{
			{
				ID:            1,
				TransactionID: "testTransaction12",
				UserID:        1,
				BookID:        1,
				Type:          1,
			},
		}, nil)

		loans, err := loanInteractor.FindByUserID(1)

		expectedLoans := []*model.Loan{
			{
				ID:            1,
				TransactionID: "testTransaction12",
				UserID:        1,
				BookID:        1,
				Type:          1,
			},
		}
		assert.Equal(loans, expectedLoans)
		assert.NoError(err)

	})
	t.Run("Cannot retrieve loans by user id", func(t *testing.T) {
		loanstore := &mock.Loan{}
		loanInteractor := NewLoan(loanstore)
		storeError := errors.New("error finding loans with given user ID")
		loanstore.On("FindLoansByUserID", 13).Return(nil, storeError)

		loans, err := loanInteractor.FindByUserID(13)

		assert.Nil(loans)
		assert.Equal(err, storeError)
	})
}

func TestFindByBookID(t *testing.T) {
	assert := assert.New(t)
	t.Run("successfully retrieved loans by book id", func(t *testing.T) {
		loanstore := &mock.Loan{}
		loanInteractor := NewLoan(loanstore)
		loanstore.On("FindLoansByBookID", 1).Return([]*model.Loan{
			{
				ID:            1,
				TransactionID: "testTransaction12",
				UserID:        1,
				BookID:        1,
				Type:          1,
			},
		}, nil)

		loans, err := loanInteractor.FindByBookID(1)

		expectedLoans := []*model.Loan{
			{
				ID:            1,
				TransactionID: "testTransaction12",
				UserID:        1,
				BookID:        1,
				Type:          1,
			},
		}
		assert.Equal(loans, expectedLoans)
		assert.NoError(err)

	})
	t.Run("Cannot retrieve loans by book id", func(t *testing.T) {
		loanstore := &mock.Loan{}
		loanInteractor := NewLoan(loanstore)
		storeError := errors.New("error finding loans with given book ID")
		loanstore.On("FindLoansByBookID", 13).Return(nil, storeError)

		loans, err := loanInteractor.FindByBookID(13)

		assert.Nil(loans)
		assert.Equal(err, storeError)
	})
}
