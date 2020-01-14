package interactor_test

import (
	"errors"
	"fmt"
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
		assert.Equal(loans, expectedLoans)
		assert.NoError(err)
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

func TestLoan_Create(t *testing.T) {
	assert := assert.New(t)
	t.Run("Loan successfully saved in database", func(t *testing.T) {
		store := &repomock.Store{}
		l := interactor.NewLoan(store)

		loan, err := model.BorrowedLoan(1, 1)
		fmt.Println(loan)
		store.On("CreateLoan", loan).Return(nil)

		err = l.CreateLoan(loan.UserID, loan.BookID, loan.Type)
		assert.NoError(err)

	})
}

func TestFindByUserID(t *testing.T) {

}

func TestFindByBookID(t *testing.T) {

}