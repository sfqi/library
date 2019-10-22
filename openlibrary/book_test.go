package openlibrary

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
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
			t.Errorf("We got error: %s", err)
		}

		if responseBook.Title != "Concrete mathematics" {
			t.Errorf("We did not get the expected response,expected %s, but got %s", expected, responseBook.Title)
		}
	})

	t.Run("book with error in title", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:0140447938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))

		defer server.Close()
		client := &Client{
			server.URL,
		}

		responseBook, err := client.FetchBook("0201558025")
		if err == nil {
			t.Error("Expected error , but got nil")
		}
		if err.Error() != "value for given key cannot be found: ISBN:0201558025" {
			t.Errorf("We got error: %s", err)
		}

		if responseBook != nil {
			t.Errorf("Expected to be nil but got: %v", responseBook)
		}
	})

	t.Run("book error decoding from fetchBook",func(t *testing.T){
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`aa{"ISBN:0140447938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))

		defer server.Close()
		client := &Client{
			server.URL,
		}

		responseBook, err := client.FetchBook("0140447938")
		if err == nil {
			t.Error("Expected error , but got nil")
		}
		if !strings.Contains(err.Error(),"error while decoding from FetchBook:"){
			t.Errorf("We got error: %s", err)
		}

		if responseBook != nil {
			t.Errorf("Expected to be nil but got: %v", responseBook)
		}

	})

	t.Run("Book with given ISBN cannot be found", func(t *testing.T){
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:0140447938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))

		defer server.Close()
		client := &Client{
			server.URL,
		}

		responseBook , err := client.FetchBook("014044793")
		if err == nil {
			t.Error("Expected error , but got nil")
		}
		if err.Error() != "value for given key cannot be found: ISBN:014044793" {
			t.Errorf("Here we get error : %s",err)
		}
		if responseBook != nil {
			t.Errorf("We expected Book to be nil, but got : %v",responseBook)
		}
	})
}
