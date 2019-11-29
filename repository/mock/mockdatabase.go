package mock

import "github.com/sfqi/library/domain/model"

type Store struct{
	Book *model.Book
	Books []*model.Book
	Err error
}

func (s *Store)FindBookByid(int)(*model.Book,error){
	return s.Book,s.Err
}

func (s *Store)CreateBook(book *model.Book) error{
	return s.Err
}

func (s *Store)UpdateBook(book *model.Book)error{
	return s.Err
}

func (s *Store)FindAllBooks()([]*model.Book,error){
	return s.Books, s.Err
}

func(s *Store)DeleteBooks(b *model.Book)error{
	return s.Err
}
