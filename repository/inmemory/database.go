package inmemory

import (
	"fmt"

	"github.com/sfqi/library/domain/model"

	"time"
)

var timeNow = time.Now()
var earlier10sec = timeNow.Add(-10 * time.Second)
var earlier15sec = timeNow.Add(-15 * time.Second)



var books = []model.Book{
	{
		Id:            1,
		Title:         "some title",
		Author:        "some author",
		Isbn:          "some isbn",
		Isbn13:        "some isbon13",
		OpenLibraryId: "again some id",
		CoverId:       "some cover ID",
		Year:          2019,
		CreatedAt:     earlier10sec,
		UpdatedAt:     earlier10sec,
	},
	{
		Id:            2,
		Title:         "other title",
		Author:        "other author",
		Isbn:          "other isbn",
		Isbn13:        "other isbon13",
		OpenLibraryId: "other some id",
		CoverId:       "other cover ID",
		Year:          2019,
		CreatedAt:     earlier15sec,
		UpdatedAt:     earlier10sec,
	},
	{

		Id:            3,
		Title:         "another title",
		Author:        "another author",
		Isbn:          "another isbn",
		Isbn13:        "another isbon13",
		OpenLibraryId: "another some id",
		CoverId:       "another cover ID",
		Year:          2019,
		CreatedAt:     timeNow,
		UpdatedAt:     timeNow,
	},
}

type DB struct {
	Id    int
	books []model.Book
}

func NewDB() *DB {
	return &DB{
		Id:    len(books),
		books: books,
	}
}

func (db *DB) FindBookById(id int) (*model.Book, error) {
	book, _, err := db.findBookByID(id)
	return book, err
}

func (db *DB) findBookByID(id int) (*model.Book, int, error) {
	for i, b := range db.books {
		if b.Id == id {
			return &b, i, nil
		}
	}
	return nil, -1, fmt.Errorf("error while findBookByID")
}

func (db *DB) FindAllBooks() ([]*model.Book,error) {
	pointers := make([]*model.Book,len(db.books))
	for i:=0;i<len(db.books);i++{
		pointers[i] = &db.books[i]
	}
	fmt.Println(pointers)
	return pointers,nil
}

func (db *DB) CreateBook(book *model.Book) error {
	db.Id++
	now := time.Now()
	book.CreatedAt = now
	book.UpdatedAt = now

	book.Id = db.Id
	db.books = append(db.books, *book)
	fmt.Println(db.books)
	return nil
}

func (db *DB) UpdateBook(toUpdate *model.Book) error {
	book, index, err := db.findBookByID(toUpdate.Id)

	book.UpdatedAt = time.Now()
	book.Title = toUpdate.Title
	book.Year = toUpdate.Year
	toUpdate = book
	if err != nil {
		return err
	}
	db.books[index] = *book
	return nil
}

func (db *DB) DeleteBook(book *model.Book) error {
	_, loc, err := db.findBookByID(book.Id)
	if err != nil {
		return err
	}
	db.books = append(db.books[:loc], db.books[loc+1:]...)
	return nil
}
