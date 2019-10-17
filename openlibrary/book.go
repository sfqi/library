package openlibrary

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const bibkeys = "books?bibkeys=ISBN:"
const queryParams = "&format=json&jscmd=data"

type OL struct {
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

func (ol *OL)FetchBook(isbn string) (*Book, error) {
	url := ol.Url + bibkeys+isbn+queryParams
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	status := response.StatusCode
	fmt.Println("Status code: ", status)
	result := make(map[string]*json.RawMessage, 0)
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("ISBN:%v", isbn)
	rawBook, ok := result[key]
	if !ok {
		return nil, errors.New("Value for given key cannot be found")
	}
	var book Book

	err = json.Unmarshal(*rawBook, &book)
	if err != nil {
		return nil, err
	}
	return &book, nil
}
