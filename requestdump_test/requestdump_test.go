package requestdump_test

import (
	 "fmt"
	"net/http"
	 "net/http/httptest"
	"testing"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestDump(t *testing.T) {
	req, err := http.NewRequest("GET", "/testmiddleware", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler := http.HandlerFunc(testHandler)
	rr := httptest.NewRecorder()
	newHandler := Dump(handler)
	newHandler.ServeHTTP(rr, req)
	fmt.Println(rr.Code)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected code to be %d, got %d", http.StatusOK, rr.Code)
	}
}