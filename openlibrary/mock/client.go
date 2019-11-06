package mock

import (
	"github.com/library/openlibrary/dto"
)

type Client struct {
	Book *dto.Book
	Err  error
}

func (cm *Client) FetchBook(isbn string) (*dto.Book, error) {
	return cm.Book, cm.Err
}
