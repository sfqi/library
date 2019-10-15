package fetcher

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const basePath = "https://openlibrary.org/api/books?bibkeys=ISBN:"
const queryParams = "&format=json&jscmd=data"

type Book struct {
	Title      string     `json:"title"`
	Identifier identifier `json:"identifiers"`
	Author     []Author   `json:"authors"`
	Cover      Cover      `json:"cover"`
	Year       string     `json:"publish_date"`
}

type identifier struct {
	ISBN10      []string `json:"isbn_10"`
	ISBN13      []string `json:"isbn_13"`
	Openlibrary []string `json:"openlibrary"`
}
type Author struct {
	Name string `json:"string"`
}

type Cover struct {
	Url string `json:"small"`
}

func FetchBook(isbn string) (*Book, error) {
	url := basePath + isbn + queryParams
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	status := response.StatusCode
	fmt.Println("Status code: ", status)
	defer response.Body.Close()
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
	sliceOfBytes, err := rawBook.MarshalJSON()
	if err != nil {
		fmt.Println("Error converting *rawMessage into []byte")
	}

	err = json.Unmarshal(sliceOfBytes, &book)
	fmt.Println(book.Identifier.ISBN10[0])
	fmt.Println(book.Identifier.ISBN13)
	return &book, nil
}
