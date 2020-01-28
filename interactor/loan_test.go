package interactor_test

import (
	"errors"
	"testing"
	"time"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/interactor"
	repomock "github.com/sfqi/library/repository/mock"
	"github.com/stretchr/testify/assert"
)

func TestFindByID(t *testing.T) {
	assert := assert.New(t)
	t.Run("Successfully retrieved loan", func(t *testing.T) {
		store := &repomock.Store{}

		l := interactor.NewLoan(store)

		store.On("FindLoanByID", 1).Return(&model.Loan{
			ID:            1,
			TransactionID: "12",
			UserID:        1,
			BookID:        1,
			Type:          22,
			CreatedAt:     time.Time{},
		}, nil)

		loan, err := l.FindByID(1)

		expectedLoan := &model.Loan{
			ID:            1,
			TransactionID: "12",
			UserID:        1,
			BookID:        1,
			Type:          22,
			CreatedAt:     time.Time{},
		}

		assert.Equal(loan, expectedLoan)
		assert.NoError(err)

	})

	t.Run("Cannot retrieve loan", func(t *testing.T) {
		store := &repomock.Store{}
		storeError := errors.New("Error finding loan ID from database")

		l := interactor.NewLoan(store)

		store.On("FindLoanByID", 12).Return(nil, storeError)

		loan, err := l.FindByID(12)

		assert.Nil(loan)
		assert.Equal(err, storeError)
	})
}

func TestFindAllLoans(t *testing.T) {
	assert := assert.New(t)
	t.Run("Successfully returned loans", func(t *testing.T) {
		store := &repomock.Store{}

		l := interactor.NewLoan(store)

		store.On("FindAllLoans").Return([]*model.Loan{
			{
				ID:            1,
				TransactionID: "12",
				UserID:        1,
				BookID:        1,
				Type:          22,
				CreatedAt:     time.Time{},
			},
			{
				ID:            2,
				TransactionID: "13",
				UserID:        2,
				BookID:        2,
				Type:          23,
				CreatedAt:     time.Time{},
			}}, nil)

		expectedLoans := []*model.Loan{
			{
				ID:            1,
				TransactionID: "12",
				UserID:        1,
				BookID:        1,
				Type:          22,
				CreatedAt:     time.Time{},
			},
			{
				ID:            2,
				TransactionID: "13",
				UserID:        2,
				BookID:        2,
				Type:          23,
				CreatedAt:     time.Time{},
			},
		}

		loans, err := l.FindAll()
		assert.NoError(err)
		assert.Equal(expectedLoans, loans)

	})

	t.Run("Error retrieving loans", func(t *testing.T) {
		store := &repomock.Store{}
		storeError := errors.New("Error finding loans")

		l := interactor.NewLoan(store)
		store.On("FindAllLoans").Return(nil, storeError)

		loans, err := l.FindAll()
		assert.Nil(loans)
		assert.Equal(err, storeError)
	})

}

func TestFindByUserID(t *testing.T) {
	assert := assert.New(t)
	t.Run("Successfully retrieved loans by user id", func(t *testing.T) {
		store := &repomock.Store{}

		l := interactor.NewLoan(store)

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

		l := interactor.NewLoan(store)

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

		l := interactor.NewLoan(store)

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

		l := interactor.NewLoan(store)

		store.On("FindLoansByBookID", 10).Return(nil, storeError)

		loans, err := l.FindByBookID(10)

		assert.Nil(loans)
		assert.Equal(err, storeError)
	})
}
