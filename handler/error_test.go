package handler

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	tests := []struct {
		error   HTTPError
		message string
	}{
		{
			error: HTTPError{
				code:     400,
				internal: errors.New("Some error occured with code 400"),
				context:  "",
			},
			message: "HTTP 400: Some error occured with code 400",
		},
		{
			error: HTTPError{
				code:     400,
				internal: errors.New("Some error occured with code 400"),
				context:  "with some context",
			},
			message: "HTTP 400: with some context Some error occured with code 400",
		},
		{
			error: HTTPError{
				code:     500,
				internal: errors.New("Some error occured with code 500"),
				context:  "",
			},
			message: "HTTP 500: Some error occured with code 500",
		},
		{
			error: HTTPError{
				code:     500,
				internal: errors.New("Some error occured with code 500"),
				context:  "with some context",
			},
			message: "HTTP 500: with some context Some error occured with code 500",
		},
	}
	for _, tt := range tests {
		t.Run("test", func(t *testing.T) {
			errMsg := tt.error.Error()
			if tt.message != errMsg {
				t.Errorf("expected %q; got %q", tt.message, errMsg)
			}
		})
	}

}
