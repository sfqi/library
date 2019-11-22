package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/sirupsen/logrus"
)

type BodyDump struct {
	Logger *logrus.Logger
}

func (b BodyDump) Dump(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			b.Logger.Error(err)
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		body := string(dump)

		b.Logger.Info(strings.ReplaceAll(body, "\r\n", "; "))

		next.ServeHTTP(w, r)
	})
}
