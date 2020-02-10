package main

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"

	"github.com/google/uuid"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/repository/postgres"
)

var books = []*model.Book{
	{
		Id:            1,
		Title:         "Information systems literacy.",
		Author:        "Hossein Bidgoli",
		Isbn:          "0023095334",
		Isbn13:        "9780023095337",
		OpenLibraryId: "OL1733511M",
		Year:          1993,
		Publisher:     "Maxwell Macmillan International",
	},
	{
		Id:            2,
		Title:         "Programming the World Wide Web",
		Author:        "Robert W. Sebesta",
		Isbn:          "0321303326",
		Isbn13:        "9780321303325",
		OpenLibraryId: "OL3393672M",
		Year:          2005,
		Publisher:     "Pearson/Addison-Wesley",
	},
}

var loans = []*model.Loan{
	{
		TransactionID: uuid.New().String(),
		UserID:        5,
		BookID:        1,
		Type:          1,
	},
	{
		TransactionID: uuid.New().String(),
		UserID:        10,
		BookID:        2,
		Type:          0,
	},
}

func areTablesEmpty(db *gorm.DB) error {
	var bookCount int
	var loanCount int
	var count int

	switch count {
	case bookCount:
		if err := db.Model(&model.Book{}).Count(&bookCount).Error; err != nil {
			return err
		}
	case loanCount:
		if err := db.Model(&model.Loan{}).Count(&loanCount).Error; err != nil {
			return err
		}
	}

	switch count {
	case bookCount:
		return errors.New("Table book is not empty!")
	case loanCount:
		return errors.New("Table loan is not empty!")
	}
	return nil
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

	tx := store.DB().Begin()

	defer func() {
		store.Close()
		if r := recover(); r != nil {
			fmt.Println(r.(error))
			tx.Rollback()
		}
	}()

	if err = areTablesEmpty(store.DB()); err != nil {
		panic(err)
	}

	for _, book := range books {
		if err := tx.Create(book).Error; err != nil {
			panic(err)
		}
	}

	if err = areTablesEmpty(store.DB()); err != nil {
		panic(err)
	}

	for _, loan := range loans {
		if err := tx.Create(loan).Error; err != nil {
			panic(err)
		}
	}

	tx.Commit()

}
