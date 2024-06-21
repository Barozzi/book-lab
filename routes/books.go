package routes

import (
	"context"
	"encoding/json"
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
	PanelizationSummary model.GoogleBookPanelizationSummary
	ImageLinks          model.GoogleBookImageLinks
	Language            string
	PreviewLink         string
	InfoLink            string
	CanonicalVolumeLink string
}

func (br *BookResponse) fromVolumeInfo(vi model.GoogleBookVolumeInfo) {
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

func BooksRouter(r chi.Router, api client.BookClientInterface) {
	r.Post("/books/author", queryByAuthor(api))
	r.Post("/books/title", queryByTitle(api))
}

func queryByAuthor(bookClient client.BookClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var bookReq BooksRequest
		err := json.NewDecoder(r.Body).Decode(&bookReq)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Fetch data from external API
		books, err := bookClient.ByAuthor(context.Background(), bookReq.Author, 25)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// No results
		if len(books.Items) == 0 {
			w.WriteHeader(http.StatusNoContent)
		}

		// Format response
		var bookResp AuthorResponse
		bookResp.Author = bookReq.Author
		for _, book := range books.Items {
			var br BookResponse
			br.fromVolumeInfo(book.VolumeInfo)
			bookResp.Books = append(bookResp.Books, br)
		}

		// write de jaysawn
		if err := json.NewEncoder(w).Encode(bookResp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func queryByTitle(bookClient client.BookClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var br BooksRequest
		err := json.NewDecoder(r.Body).Decode(&br)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Fetch data from external API
		books, err := bookClient.ByTitle(context.Background(), br.Title, 1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// No results
		if len(books.Items) == 0 {
			w.WriteHeader(http.StatusNoContent)
		}

		// Format Response
		var resp TitleResponse
		resp.Title = br.Title
		for _, book := range books.Items {
			var br BookResponse
			br.fromVolumeInfo(book.VolumeInfo)
			resp.Books = append(resp.Books, br)
		}

		// Gift the findings to our user, but this is an internal api so gift to me
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
