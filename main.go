package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type Book2 struct {
	Title         string `json:"title"`
	AuthorId      string `json:"author_id"`
	Isbn          string `json:"isbn_10"`
	Isbn13        string `json:"isbn_13"`
	OpenLibraryId string
	CoverId       string `json:"cover"`
	Year          string `json:"publish_date"`
}

const basePath = "https://openlibrary.org/api/books?bibkeys=ISBN:"
const queryParams = "&format=json&jscmd=data"

func fetchBooks(isbn string) (*Book, error) {
	url := basePath + isbn + queryParams
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	status := response.StatusCode
	fmt.Println("Status code: ",status)
	defer response.Body.Close()
	result := make(map[string]*json.RawMessage, 0)
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("ISBN:%v",isbn) // our structure's root is "ISBN:xxxxx", we must concatenate in one key
	rawBook, ok := result[key]
	if !ok {
		return nil,errors.New("Value for given key cannot be found")
	}

	book := parseBook(rawBook)
	fmt.Println(book.Identifier.ISBN10[0])
	fmt.Println(book.Identifier.ISBN13)
	return book, nil
}

type Book struct {
	Title string `json:"title"`
	Identifier identifier `json:"identifiers"`
	Author []Author `json:"authors"`
	Cover Cover `json:"cover"`
	Year string `json:"publish_date"`
}

type identifier struct {
	ISBN10 []string `json:"isbn_10"`
	ISBN13 []string `json:"isbn_13"`
	Openlibrary []string `json:"openlibrary"`
}
type Author struct {
	Name string `json:"string"`
}

type Cover struct {
	Url string `json:"small"`
}

func parseBook(values *json.RawMessage) *Book {
	var book Book
	json.Unmarshal(*values,&book)
	return &book
}
// I created struct Books, where it will be stored whole concept of Book, so if we need Cover_Id and
// library id, we will just parse fields from Books. exp: cover_id = strings.Split(Books.Cover,"/")[5]
func main() {
	book, err := fetchBooks("9780261102736")
	if err != nil {
		panic(err)
	}
	//printMap(*book)
	fmt.Println(*book)

}
