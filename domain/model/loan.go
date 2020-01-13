package model

import (
	"github.com/google/uuid"
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

func NewLoan(userId int, bookId int, state loanType) *Loan {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return nil
	}
	return &Loan{
		TransactionID: uuid.String(),
		UserID:        userId,
		BookID:        bookId,
		Type:          state,
		CreatedAt:     time.Now(),
	}
}

func ReturnedLoan(userid int, bookid int) *Loan {
	return NewLoan(userid, bookid, returned)
}

func BorrowedLoan(userid int, bookid int) *Loan {
	return NewLoan(userid, bookid, borrowed)
}
