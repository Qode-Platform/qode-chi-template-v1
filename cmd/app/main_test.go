package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func code(path string) int {
	w := httptest.NewRecorder()
	newRouter().ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w.Code
}

func TestServesAtRoot(t *testing.T) {
	for _, path := range []string{"/health", "/"} {
		if got := code(path); got != http.StatusOK {
			t.Fatalf("GET %s = %d", path, got)
		}
	}
}
