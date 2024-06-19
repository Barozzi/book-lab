package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	client "example.com/book-learn/clients"
	model "example.com/book-learn/models"
	"github.com/go-chi/chi/v5"
)

type BooksRequest struct {
	Author string `json:"author"`
	Title  string `json:"title"`
}

type AuthorResponse struct {
	Author string
	Books  []BookResponse
}

type TitleResponse struct {
	Title string
	Books []BookResponse
}

type BookResponse struct {
	Title               string
	Authors             []string
	PublishedDate       string
	Description         string
	PageCount           int
	Categories          []string
	ContentVersion      string
	PanelizationSummary model.PanelizationSummary
	ImageLinks          model.ImageLinks
	Language            string
	PreviewLink         string
	InfoLink            string
	CanonicalVolumeLink string
}

func (br *BookResponse) fromVolumeInfo(vi model.VolumeInfo) {
	br.Title = vi.Title
	br.Authors = vi.Authors
	br.PublishedDate = vi.PublishedDate
	br.Description = vi.Description
	br.PageCount = vi.PageCount
	br.Categories = vi.Categories
	br.ContentVersion = vi.ContentVersion
	br.PanelizationSummary = vi.PanelizationSummary
	br.ImageLinks = vi.ImageLinks
	br.Language = vi.Language
	br.PreviewLink = vi.PreviewLink
	br.InfoLink = vi.InfoLink
	br.CanonicalVolumeLink = vi.CanonicalVolumeLink
}

func BooksRouter(r chi.Router) {
	r.Post("/books/author", queryByAuthor)
	r.Post("/books/title", queryByTitle)
}

func queryByAuthor(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var br BooksRequest
	err := json.NewDecoder(r.Body).Decode(&br)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch data from external API
	var bc client.BookClient
	books, err := bc.ByAuthor(context.Background(), br.Author, 25)
	if err != nil {
		fmt.Printf("error fetching book data from external api: %s", err.Error())
		return
	}

	// Format as JSON
	var resp AuthorResponse
	resp.Author = br.Author
	for _, book := range books.Items {
		var br BookResponse
		br.fromVolumeInfo(book.VolumeInfo)
		resp.Books = append(resp.Books, br)
	}
	jsonData, err := json.Marshal(resp)
	fmt.Fprint(w, string(jsonData))
}

func queryByTitle(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var br BooksRequest
	err := json.NewDecoder(r.Body).Decode(&br)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch data from external API
	var bc client.BookClient
	books, err := bc.ByTitle(context.Background(), br.Title, 1)
	if err != nil {
		fmt.Printf("error fetching book data from external api: %s", err.Error())
		return
	}

	// Format as JSON
	var resp TitleResponse
	resp.Title = br.Title
	for _, book := range books.Items {
		var br BookResponse
		br.fromVolumeInfo(book.VolumeInfo)
		resp.Books = append(resp.Books, br)
	}
	jsonData, err := json.Marshal(resp)
	fmt.Fprint(w, string(jsonData))
}
