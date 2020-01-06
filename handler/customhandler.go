package handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type customHandler func(http.ResponseWriter, *http.Request) error

func (ch *customHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

type ErrorHandler struct {
	Logger *logrus.Logger
}

func (eh *ErrorHandler) Wrap(handler customHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		eh.Logger.Println(r.URL.Path + r.Method)
		handler.ServeHTTP(w, r)
	})
}
