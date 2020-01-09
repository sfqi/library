package handler

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
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
			eh.Logger.Error(err) //1. logging
			if err.code < http.StatusInternalServerError {
				fmt.Fprintf(w, err.Error())
			}
		}
	})
}
