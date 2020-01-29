package middleware

import (
	"context"
	"fmt"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sfqi/library/domain/model"

	"net/http"
)

type userInteractor interface {
	FindByID(id int) (*model.User, error)
}

type UserLoader struct {
	Interactor userInteractor
}

func (ul UserLoader) GetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Println("Error while converting ID to integer ", err)
			http.Error(w, "Error while converting url parameter into integer", http.StatusBadRequest)
			return
		}

		user, err := ul.Interactor.FindByID(id)
		if err != nil {
			fmt.Println("Cannot find user with given ID ", err)
			http.Error(w, "User with given ID can not be found", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
