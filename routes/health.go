package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func HealthRouter(r chi.Router) {
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "pong")
	})
	r.Get("/info", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "book-lab-api-v1")
	})
}
