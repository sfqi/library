package openlibrary

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const bookPath="/books"
const bibkeys = "?bibkeys=ISBN:%s"
const formatParams = "&format=json&jscmd=data"
var fetchBookPath = bookPath+bibkeys + formatParams
type Client struct {
	Url string
}

type Book struct {
	Title      string     `json:"title"`
	Identifier identifier `json:"identifiers"`
	Author     []author   `json:"authors"`
	Cover      cover      `json:"cover"`
	Year       string     `json:"publish_date"`
}

type identifier struct {
	ISBN10      []string `json:"isbn_10"`
	ISBN13      []string `json:"isbn_13"`
	Openlibrary []string `json:"openlibrary"`
}
type author struct {
	Name string `json:"string"`
}

type cover struct {
	Url string `json:"small"`
}

func (client *Client)FetchBook(isbn string) (*Book, error) {
	url := fmt.Sprintf(client.Url + fetchBookPath,isbn)
	fmt.Println(url)

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
		errorKey := fmt.Errorf("value for given key cannot be found: %s ", result[key])
		return nil, errorKey
	}

	var book Book

	err = json.Unmarshal(*rawBook, &book)
	if err != nil {
		err := fmt.Errorf("error while Unmarshaling from FetchBook: %v", err)
		return nil, err
	}
	return &book, nil
}
