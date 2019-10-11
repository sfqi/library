package main

import (
	"encoding/json"
	"net/http"
)

type BookModel struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Isbn          string `json:"isbn_10"`
	Isbn13        string `json:"isbn_13"`
	OpenLibraryId string
	CoverId       string `json:"cover"`
	Year          string `json:"publish_date"`
}

func GetBooks(w http.ResponseWriter, r *http.Request) {

	books := []BookModel{
		BookModel{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          "2019",
		},
		BookModel{
			Id:            2,
			Title:         "other title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          "2019",
		},
		BookModel{
			Id:            3,
			Title:         "another title",
			Author:        "another author",
			Isbn:          "another isbn",
			Isbn13:        "another isbon13",
			OpenLibraryId: "another some id",
			CoverId:       "another cover ID",
			Year:          "2019",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func main() {

	http.HandleFunc("/books", GetBooks)
	http.ListenAndServe(":8080", nil)

}
