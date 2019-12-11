package openlibrary_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/sfqi/library/openlibrary/dto"
)

type ClientMock struct{
	mock.Mock
	Book *dto.Book
	Err  error
}

func(cm *ClientMock)FetchBook(isbn string) (*dto.Book, error){
	fmt.Println("Input isbn is: " + isbn)
	args := cm.Called(isbn)
	err := args.Error(1)
	if err != nil{
		return nil, err
	}
	// Cim imamo error, znamo da je knjiga nil, sa ovim if-om izbegavamo
	// da program pukne zbog konverzije  -- Izbrisacu
	book := args.Get(0).(dto.Book)
	return &book,args.Error(1)
}

func TestFetchBook(t *testing.T) {
	t.Run("book is successfully fetched", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:0201558025": {"title": "Concrete mathematics"}}`))
		}))
		defer server.Close()

		mockClient := ClientMock{}
		mockClient.On("FetchBook","0201558025").Return(dto.Book{
			Title:      "Concrete mathematics",
			Year:       2019,
		},nil)
		mockClient.FetchBook("0201558025")
		mockClient.AssertExpectations(t)
	})

	t.Run("error decoding from fetchBook", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`aa{"ISBN:0140447938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))

		defer server.Close()

		mockClient := ClientMock{}
		mockClient.On("FetchBook","aaa0201558025").Return(nil,errors.New("error while decoding from FetchBook:"))
		mockClient.FetchBook("aaa0201558025")
		mockClient.AssertExpectations(t)
	})

	t.Run("server response doesn't contain the expected ISBN key", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"ISBN:014044723123938": {"title": "War and Peace (Penguin Classics)"}}`))
		}))
		defer server.Close()

		mockClient := ClientMock{}
		mockClient.On("FetchBook","02015580251111").Return(nil,errors.New("value for given key cannot be found:"))
		mockClient.FetchBook("02015580251111")
		mockClient.AssertExpectations(t)
	})
}

