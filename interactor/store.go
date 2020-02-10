package interactor

import "github.com/sfqi/library/domain/model"

type Store interface {
	FindBookById(int) (*model.Book, error)
	CreateBook(*model.Book) error
	UpdateBook(*model.Book) error
	FindAllBooks() ([]*model.Book, error)
	DeleteBook(*model.Book) error
	FindLoanByID(int) (*model.Loan, error)
	FindAllLoans() ([]*model.Loan, error)
	CreateLoan(*model.Loan) error
	FindLoansByBookID(int) ([]*model.Loan, error)
	FindLoansByUserID(int) ([]*model.Loan, error)
	Transaction() Store
	Commit() error
	Rollback()
	FindBookByIDForUpdate(int) (*model.Book, error)
}
