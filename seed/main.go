package main

import (
	"github.com/google/uuid"
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

var users = []*model.User{
	{
		Id:        1,
		Email:     "joe@doe.com",
		Name:      "Joe",
		LastName:  "Doe",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
	{
		Id:        2,
		Email:     "jane@doe.com",
		Name:      "Jane",
		LastName:  "Doe",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
	{
		Id:        3,
		Email:     "john@smith.com",
		Name:      "John",
		LastName:  "Smith",
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
	},
}

var loans = []*model.Loan{
	{
		ID:            1,
		TransactionID: uuid.New().String(),
		UserID:        1,
		BookID:        2,
		Type:          1,
	},
	{
		ID:            2,
		TransactionID: uuid.New().String(),
		UserID:        2,
		BookID:        1,
		Type:          0,
	},
	{
		ID:            3,
		TransactionID: uuid.New().String(),
		UserID:        1,
		BookID:        1,
		Type:          1,
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
	for _, book := range books {
		if err := store.CreateBook(book); err != nil {
			panic(err)
		}
	}
	for _, loan := range loans {
		if err := store.CreateLoan(loan); err != nil {
			panic(err)
		}
	}

	for _, user := range users {
		if err := store.CreateUser(user); err != nil {
			panic(err)
		}
	}

}
