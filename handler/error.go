package handler

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
)

type customHandler func(http.ResponseWriter, *http.Request) *HTTPError

type HTTPError struct {
	code      int
	internal  error
	context   string
	publicMsg string
}

func (h HTTPError) Error() string {
	if h.context != "" {
		return fmt.Sprintf("HTTP %d: %s: %s: %s", h.code, h.context, h.internal, h.publicMsg)
	}
	return fmt.Sprintf("HTTP %d: %s: %s", h.code, h.internal, h.publicMsg)
}

func newHTTPError(code int, err error, msg string) *HTTPError {
	return &HTTPError{
		code:      code,
		internal:  err,
		publicMsg: msg,
	}
}

func (e *HTTPError) Wrap(ctx string) *HTTPError {
	e.context = ctx
	return e
}

func (e *HTTPError) Code() int {
	return e.code
}

type ErrorHandler struct {
	Logger *logrus.Logger
}

func (eh ErrorHandler) Wrap(handler customHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorMsg string
		err := handler(w, r)

		if err != nil {
			eh.Logger.Error(err)

			if err.code < http.StatusInternalServerError {
				errorMsg = err.publicMsg
			}

			http.Error(w, errorMsg, err.code)
		}
	})
}
