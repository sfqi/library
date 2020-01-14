package model

import (
	"github.com/google/uuid"
	"time"
)

type LoanType int

const (
	borrowed LoanType = 0
	returned LoanType = 1
)

type Loan struct {
	ID            int
	TransactionID string
	UserID        int
	BookID        int
	Type          LoanType
	CreatedAt     time.Time
}

func NewLoan(userId int, bookId int, state LoanType) (*Loan, error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	return &Loan{
		TransactionID: uuid.String(),
		UserID:        userId,
		BookID:        bookId,
		Type:          state,
	}, nil
}

func ReturnedLoan(userid int, bookid int) (*Loan, error) {
	return NewLoan(userid, bookid, returned)
}

func BorrowedLoan(userid int, bookid int) (*Loan, error) {
	return NewLoan(userid, bookid, borrowed)
}
