package mock

import (
	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/handler/dto"
	"github.com/stretchr/testify/mock"
)

type Book struct {
	mock.Mock
}

func (b *Book) FindAll() ([]*model.Book, error) {
	args := b.Called()
	if args.Get(0) != nil {
		return args.Get(0).([]*model.Book), nil
	}
	return nil, args.Error(1)
}

func (b *Book) Create(bookRequest dto.CreateBookRequest) (*model.Book, error) {
	args := b.Called(bookRequest)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Book), nil
	}
	return nil, args.Error(1)
}

func (b *Book) Update(book *model.Book, updateBookRequest dto.UpdateBookRequest) (*model.Book, error) {
	args := b.Called(book, updateBookRequest)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Book), nil
	}
	return nil, args.Error(1)
}

func (b *Book) FindById(id int) (*model.Book, error) {
	args := b.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*model.Book), nil
	}
	return nil, args.Error(1)
}

func (b *Book) Delete(book *model.Book) error {
	args := b.Called(book)
	return args.Error(0)
}
