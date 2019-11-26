package middleware

import (
	"fmt"
	"net/http"

	"github.com/sfqi/library/repository/inmemory"
)

type BookLoader struct {
	Db *inmemory.DB
}

func (bl BookLoader) GetBook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		book, _ := bl.Db.FindBookByID(1)

		fmt.Println(book)

		next.ServeHTTP(w, r)
	})
}
