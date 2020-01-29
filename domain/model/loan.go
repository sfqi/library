package model

import (
	"errors"
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

func (l *Loan) PrintType() (string, error) {
	switch l.Type {
	case 0:
		return "borrowed", nil
	case 1:
		return "returned", nil
	default:
		return "", errors.New("unknown loan type")
	}
}

func ReturnedLoan(userid int, bookid int, uuid string) *Loan {
	return newLoan(userid, bookid, uuid, returned)
}

func BorrowedLoan(userid int, bookid int, uuid string) *Loan {
	return newLoan(userid, bookid, uuid, borrowed)
}
