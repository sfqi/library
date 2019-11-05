package olmock

import (
	"github.com/library/handler/dto"
)

type ClientMock struct{
	Book *dto.Book
	Err error
}

func(cm *ClientMock)FetchBook(isbn string)(*dto.Book, error){
	return cm.Book, cm.Err
}
