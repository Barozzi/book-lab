package main

import (
	"fmt"
	"net/http"
	"os"

	client "example.com/book-learn/clients"
	"example.com/book-learn/routes"
	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	bookClient := client.GoogleBookClient{
		GetData:  http.Get,
		PactMode: os.Getenv("PACT_MODE") == "true",
	}

	r.Route("/api", func(r chi.Router) {
		routes.BooksRouter(r, bookClient)
		routes.HealthRouter(r)
	})

	// Server it up
	fmt.Println("listening on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
