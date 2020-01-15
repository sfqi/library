package model

import (
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

func NewLoan(userId int, bookId int, uuid string, state LoanType) (*Loan, error) {
	return &Loan{
		UserID:        userId,
		BookID:        bookId,
		TransactionID: uuid,
		Type:          state,
	}, nil
}

func ReturnedLoan(userid int, bookid int, uuid string) (*Loan, error) {
	return NewLoan(userid, bookid, uuid, returned)
}

func BorrowedLoan(userid int, bookid int, uuid string) (*Loan, error) {
	return NewLoan(userid, bookid, uuid, borrowed)
}
