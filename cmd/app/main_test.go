package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func code(path string) int {
	w := httptest.NewRecorder()
	newRouter().ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w.Code
}

func TestServesAtRootWhenUnset(t *testing.T) {
	os.Unsetenv("BASE_PATH")
	if got := code("/health"); got != http.StatusOK {
		t.Fatalf("GET /health = %d", got)
	}
}

func TestColonInPrefixIsLiteral(t *testing.T) {
	t.Setenv("BASE_PATH", "/direct/agent-7:3000")
	if got := code("/direct/agent-7:3000/health"); got != http.StatusOK {
		t.Fatalf("prefixed = %d", got)
	}
	if got := code("/health"); got != http.StatusNotFound {
		t.Fatalf("bare = %d, want 404", got)
	}
}
