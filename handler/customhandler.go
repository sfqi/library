package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type customHandler func(http.ResponseWriter, *http.Request) *HTTPError

type ErrorHandler struct {
	Logger *logrus.Logger
}

func (eh *ErrorHandler) Wrap(handler customHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			w.WriteHeader(err.code)
			eh.Logger.Error(err)
			if err.code < http.StatusInternalServerError {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	})
}
