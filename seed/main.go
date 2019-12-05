package main

import (
	"github.com/sfqi/library/domain/model"
	"fmt"
	"time"
	"github.com/sfqi/library/repository/postgres"
)

	var books = []*model.Book{
		{
			Title:         "The Go Programming Language",
			Author:        "Donovan and Kernighan",
			Isbn:          "==========",
			Isbn13:        "9780134190570",
			OpenLibraryId: "0001",
			CoverId:       "0001",
			Year:          2015,
			Publisher: 		"Addison-Wesley Professional",
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
		},
		{
			Title:         "Go Web Programming",
			Author:        "Sau Sheong Chang",
			Isbn:          "---------",
			Isbn13:        "9781617292569",
			OpenLibraryId: "0002",
			CoverId:       "0002",
			Year:          2016,
			Publisher: 		"Manning Publications",
			CreatedAt:     time.Time{},
			UpdatedAt:     time.Time{},
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