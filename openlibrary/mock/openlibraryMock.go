package mock

import (
	"errors"
	"github.com/library/handler/dto"
	"strconv"
)

type ClientMock struct{
	url string
	books []dto.Book
}

func(cm *ClientMock)FetchBook(isbn string)(*dto.Book, error){
	if len(isbn)!=10{
		return nil, errors.New("Isbn is not valid lenght")
	}else if _,err:=strconv.Atoi(isbn);err != nil{
		return nil, err
	}else if isbn==""{
		return nil,errors.New("Book with that isbn can not be found")
	}else{
		return &dto.Book{
			Title:      "TestKnjiga",
			Identifier: nil,
			Author:     nil,
			Cover:      nil,
			Year:       "2019",
		}, nil
	}

}
