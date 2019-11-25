package postgres

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/sfqi/library/domain/model"
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
func (store *Store) FindById(id int) (*model.Book, error) {
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

func(store *Store)DeleteBook(book *model.Book)error{
	err := store.db.Where("id = ?", book.Id).Delete(&model.Book{}).Error
	if err != nil{
		return err
	}
	return nil
}
