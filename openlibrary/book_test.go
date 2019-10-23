package openlibrary

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchBook(t *testing.T) {
	t.Run("book is successfully fetched", func(t *testing.T) {
		expected := `{ISBN:0201558025": {"title": "Concrete mathematics"}}`
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:0201558025": {"title": "Concrete mathematics"}}`))
		}))

		defer server.Close()
		client := &Client{
			server.URL,
		}
		responseBook, err := client.FetchBook("0201558025")

		if err != nil {
			t.Errorf("Got error: %s", err)
		}

		if responseBook.Title != "Concrete mathematics" {
			t.Errorf("We did not get the expected response,expected %s, got %s", expected, responseBook.Title)
		}
	})

	t.Run("book with error in title", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:0140447938": {"title": "Warr and Peace (Penguin Classics)"}}`))
		}))
		responseTitle := "Warr and Peace (Penguin Classics)" // i added one more 'r' in war
		defer server.Close()
		client := &Client{
			server.URL,
		}

		responseBook, err := client.FetchBook("0140447938")
		if err != nil{
			if responseBook.Title == responseTitle{
				t.Errorf("Title of the book: %s, should not match with response %s",responseBook.Title,responseTitle)
			}
		}
	})

	t.Run("error decoding from fetchBook",func(t *testing.T){
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`aa{"ISBN:0140447938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))

		defer server.Close()
		client := &Client{
			server.URL,
		}

		responseBook, err := client.FetchBook("0140447938")
		chechIfErrorIsNIl(err,t)
		if !strings.Contains(err.Error(),"error while decoding from FetchBook:"){
			t.Errorf("Got error: %s", err)
		}

		checkIfBookIsNil(responseBook,t)

	})

	t.Run("server response doesn't contain the expected ISBN key", func(t *testing.T){
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:014044723123938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))
		//0140447938
		defer server.Close()
		client := &Client{
			server.URL,
		}
		//whatever we pass here as invalid isbn, we will get the error we want
		responseBook , err := client.FetchBook("0140447932111xxxx")
		chechIfErrorIsNIl(err,t)
		if err.Error() != "value for given key cannot be found: ISBN:0140447932111xxxx" {
			t.Errorf("Got error : %s",err)
		}
		checkIfBookIsNil(responseBook,t)
	})
}

func checkIfBookIsNil(b *Book,t *testing.T){
	if b != nil{
		t.Errorf("Expected Book to be nil, got : %v",b)
	}
}

func chechIfErrorIsNIl(err error,t *testing.T){
	if err == nil {
		t.Error("Expected error ,  got nil")
	}
}
