package middleware

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"
	"github.com/sirupsen/logrus"
)

type bookInteractor interface {
	FindById(id int) (*model.Book, error)
}

type BookLoader struct {
	Interactor bookInteractor
	Logger     *logrus.Logger
}

func (bl BookLoader) GetBook(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			bl.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		book, err := bl.Interactor.FindById(id)
		if err != nil {
			bl.Logger.Error(err)
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		ctx := context.WithValue(r.Context(), "book", book)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
