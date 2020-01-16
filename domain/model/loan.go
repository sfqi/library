package model

import (
	"time"
)

type loanType int

const (
	borrowed loanType = 0
	returned loanType = 1
)

type Loan struct {
	ID            int
	TransactionID string
	UserID        int
	BookID        int
	Type          loanType
	CreatedAt     time.Time
}

func NewLoan(userId int, bookId int, uuid string, loanType loanType) (*Loan, error) {
	return &Loan{
		UserID:        userId,
		BookID:        bookId,
		TransactionID: uuid,
		Type:          loanType,
	}, nil
}

func ReturnedLoan(userid int, bookid int, uuid string) (*Loan, error) {
	return NewLoan(userid, bookid, uuid, returned)
}

func BorrowedLoan(userid int, bookid int, uuid string) (*Loan, error) {
	return NewLoan(userid, bookid, uuid, borrowed)
}
