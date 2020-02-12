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
		return fmt.Sprintf("HTTP %d: %s: %s", h.code, h.context, h.internal)
	}

	return fmt.Sprintf("HTTP %d: %s", h.code, h.internal)
}

func newHTTPError(code int, err error) *HTTPError {
	return &HTTPError{
		code:     code,
		internal: err,
	}
}

func (e *HTTPError) Wrap(ctx string) *HTTPError {
	e.context = ctx
	return e
}

func (e *HTTPError) Code() int {
	return e.code
}

func (e *HTTPError) PublicErrMsg(msg string) *HTTPError {
	e.publicMsg = msg
	return e
}

func (h HTTPError) publicError() string {
	if h.publicMsg != "" {
		return h.publicMsg
	}

	return h.Error()
}

type ErrorHandler struct {
	logger *logrus.Logger
}

func NewErrorHandler(loger *logrus.Logger) *ErrorHandler {
	return &ErrorHandler{loger}
}

func (eh ErrorHandler) Wrap(handler customHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var errorMsg string
		err := handler(w, r)

		if err != nil {
			eh.logger.Error(err)

			if err.code < http.StatusInternalServerError {
				errorMsg = err.publicError()
			}

			http.Error(w, errorMsg, err.code)
		}
	})
}
