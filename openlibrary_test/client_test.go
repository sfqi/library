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
	assert := assert.New(t)
	require := require.New(t)
	t.Run("book is successfully fetched", func(t *testing.T) {

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:0201558025": {"title": "Concrete mathematics"}}`))
		}))
		defer server.Close()
		expectedBook := &dto.Book{Title: "Concrete mathematics"}

		client := openlibrary.NewClient(server.URL)

		responseBook, err := client.FetchBook("0201558025")

		require.NoError(err)

		require.Equal(expectedBook, responseBook)
	})

	t.Run("error decoding from fetchBook", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`aa{"ISBN:0140447938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))

		defer server.Close()

		client := openlibrary.NewClient(server.URL)
		responseBook, err := client.FetchBook("0140447938")
		require.Error(err)

		assert.Nil(responseBook)

	})

	t.Run("server response doesn't contain the expected ISBN key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:014044723123938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))
		defer server.Close()

		client := openlibrary.NewClient(server.URL)
		responseBook, err := client.FetchBook("0140447932111xxxx")
		require.Error(err)

		assert.Nil(responseBook)
	})
}
