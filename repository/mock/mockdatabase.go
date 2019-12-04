package mock

import (
	"github.com/sfqi/library/domain/model"
)

type Store struct{
	Books []*model.Book
	Err error
}

func NewStore(books []*model.Book,err error)*Store{
	return &Store{
		Books: books,
		Err:err,
	}
}

func (s *Store)FindBookById(id int) (*model.Book, error){
	return s.findBookByID(id),s.Err
}

func (s *Store) findBookByID(id int) (*model.Book) {
	for _, b := range s.Books {
		if b.Id == id {
			return b
		}
	}
	return nil
}

func (s *Store)CreateBook(book *model.Book) error{
	return s.Err
}

func (s *Store)UpdateBook(book *model.Book) error{
	return s.Err
}

func (s *Store)FindAllBooks() ([]*model.Book, error){
	return s.Books, s.Err
}

func(s *Store)DeleteBook(book *model.Book)error{
	return s.Err
}
