package interactor

import (
	"fmt"
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
	openlibrarydto "github.com/sfqi/library/openlibrary/dto"
	"regexp"
	"strconv"
	"strings"
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

var yearRgx = regexp.MustCompile(`[0-9]{4}`)

func NewBook(store store, olc openlibraryClient) *Book {
	return &Book{
		store:   store,
		openlib: olc,
	}
}

func (b *Book) FindAll() ([]*model.Book, error) {
	return b.store.FindAllBooks()
}

func (b *Book) Create(bookRequest dto.CreateBookRequest) (*model.Book, error) {
	openLibraryBook, err := b.openlib.FetchBook(bookRequest.ISBN)
	if err != nil {
		return nil, err
	}

	book, err := toBook(openLibraryBook)
	if err != nil {
		return nil, err
	}

	if err := b.store.CreateBook(book); err != nil {
		return nil, err
	}
	return book, nil
}

func (b *Book) Update(*model.Book) error {

	return nil
}

func (b *Book) FindById(id int) (*model.Book, error) {

	return b.store.FindBookById(id)
}

func (b *Book) Delete(book *model.Book) error {
	return b.store.DeleteBook(book)

}

func toBook(book *openlibrarydto.Book) (*model.Book, error) {
	isbn10 := ""
	if book.Identifier.ISBN10 != nil {
		isbn10 = book.Identifier.ISBN10[0]
	}
	isbn13 := ""
	if book.Identifier.ISBN13 != nil {
		isbn13 = book.Identifier.ISBN13[0]
	}

	CoverId := ""
	if book.Cover.Url != "" {
		part1 := strings.Split(book.Cover.Url, "/")[5]
		part2 := strings.Split(part1, ".")[0]
		CoverId = strings.Split(part2, "-")[0]
	}
	libraryId := ""
	if book.Identifier.Openlibrary != nil {
		libraryId = book.Identifier.Openlibrary[0]
	}
	author := ""
	if book.Author != nil {
		author = book.Author[0].Name
	}

	yearString := yearRgx.FindString(book.Year)
	if yearString == "" {
		return nil, fmt.Errorf("no conversion is aloved")
	}
	year, err := strconv.Atoi(yearString)
	if err != nil {
		return nil, err
	}

	return &model.Book{
		Title:         book.Title,
		Author:        author,
		Isbn:          isbn10,
		Isbn13:        isbn13,
		OpenLibraryId: libraryId,
		CoverId:       CoverId,
		Year:          year,
	}, nil

}
