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

type openlibraryClient interface {
	FetchBook(isbn string) (*openlibrarydto.Book, error)
}

type Book struct {
	store   store
	openlib openlibraryClient
	books   []model.Book
}

func NewBook(store store, olc openlibraryClient) *Book {
	return &Book{
		store:   store,
		openlib: olc,
	}
}

func (b *Book) FindAll() ([]*model.Book, error) {
	//TODO : implement function
	return nil, nil
}

func (b *Book) Create(*model.Book) error {

	return nil
}

func (b *Book) Update(*model.Book) error {

	return nil
}

func (b *Book) FindById(id int) (*model.Book, error) {

	return b.store.FindBookById(id)
	//2 scenarija, jedan kad se vrati dobro, drugi kad vrati gresku
}

func (b *Book) Delete(*model.Book) error {

	return nil
}
