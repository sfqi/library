package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"
)

func BodyDump(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, err := httputil.DumpRequest(r, true)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		body := string(dump)
		fmt.Println(body)
		next.ServeHTTP(w, r)
	})
}
