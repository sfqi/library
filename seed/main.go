package main

import (
	"github.com/jinzhu/gorm"
	"time"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/repository/postgres"
)

var books = []*model.Book{
	{
		Title:         "Information systems literacy.",
		Author:        "Hossein Bidgoli",
		Isbn:          "0023095334",
		Isbn13:        "9780023095337",
		OpenLibraryId: "OL1733511M",
		Year:          1993,
		Publisher:     "Maxwell Macmillan International",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
	{
		Title:         "Programming the World Wide Web",
		Author:        "Robert W. Sebesta",
		Isbn:          "0321303326",
		Isbn13:        "9780321303325",
		OpenLibraryId: "OL3393672M",
		Year:          2005,
		Publisher:     "Pearson/Addison-Wesley",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	},
}

func main() {

	config := postgres.PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "postgres",
		Name:     "library",
	}

	store, err := postgres.Open(config)
	if err != nil {
		panic(err)
	}
	defer store.Close()

	err = CreateBooks(store.GetDB(), books)
	if err != nil {
		panic(err)
	}

}

func CreateBooks(db *gorm.DB, books []*model.Book) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	for _, book := range books {
		if err := tx.Create(book).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func CreateLoans(db *gorm.DB, loans []*model.Loan) error {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}
	for _, loan := range loans {
		if err := tx.Create(loan).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
