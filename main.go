package main

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"net/http"
	"strings"
)

type Book struct {
	Title string `json:"title"`
	AuthorId string `json:"author_id"`
	Isbn string `json:"isbn_10"`
	Isbn13 string `json:"isbn_13"`
	OpenLibraryId string
	CoverId string `json:"cover"`
	Year string `json:"publish_date"`
}

const BaseApi = "https://openlibrary.org/api/books?bibkeys=ISBN:"

func GetBooksFromApi(isbn string) (*Book,error){
	url := BaseApi + isbn + "&format=json&jscmd=data"
	fmt.Println(url)
	response,err := http.Get(url)
	if err != nil{
		fmt.Println("Error getting response from url")
		return  nil,err
	}
	defer response.Body.Close()
	result := make(map[string]interface{},0)

	err = json.NewDecoder(response.Body).Decode(&result)
	var book Book
	for _,value := range result {
		 book = renderBook(value.(map[string]interface{})) // HERE WE PASS ALL DATA TO RENDERING, TO GET BOOK's ATTRIBUTES
	}
	return &book,nil
}

func renderBook(values map[string]interface{})Book{
	var book Book
	book.Title =  fmt.Sprintf("%v",values["title"])
	book.Year = fmt.Sprintf("%v",values["publish_date"])
	identifiers := values["identifiers"]
	// CASTING map[string]interface{} to map[string][]string is not allowed, i dont understand why
	// so i had to marshal date in []byte, and then unmarshall it in map[string][]string
	jsonformat , err:= json.Marshal(identifiers)
	if err != nil {
		log.Fatal(err)
	}
	var mapa map[string][]string
	err = json.Unmarshal(jsonformat,&mapa)
	if err != nil {
		fmt.Println("Error unmarhsaling to map")
	}
	// Now we have our map, we iterate thru and populate needed fields for Book

	for key,value := range mapa{
		// as we can see, we have pairs key:value, only one value for every key.. comparing to key we assign value of our fields that are matter
		// for every key we have slice of values, but we need only first value from slice
		if key == "isbn_10"{
			book.Isbn = value[0]
		}
		if key == "isbn_13"{
			book.Isbn13 = value[0]
		}
		if key == "openlibrary" {
			book.OpenLibraryId = value[0]
		}

	}

	//NOW WE MUST GET AUTHOR ID
	//Posto je 'authors' interfejs, moramo da ga marsalujemo, kako bismo posle unmarshalovali u mapu,
	//jer nam dodadtno kastovanje nije dozvoljeno
	// mapp1 je niz mapa, tako je predstavljeno  u json formatu koji dobijemo u responsu
	authors := values["authors"]
	jsonformat1 , err:= json.Marshal(authors)
	if err != nil {
		fmt.Println(err)
	}
	map1 := make([]map[string]string,0)
	er2 := json.Unmarshal(jsonformat1,&map1)
	if er2 != nil {
		fmt.Println(er2)
	}
	for key1, value1 := range map1[0] {
		if key1 == "url" {
			parts := strings.Split(value1,"/") //0: https   1:____  2: openlibrary.org   3: authors  4. author_ID
			book.AuthorId = parts[4]
		}
	}

	//And last, we need to get CoverID
	coverInterface := values["cover"]
	jsonfromat2 , err2 := json.Marshal(&coverInterface)
	if err2 != nil {
		fmt.Println("Error marshaling")
	}
	var lastMap map[string]string
	err2 = json.Unmarshal(jsonfromat2,&lastMap)
	if err2 != nil {
		fmt.Println(err2)
	}
	for key3,value3 := range lastMap {
		// Which comaparison we do , doesnt matter, because ID is the same, but i compare to first one
		if key3 == "small" {
			parts := strings.Split(value3,"/") // now we need to split last part : 2341302-S.jpg
			// We need only number 234.. so we need to split again, but on separator '-'
			book.CoverId = strings.Split(parts[len(parts)-1],"-")[0]
		}
	}

	fmt.Println("Title  : ",book.Title)
	fmt.Println("Date of publishing: ",book.Year)
	fmt.Println("ISBN13  : ",book.Isbn13)
	fmt.Println("ISBN10  : ",book.Isbn)
	fmt.Println("AuthorID  : ",book.AuthorId)
	fmt.Println("CoverID  : ",book.CoverId)


	return book
}

func main() {
	book,err := GetBooksFromApi("9780261102736")
	if err != nil {
		fmt.Errorf(err.Error())
	}
	//printMap(*book)
	fmt.Println(*book)


}
