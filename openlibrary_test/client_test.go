package openlibrary_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sfqi/library/openlibrary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sfqi/library/openlibrary/dto"
)

func TestFetchBook(t *testing.T) {
	t.Run("book is successfully fetched", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:0201558025": {"title": "Concrete mathematics"}}`))
		}))
		expectedBook := &dto.Book{Title: "Concrete mathematics"}
		defer server.Close()

		client := openlibrary.NewClient(server.URL)

		responseBook, err := client.FetchBook("0201558025")

		if err != nil {
			require.Error(t, err, "Got error")
		}

		require.Equal(t, expectedBook, responseBook, "Response body differs")

		//assert.Equal(t, responseBook, server.Client, "Response body differs")
	})

	t.Run("error decoding from fetchBook", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`aa{"ISBN:0140447938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))

		defer server.Close()

		client := openlibrary.NewClient(server.URL)
		responseBook, err := client.FetchBook("0140447938")

		assert.Contains(t, err, "error while decoding from FetchBook:")
		assert.Nil(t, responseBook)

	})

	t.Run("server response doesn't contain the expected ISBN key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:014044723123938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))
		defer server.Close()

		client := openlibrary.NewClient(server.URL)
		responseBook, err := client.FetchBook("0140447932111xxxx")

		assert.Contains(t, err, "value for given key cannot be found:")
		assert.Nil(t, responseBook)
	})
}
