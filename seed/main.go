package main

import (
	"github.com/sfqi/library/domain/model"
	"fmt"
	"time"
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
			Publisher: 	   "Maxwell Macmillan International",
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
			Publisher: 	   "Pearson/Addison-Wesley",
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
	if err != nil{
		panic(err)
	}
	fmt.Println("Successfully connected")

	for _, book := range books {
		if err := store.CreateBook(book); err != nil {
			panic(err)
		}
	}

}