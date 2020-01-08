package handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type customHandler func(http.ResponseWriter, *http.Request) error

type ErrorHandler struct {
	Logger *logrus.Logger
}

func (eh *ErrorHandler) Wrap(handler customHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		eh.Logger.Error(err)
	})
}
