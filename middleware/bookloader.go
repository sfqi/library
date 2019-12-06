package middleware

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"

	"net/http"
)

type store interface {
	FindBookById(int) (*model.Book, error)
}

type BookLoader struct {
	Db store
}

func (bl BookLoader) GetBook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			errorConvertingId(w, err)
			return
		}
		book, err := bl.Db.FindBookById(id)
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
	fmt.Println("Cannot find book with given Id ", err)
	http.Error(w, "Book with given Id can not be found", http.StatusBadRequest)
}
