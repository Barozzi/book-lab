package main

import (
	"fmt"
	"net/http"

	"example.com/book-learn/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		routes.BooksRouter(r)
		routes.HealthRouter(r)
	})

	// Server it up
	fmt.Println("listening on http://localhost:8080/api")
	http.ListenAndServe(":8080", r)
}
