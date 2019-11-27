package middleware

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"net/http"

	"github.com/sfqi/library/repository/inmemory"
)

type BookLoader struct {
	Db *inmemory.DB
}

func (bl BookLoader) GetBook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
		if err != nil {
			errorConvertingId(w, err)
			return
		}
		book, err := bl.Db.FindBookByID(id)
		if err != nil {
			errorFindingBook(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), "book", book)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func errorConvertingId(w http.ResponseWriter, err error) {
	fmt.Println("Error while converting Id to integer ", err)
	http.Error(w, "Error while converting url parameter into integer", http.StatusBadRequest)
}

func errorFindingBook(w http.ResponseWriter, err error) {
	fmt.Println("Cannot find book with given Id ")
	http.Error(w, "Book with given Id can not be found", http.StatusBadRequest)
}
