package interactor

import (
	"errors"
	"github.com/sfqi/library/domain/model"
	repomock "github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindByID(t *testing.T) {

}

func TestFindAll(t *testing.T) {

}

func TestFindByUserID(t *testing.T) {
	assert := assert.New(t)
	t.Run("Successfully retrieved loans by user id", func(t *testing.T) {
		store := &repomock.Store{}

		l := NewLoan(store)

		store.On("FindLoansByUserID", 1).Return([]*model.Loan{
			{
				ID:            2,
				TransactionID: "asdasd22",
				UserID:        1,
				BookID:        1,
				Type:          1,
			},
		}, nil)

		loans, err := l.FindByUserID(1)

		expectedLoans := []*model.Loan{
			{
				ID:            2,
				TransactionID: "asdasd22",
				UserID:        1,
				BookID:        1,
				Type:          1,
			},
		}

		assert.Equal(loans, expectedLoans)
		assert.NoError(err)

	})

	t.Run("Cannot retrieve loans by user id", func(t *testing.T) {
		store := &repomock.Store{}
		storeError := errors.New("Error finding loans with given user ID from database")

		l := NewLoan(store)

		store.On("FindLoansByUserID", 10).Return(nil, storeError)

		loans, err := l.FindByUserID(10)

		assert.Nil(loans)
		assert.Equal(err, storeError)
	})
}

func TestFindByBookID(t *testing.T) {
	assert := assert.New(t)
	t.Run("Successfully retrieved loans by book id", func(t *testing.T) {
		store := &repomock.Store{}

		l := NewLoan(store)

		store.On("FindLoansByBookID", 1).Return([]*model.Loan{
			{
				ID:            2,
				TransactionID: "asdasd22",
				UserID:        1,
				BookID:        1,
				Type:          1,
			},
		}, nil)

		loans, err := l.FindByBookID(1)

		expectedLoans := []*model.Loan{
			{
				ID:            2,
				TransactionID: "asdasd22",
				UserID:        1,
				BookID:        1,
				Type:          1,
			},
		}

		assert.Equal(loans, expectedLoans)
		assert.NoError(err)

	})

	t.Run("Cannot retrieve loans by book id", func(t *testing.T) {
		store := &repomock.Store{}
		storeError := errors.New("Error finding loans with given book ID from database")

		l := NewLoan(store)

		store.On("FindLoansByBookID", 10).Return(nil, storeError)

		loans, err := l.FindByBookID(10)

		assert.Nil(loans)
		assert.Equal(err, storeError)
	})
}
