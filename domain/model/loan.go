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

func newLoan(userId int, bookId int, uuid string, loanType loanType) *Loan {
	return &Loan{
		UserID:        userId,
		BookID:        bookId,
		TransactionID: uuid,
		Type:          loanType,
	}
}

func ReturnedLoan(userid int, bookid int, uuid string) *Loan {
	return newLoan(userid, bookid, uuid, returned)
}

func BorrowedLoan(userid int, bookid int, uuid string) *Loan {
	return newLoan(userid, bookid, uuid, borrowed)
}
