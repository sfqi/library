package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/library/openlibrary"
)

type BookModel struct {
	Id            int    `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	Isbn          string `json:"isbn_10"`
	Isbn13        string `json:"isbn_13"`
	OpenLibraryId string `json:"olid"`
	CoverId       string `json:"cover"`
	Year          string `json:"publish_date"`
}

var client openlibrary.Client


type db struct{
	id int
	books []BookModel
}
func(bm *db)FindBookById(id int)(book BookModel, location int, found bool){
	for i, b := range bm.books {
		if b.Id == id {
			book = b
			location = i
			found = true
			break
		}
	}

	return book,location,found
}

var books = []BookModel{
	{
		Id:            1,
		Title:         "some title",
		Author:        "some author",
		Isbn:          "some isbn",
		Isbn13:        "some isbon13",
		OpenLibraryId: "again some id",
		CoverId:       "some cover ID",
		Year:          "2019",
	},
	{
		Id:            2,
		Title:         "other title",
		Author:        "other author",
		Isbn:          "other isbn",
		Isbn13:        "other isbon13",
		OpenLibraryId: "other some id",
		CoverId:       "other cover ID",
		Year:          "2019",
	},
	{
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

var shelf = &db{
	id: len(books), // it should be 0, but we added 3 initial books with id 1,2,3 .. so next book we add will have id=4
	books: books,
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(shelf.books)
	if err != nil {
		fmt.Println("error while getting books: ", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
}

type createBookRequest struct {
	ISBN string `json:"ISBN"`
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var createBook createBookRequest
	if err := json.NewDecoder(r.Body).Decode(&createBook); err != nil {
		errorDecodingBook(w, err)
		return
	}
	fmt.Println(createBook.ISBN)
	book, err := client.FetchBook(createBook.ISBN)
	if err != nil {
		fmt.Println("error while fetching book: ", err)
		http.Error(w, "Error while fetching book: "+err.Error(), http.StatusInternalServerError)
		return
	}

	bookToAdd := CreateBookModelFromBook(*book)
	fmt.Println("PRINTING BOOK MODEL")
	fmt.Println(bookToAdd)
	shelf.books = append(shelf.books,bookToAdd)

	if err := json.NewEncoder(w).Encode(book); err != nil {
		errorEncoding(w, err)
		return
	}
}
//FetchBook returns book, with whole attributes, but in BookModel we need some parts
//prim: we need to render Cover, to get CoverId
func CreateBookModelFromBook(b openlibrary.Book)(bm BookModel){
	//every time we createBook, Id should increase
	shelf.id = shelf.id + 1
	//Here we must hardcode this, because not everybook has both of ISBNs, so in that case we will face segmentation fault
	//This is one way to do it manualy
	isbn10 := ""
	if b.Identifier.ISBN10 != nil{
		isbn10 = b.Identifier.ISBN10[0]
	}
	isbn13 := ""
	if b.Identifier.ISBN13 != nil{
		isbn13 = b.Identifier.ISBN13[0]
	}

////////////    	Rendering CoverId
	part1 := strings.Split(b.Cover.Url,"/")[5]
	part2 := strings.Split(part1,".")[0]
	CoverId := strings.Split(part2,"-")[0]
////////////
	libraryId := ""
	if b.Identifier.Openlibrary != nil {
		libraryId = b.Identifier.Openlibrary[0]
	}

	bookToAdd :=BookModel{
		Id:            shelf.id,
		Title:         b.Title,
		Author:        b.Author[0].Name,
		Isbn:          isbn10,
		Isbn13:        isbn13,
		OpenLibraryId: libraryId,
		CoverId:       CoverId,
		Year:          b.Year,
	}
	return bookToAdd
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book BookModel
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		errorDecodingBook(w, err)
		return
	}

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		errorConvertingId(w, err)
		return
	}
	bookWithId,location, found := shelf.FindBookById(id)
	fmt.Println(bookWithId.Id ,bookWithId.Title,bookWithId.Author)
	if !found {
		errorFindingBook(w, err)
		return
	}

	shelf.books[location]=book
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		errorEncoding(w, err)
		return
	}
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		errorConvertingId(w, err)
		return
	}
	book,_,found := shelf.FindBookById(id)
	if !found {
		errorFindingBook(w, err)
		return
	}
	err = json.NewEncoder(w).Encode(book)
	if err != nil {
		errorEncoding(w, err)
		return
	}
}

func main() {
	// Setting env var
	openLibraryURL := os.Getenv("LIBRARY")
	client = *openlibrary.NewClient(openLibraryURL)

	r := mux.NewRouter()

	r.HandleFunc("/books", GetBooks).Methods("GET")
	r.HandleFunc("/books", CreateBook).Methods("POST")
	r.HandleFunc("/books/{id}", UpdateBook).Methods("PUT")
	r.HandleFunc("/book/{id}", GetBook).Methods("GET")
	http.ListenAndServe(":8080", r)
}

// Handling errors ***************
func errorDecodingBook(w http.ResponseWriter, err error) {
	fmt.Println("error while decoding book from response body: ", err)
	http.Error(w, "Error while decoding from request body", http.StatusBadRequest)
}

func errorEncoding(w http.ResponseWriter, err error) {
	fmt.Println("error while encoding book: ", err)
	http.Error(w, "Internal server error:"+err.Error(), http.StatusInternalServerError)
}

func errorConvertingId(w http.ResponseWriter, err error) {
	fmt.Println("Error while converting Id to integer ", err)
	http.Error(w, "Error while converting url parameter into integer", http.StatusBadRequest)
}

func errorFindingBook(w http.ResponseWriter, err error) {
	fmt.Println("Cannot find book with given Id ")
	http.Error(w, "Book with given Id can not be found", http.StatusBadRequest)
}
