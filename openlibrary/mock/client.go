package mock

import (
	"github.com/sfqi/library/openlibrary/dto"
	"github.com/stretchr/testify/mock"
)

type Client struct {
	mock.Mock
}

func (c *Client) FetchBook(isbn string) (*dto.Book, error) {
	args := c.Called(isbn)
	if args.Get(0) != nil {
		return args.Get(0).(*dto.Book), nil
	}
	return nil, args.Error(1)
}
