package middleware

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"net/http"

	"github.com/sfqi/library/domain/model"
	"github.com/sfqi/library/repository/inmemory"
)

type BookLoader struct {
	Db *inmemory.DB
}

type ContextBody struct {
	Id   int
	Book model.Book
}

func (bl BookLoader) GetBook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(strings.Split(r.URL.Path, "/")[2])
		if err != nil {
			errorConvertingId(w, err)
		}
		book, err := bl.Db.FindBookByID(id)
		if err != nil {
			errorFindingBook(w, err)
		}
		contextBody := ContextBody{
			Id:   id,
			Book: *book,
		}

		ctx := context.WithValue(r.Context(), "context", contextBody)
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
