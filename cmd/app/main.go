// A minimal chi service shaped for the fleet.
//
// Note for anyone porting this pattern: chi uses {param} syntax, so the colon
// in BASE_PATH (/direct/<agent>:<port>) is a literal path character here. The
// same prefix crashes path-to-regexp-based routers, which read ":3000" as a
// parameter.
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

// basePath returns the fleet's ingress prefix, normalised to "" or
// "/leading/no-trailing-slash". The fleet injects BASE_PATH as
// /direct/<agent>:<port> and nginx forwards it UNCHANGED, so every route must
// live under it. Empty means standalone: serve at the host root.
func basePath() string {
	raw := strings.Trim(strings.TrimSpace(os.Getenv("BASE_PATH")), "/")
	if raw == "" {
		return ""
	}
	return "/" + raw
}

func port() string {
	if p := strings.TrimSpace(os.Getenv("PORT")); p != "" {
		return p
	}
	return "8080"
}

func newRouter() http.Handler {
	app := chi.NewRouter()
	app.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"status": "ok"})
	})
	app.Get("/", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]string{"service": "chi-template", "base_path": basePath()})
	})

	root := chi.NewRouter()
	root.Use(middleware.Recoverer)
	if bp := basePath(); bp != "" {
		root.Mount(bp, app)
	} else {
		root.Mount("/", app)
	}
	return root
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func main() {
	addr := ":" + port()
	log.Printf("chi-template listening on %s (base_path=%q)", addr, basePath())
	if err := http.ListenAndServe(addr, newRouter()); err != nil {
		log.Fatal(err)
	}
}
