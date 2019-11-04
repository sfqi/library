package openlibrary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	openlibrary "github.com/library/openlibrary/dto"
)

const bookPath = "/books"
const bibkeys = "?bibkeys=ISBN:%s"
const formatParams = "&format=json&jscmd=data"

var fetchBookPath = bookPath + bibkeys + formatParams

type Client struct {
	url string
}

func NewClient(url string) *Client {
	url = strings.TrimSuffix(url, "/")
	return &Client{
		url: url,
	}

}

func (client *Client) FetchBook(isbn string) (*openlibrary.Book, error) {
	url := fmt.Sprintf(client.url+fetchBookPath, isbn)

	response, err := http.Get(url)
	if err != nil {
		err := fmt.Errorf("error while getting url: %v", err)
		return nil, err
	}
	defer response.Body.Close()
	status := response.StatusCode
	fmt.Println("Status code: ", status)
	result := make(map[string]*json.RawMessage, 0)
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		err := fmt.Errorf("error while decoding from FetchBook: %v", err)
		return nil, err
	}

	key := "ISBN:" + isbn
	rawBook, ok := result[key]
	if !ok {
		errorKey := fmt.Errorf("value for given key cannot be found: %s", key)
		return nil, errorKey
	}

	var book openlibrary.Book

	err = json.Unmarshal(*rawBook, &book)
	if err != nil {
		err := fmt.Errorf("error while Unmarshaling from FetchBook: %v", err)
		return nil, err
	}
	return &book, nil
}
