package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		error   HTTPError
		message string
	}{
		{
			name: "Status code 400 with no context",
			error: HTTPError{
				code:     400,
				internal: errors.New("Some error occured with code 400"),
				context:  "",
			},
			message: "HTTP 400: Some error occured with code 400",
		},
		{
			name: "Status code 400 with context",
			error: HTTPError{
				code:     400,
				internal: errors.New("Some error occured with code 400"),
				context:  "with some context",
			},
			message: "HTTP 400: with some context Some error occured with code 400",
		},
		{
			name: "Status code 500 with no context",
			error: HTTPError{
				code:     500,
				internal: errors.New("Some error occured with code 500"),
				context:  "",
			},
			message: "HTTP 500: Some error occured with code 500",
		},
		{
			name: "Status code 500 with context",
			error: HTTPError{
				code:     500,
				internal: errors.New("Some error occured with code 500"),
				context:  "with some context",
			},
			message: "HTTP 500: with some context Some error occured with code 500",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errMsg := tt.error.Error()
			assert.Equal(tt.message, errMsg)
		})
	}

}

func testHandler(w http.ResponseWriter, r *http.Request) *HTTPError {
	w.WriteHeader(http.StatusOK)
	return nil
}

func testHandler2(w http.ResponseWriter, r *http.Request) *HTTPError {
	return &HTTPError{
		code:     500,
		internal: errors.New("error with status code 500"),
		context:  "",
	}
}
func testHandler3(w http.ResponseWriter, r *http.Request) *HTTPError {
	return &HTTPError{
		code:     400,
		internal: errors.New("error with status code 400"),
		context:  " with some context",
	}
}

func TestWrap(t *testing.T) {
	assert := assert.New(t)
	t.Run("Test Wrap", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/testerrorhandler", nil)
		if err != nil {
			t.Fatal(err)
		}
		logger := logrus.New()
		rr := httptest.NewRecorder()
		errHandler := ErrorHandler{logger}
		newHandler := errHandler.Wrap(testHandler)
		newHandler.ServeHTTP(rr, req)
		assert.Equal(http.StatusOK, rr.Code)
		fmt.Println(rr.Code)
	})
	t.Run("Status code 500", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/testerrorhandler2", nil)
		if err != nil {
			t.Fatal(err)
		}
		logger := logrus.New()
		rr := httptest.NewRecorder()
		errHandler := ErrorHandler{logger}
		newHandler := errHandler.Wrap(testHandler2)
		newHandler.ServeHTTP(rr, req)
		assert.Equal(http.StatusInternalServerError, rr.Code)
		expected := "\n"
		assert.Equal(expected, rr.Body.String())
	})
	t.Run("Status code 500", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/testerrorhandler3", nil)
		if err != nil {
			t.Fatal(err)
		}
		logger := logrus.New()
		rr := httptest.NewRecorder()
		errHandler := ErrorHandler{logger}
		newHandler := errHandler.Wrap(testHandler3)
		newHandler.ServeHTTP(rr, req)
		assert.Equal(http.StatusBadRequest, rr.Code)
		expected := "HTTP 400:  with some context error with status code 400\n"
		assert.Equal(expected, rr.Body.String())
	})

}
