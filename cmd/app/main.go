// A minimal chi service shaped for the fleet: it serves at the root of its own
// hostname, so routes mount directly on the router.
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func port() string {
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		return p
	}
	return "8080"
}

func newRouter() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"service": "chi-template"})
	})
	return r
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	addr := ":" + port()
	log.Printf("chi-template listening on %s", addr)
	if err := http.ListenAndServe(addr, newRouter()); err != nil {
		log.Fatal(err)
	}
}
