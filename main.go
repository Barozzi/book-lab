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
	})

	// Server it up
	fmt.Println("listening on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
