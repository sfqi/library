package inmemory

import (
	"fmt"
	"github.com/sfqi/library/domain/model"
	"time"
)

func (db *DB) FindLoanById(id int) (*model.Loan, error) {
	loan, _, err := db.findLoanById(id)
	return loan, err
}

func (db *DB) findLoanById(id int) (*model.Loan, int, error) {
	for i, l := range db.loans {
		if l.Id == id {
			return &l, i, nil
		}
	}
	return nil, -1, fmt.Errorf("error while finding Loan")
}

func (db *DB) FindAllLoans() ([]*model.Loan, error) {
	pointers := make([]*model.Loan, len(db.loans))
	for i := 0; i < len(db.loans); i++ {
		pointers[i] = &db.loans[i]
	}
	fmt.Println(pointers)
	return pointers, nil
}

func (db *DB) CreateLoan(loan *model.Loan) error {
	db.Id++
	now := time.Now()
	loan.Created_at = now

	loan.Id = db.Id
	db.loans = append(db.loans, *loan)
	return nil
}

func (db *DB) UpdateLoan(toUpdate *model.Loan) error {
	loan, index, err := db.findLoanById(toUpdate.Id)
	if err != nil {
		return err
	}
	loan.Transaction_id = toUpdate.Transaction_id
	loan.Type = toUpdate.Type
	loan.Book_id = toUpdate.Book_id
	loan.User_id = toUpdate.User_id
	toUpdate = loan
	db.loans[index] = *loan
	return nil
}

func (db *DB) DeleteLoan(loan *model.Loan) error {
	_, loc, err := db.findLoanById(loan.Id)
	if err != nil {
		return err
	}
	db.loans = append(db.loans[:loc], db.loans[loc+1:]...)
	db.Id--
	return nil
}
