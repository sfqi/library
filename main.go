package main

import (
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
)

type Book struct {
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

func fetchBooks(isbn string) (*Books, error) {
	url := basePath + isbn + queryParams
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	result := make(map[string]*json.RawMessage, 0)
	err = json.NewDecoder(response.Body).Decode(&result)
	var book Books
	key := fmt.Sprintf("ISBN:%v",isbn) // our structure's root is "ISBN:xxxxx", we must concatenate in one key
	fmt.Println("Raw message")

	book = parseBook(result[key])

	return &book, nil
}

type Books struct {
	Title string `json:"title"`
	Identf identifier `json:"identifiers"`
	Authr []Author `json:"authors"`
	Cvr Cover `json:"cover"`
	Year string `json:"publish_date"`
}

type identifier struct {
	ISBN10 []string `json:"isbn_10"`
	ISBN13 []string `json:"isbn_13"`
	Openlibrary []string `json:"openlibrary"`
}
type Author struct {
	Url string `json:"url"`
	Name string `json:"string"`
}

type Cover struct {
	Small string `json:"small"`
}

func parseBook(values *json.RawMessage) Books {
	var book Books
	json.Unmarshal(*values,&book)

	fmt.Println(book.Title)
	fmt.Println(book.Authr[0].Name)
	fmt.Println(book.Authr[0].Url)
	fmt.Println(book.Identf.ISBN10[0])
	fmt.Println(book.Identf.ISBN13[0])
	fmt.Println(book.Identf.Openlibrary[0])
	fmt.Println(book.Year)

	return book
}
// I created struct Books, where it will be stored whole concept of Book, so if we need Cover_Id and
// library id, we will just parse fields from Books. exp: cover_id = strings.Split(Books.Cover,"/")[5]
func main() {
	book, err := fetchBooks("9780261102736")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	//printMap(*book)
	fmt.Println(*book)

}
