package model

import (
	"github.com/google/uuid"
	"time"
)

type State int
const(
	borrow State = 0
	giveBack State = 1 // Not allowing me to call "return"
)

type Loan struct{
	Id int
	Transaction_id string
	User_id int
	Book_id int
	Type State
	Created_at time.Time
}

func NewLoan() *Loan{
	return &Loan{
		Transaction_id: uuid.New().String(),
	}
}