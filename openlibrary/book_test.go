package openlibrary

import (
	"net/http"
	"net/http/httptest"
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
			t.Errorf("We got error: %s", err)
		}

		if responseBook.Title != "Concrete mathematics" {
			t.Errorf("We did not get the expected response,expected %s, but got %s", expected, responseBook.Title)
		}
	})

	t.Run("book with error in title", func(t *testing.T) {
		expected := `{"ISBN:0140447938": {"title": "War and Peace (Penguin Classics)"}}`

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:0140447938": {"title": "War and Peace (Penguin Classicssss)"}}`))
		}))

		defer server.Close()
		client := &Client{
			server.URL,
		}

		responseBook, err := client.FetchBook("0140447938")
		if err != nil {
			t.Errorf("We got error: %s", err)
		}

		if responseBook.Title != "War and Peace (Penguin Classics)" {
			t.Errorf("We did not get the expected response,expected %s, but got %s", expected, responseBook.Title)
		}

	})
}
