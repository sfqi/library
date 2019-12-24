package interactor

import (
	"github.com/sfqi/library/domain/model"
	openlibrarydto "github.com/sfqi/library/openlibrary/dto"
)

type store interface {
	FindBookById(int) (*model.Book, error)
	CreateBook(*model.Book) error
	UpdateBook(*model.Book) error
	FindAllBooks() ([]*model.Book, error)
	DeleteBook(*model.Book) error
}

type openLibraryClient interface {
	FetchBook(isbn string) (*openlibrarydto.Book, error)
}

type Book struct {
	Db  store
	Olc openLibraryClient
}

func NewBook(store store, olc openLibraryClient) *Book {
	return &Book{
		Db:  store,
		Olc: olc,
	}
}

func (bi *Book) FindAll() ([]*model.Book, error) {
	//TODO : implement function
	return nil, nil
}

func (bi *Book) Create(*model.Book) error {

	return nil
}

func (bi *Book) Update(*model.Book) error {

	return nil
}

func (bi *Book) FindById(int) (*model.Book, error) {

	return nil, nil
}

func (bi *Book) Delete(*model.Book) error {

	return nil
}
