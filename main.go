package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/pkg/errors"
)

const basePath = "https://openlibrary.org/api/books?bibkeys=ISBN:"
const queryParams = "&format=json&jscmd=data"

type BookModel struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Isbn          string `json:"isbn_10"`
	Isbn13        string `json:"isbn_13"`
	OpenLibraryId string
	CoverId       string `json:"cover"`
	Year          string `json:"publish_date"`
}

func GetBooks(w http.ResponseWriter, r *http.Request) {

	books := []BookModel{
		BookModel{
			Id:            1,
			Title:         "some title",
			Author:        "some author",
			Isbn:          "some isbn",
			Isbn13:        "some isbon13",
			OpenLibraryId: "again some id",
			CoverId:       "some cover ID",
			Year:          "2019",
		},
		BookModel{
			Id:            2,
			Title:         "other title",
			Author:        "other author",
			Isbn:          "other isbn",
			Isbn13:        "other isbon13",
			OpenLibraryId: "other some id",
			CoverId:       "other cover ID",
			Year:          "2019",
		},
		BookModel{
			Id:            3,
			Title:         "another title",
			Author:        "another author",
			Isbn:          "another isbn",
			Isbn13:        "another isbon13",
			OpenLibraryId: "another some id",
			CoverId:       "another cover ID",
			Year:          "2019",
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

type createBookRequest struct {
	ISBN string `json:"ISBN"`
}

func WriteBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var createBook createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&createBook); err != nil {
		fmt.Println(err)
		return
	}
	book, err := fetchBook(createBook.ISBN)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := json.NewEncoder(w).Encode(book); err != nil {
		fmt.Println(err)
		return
	}
}

func fetchBook(isbn string) (*Book, error) {
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

	key := fmt.Sprintf("ISBN:%v", isbn) // our structure's root is "ISBN:xxxxx", we must concatenate in one key
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

// I created struct Books, where it will be stored whole concept of Book, so if we need Cover_Id and
// library id, we will just parse fields from Books. exp: cover_id = strings.Split(Books.Cover,"/")[5]
func main() {
	r := mux.NewRouter()
	book, err := fetchBook("9780261102736")
	if err != nil {
		panic(err)
	}
	fmt.Println(*book)

	r.HandleFunc("/books", GetBooks).Methods("GET")
	r.HandleFunc("/books", WriteBooks).Methods("POST")

	http.ListenAndServe(":8080", r)
}
