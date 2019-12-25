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
}

var yearRgx = regexp.MustCompile(`[0-9]{4}`)

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

func (b *Book) Create(bookRequest dto.CreateBookRequest) (*model.Book, error) {
	openLibraryBook, err := b.openlib.FetchBook(bookRequest.ISBN)
	if err != nil {
		return nil, err
	}

	book := toBook(openLibraryBook)

	if err := b.store.CreateBook(book); err != nil {
		return nil, err
	}
	return book, nil
}

func (b *Book) Update(*model.Book) error {

	return nil
}

func (b *Book) FindById(int) (*model.Book, error) {

	return nil, nil
}

func (b *Book) Delete(*model.Book) error {

	return nil
}

func toBook(book *openlibrarydto.Book) (bm *model.Book) {
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

	year := 0
	var err error
	yearString := yearRgx.FindString(book.Year)
	if yearString != "" {
		year, err = strconv.Atoi(yearString)
		if err != nil {
			fmt.Println("error while converting year from string to int", err)
			return nil
		}
	}

	bookToAdd := model.Book{
		Title:         book.Title,
		Author:        author,
		Isbn:          isbn10,
		Isbn13:        isbn13,
		OpenLibraryId: libraryId,
		CoverId:       CoverId,
		Year:          year,
	}
	return &bookToAdd
}
