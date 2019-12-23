package mock

import (
	"github.com/sfqi/library/domain/model"
	"github.com/stretchr/testify/mock"
)

type Store struct {
	mock.Mock
}

func (s *Store) FindBookById(id int)(*model.Book, error) {
	args := s.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Book), nil
	}
	return nil, args.Error(1)
}

func (s *Store) CreateBook(book *model.Book) error {
	args := s.Called(book)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *Store) UpdateBook(book *model.Book) error {
	args := s.Called(book)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (s *Store) FindAllBooks() ([]*model.Book, error) {
	args := s.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Book), nil
	}
	return nil, args.Error(1)
}

func (s *Store) DeleteBook(book *model.Book) error {
	args := s.Called(book)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}