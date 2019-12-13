package mock

import (
	"github.com/sfqi/library/domain/model"
	"github.com/stretchr/testify/mock"
)

type StoreMock struct {
	mock.Mock
}

func (sm *StoreMock) FindBookById(id int) (*model.Book, error) {
	args := sm.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Book), nil
	}
	return nil, args.Error(1)
}

func (sm *StoreMock) findBookByID(id int) *model.Book {
	args := sm.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Book)
	}
	return nil
}

func (sm *StoreMock) CreateBook(book *model.Book) error {
	args := sm.Called(book)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (sm *StoreMock) UpdateBook(book *model.Book) error {
	args := sm.Called(book)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (sm *StoreMock) FindAllBooks() ([]*model.Book, error) {
	args := sm.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Book), nil
	}
	return nil, args.Error(1)
}

func (sm *StoreMock) DeleteBook(book *model.Book) error {
	args := sm.Called(book)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}
