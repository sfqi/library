package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sfqi/library/domain/model"
	i "github.com/sfqi/library/interfaces"
)

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type Store struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

func Open(config PostgresConfig) (*Store, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.Name)

	db, err := gorm.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	store := NewStore(db)
	return store, nil
}

func (store *Store) Close() error {
	return store.db.Close()
}

func (store *Store) FindBookById(id int) (*model.Book, error) {
	b := model.Book{}
	if err := store.db.First(&b, id).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (store *Store) CreateBook(book *model.Book) error {
	return store.db.Create(&book).Error

}

func (store *Store) UpdateBook(book *model.Book) error {
	return store.db.Save(&book).Error
}

func (store *Store) FindAllBooks() ([]*model.Book, error) {
	books := []*model.Book{}
	if err := store.db.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

func (store *Store) DeleteBook(book *model.Book) error {
	return store.db.Where("id = ?", book.Id).Delete(&model.Book{}).Error
}

func (store *Store) CreateLoan(loan *model.Loan) error {
	return store.db.Create(&loan).Error
}

func (store *Store) FindLoanByID(id int) (*model.Loan, error) {
	loan := &model.Loan{}
	if err := store.db.First(&loan, id).Error; err != nil {
		return nil, err
	}
	return loan, nil
}

func (store *Store) FindAllLoans() ([]*model.Loan, error) {
	loans := []*model.Loan{}
	if err := store.db.Find(&loans).Error; err != nil {
		return nil, err
	}
	return loans, nil
}

func (store *Store) FindLoansByBookID(bookID int) ([]*model.Loan, error) {
	loans := []*model.Loan{}
	err := store.db.Find(&loans, "book_id = $1", bookID).Error
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (store *Store) FindLoansByUserID(userID int) ([]*model.Loan, error) {
	loans := []*model.Loan{}
	err := store.db.Find(&loans, "user_id = $1", userID).Error
	if err != nil {
		return nil, err
	}
	return loans, nil
}

func (store *Store) CreateUser(user *model.User) error {
	return store.db.Create(&user).Error

}

func (store *Store) FindUserByID(id int) (*model.User, error) {
	user := &model.User{}
	if err := store.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (store *Store) Transaction() i.Store {
	tx := store.db.Begin()
	return NewStore(tx)
}

func (store *Store) Commit() {
	store.db.Commit()
}

func (store *Store) Rollback() {
	store.db.Rollback()
}
