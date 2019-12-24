package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/sfqi/library/middleware"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestDump(t *testing.T) {
	assert := assert.New(t)
	req, err := http.NewRequest("GET", "/testmiddleware", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(testHandler)
	logger := logrus.New()
	rr := httptest.NewRecorder()
	newHandler := middleware.BodyDump{logger}.Dump(handler)
	newHandler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)
}
