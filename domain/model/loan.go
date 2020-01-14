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

func NewLoan(userId int, bookId int, state LoanType) (*Loan, error) {
	return &Loan{
		UserID: userId,
		BookID: bookId,
		Type:   state,
	}, nil
}

func ReturnedLoan(userid int, bookid int) (*Loan, error) {
	return NewLoan(userid, bookid, returned)
}

func BorrowedLoan(userid int, bookid int) (*Loan, error) {
	return NewLoan(userid, bookid, borrowed)
}
