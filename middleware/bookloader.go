package middleware

import (
	"context"

	"net/http"

	"github.com/sfqi/library/repository/inmemory"
)

type BookLoader struct {
	Db *inmemory.DB
}

func (bl BookLoader) GetBook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		book, err := bl.Db.FindBookByID(1)
		if err != nil {
			return
		}
		ctx := context.WithValue(r.Context(), "book", book)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
