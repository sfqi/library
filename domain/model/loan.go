package model

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

type loanType int
const(
	borrowed loanType = 0
	returned loanType = 1
)

type Loan struct{
	Id int
	TransactionID string
	UserID int
	BookID int
	Type loanType
	CreatedAt time.Time
}

func NewLoan() *Loan{
	uuid ,err := uuid.NewUUID()
	if err != nil{
		fmt.Println(err)
	}
	return &Loan{
		TransactionID: uuid.String(),
	}
}

func ReturnedLoan(userid int, bookid int) *Loan{
	loan := NewLoan()
	loan.CreatedAt = time.Now()
	loan.BookID = bookid
	loan.UserID=userid
	loan.Type = 1
	return loan
}

func BorrowedLoan(userid int, bookid int) *Loan{
	loan := NewLoan()
	loan.CreatedAt = time.Now()
	loan.BookID = bookid
	loan.UserID=userid
	loan.Type = 0
	return loan
}