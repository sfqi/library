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

func (l *Loan) PrintType() string {
	switch l.Type {
	case 0:
		return "borrowed"
	case 1:
		return "returned"
	default:
		return "unknown"
	}
}

func ReturnedLoan(userid int, bookid int, uuid string) *Loan {
	return newLoan(userid, bookid, uuid, returned)
}

func BorrowedLoan(userid int, bookid int, uuid string) *Loan {
	return newLoan(userid, bookid, uuid, borrowed)
}
