package interactor

import (
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
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

func (b *Book) Update(book *model.Book, updateBookRequest dto.UpdateBookRequest) (*model.Book, error) {
	book.Title = updateBookRequest.Title
	book.Year = updateBookRequest.Year
	if err := b.store.UpdateBook(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (b *Book) FindById(int) (*model.Book, error) {

	return nil, nil
}

func (b *Book) Delete(*model.Book) error {

	return nil
}
