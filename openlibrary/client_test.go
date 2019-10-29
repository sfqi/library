package openlibrary

import (
	"github.com/library/handler/dto"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/library/handler/dto"
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

	t.Run("error decoding from fetchBook", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`aa{"ISBN:0140447938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))

		defer server.Close()
		client := &Client{
			server.URL,
		}

		responseBook, err := client.FetchBook("0140447938")

		checkError(t, err, "error while decoding from FetchBook:")

		checkIfBookIsNil(responseBook, t)

	})

	t.Run("server response doesn't contain the expected ISBN key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:014044723123938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))
		defer server.Close()
		client := &Client{
			server.URL,
		}
		responseBook, err := client.FetchBook("0140447932111xxxx")

		checkError(t, err, "value for given key cannot be found:")

		checkIfBookIsNil(responseBook, t)
	})
}

func checkError(t *testing.T, err error, expected string) {
	if err == nil {
		t.Error("Expected error ,  got nil")
	}
	if !strings.Contains(err.Error(), expected) {
		t.Errorf("Got error: %s", err)
	}
}

func checkIfBookIsNil(b *dto.Book, t *testing.T) {
	if b != nil {
		t.Errorf("Expected Book to be nil, got : %v", b)
	}
}
