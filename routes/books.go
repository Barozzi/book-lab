package routes

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	client "example.com/book-learn/clients"
	model "example.com/book-learn/models"
	"github.com/go-chi/chi/v5"
)

type AuthorResponse struct {
	Author     string         `json:"author"`
	TotalItems int            `json:"totalItems"`
	Books      []BookResponse `json:"books"`
}

type TitleResponse struct {
	Title      string         `json:"title"`
	TotalItems int            `json:"totalItems"`
	Books      []BookResponse `json:"books"`
}

type BookResponse struct {
	Title               string                              `json:"title"`
	Authors             []string                            `json:"authors"`
	PublishedDate       string                              `json:"publishedDate"`
	Description         string                              `json:"description"`
	PageCount           int                                 `json:"pageCount"`
	Categories          []string                            `json:"categories"`
	ContentVersion      string                              `json:"contentVersion"`
	PanelizationSummary model.GoogleBookPanelizationSummary `json:"panelizationSummary"`
	ImageLinks          model.GoogleBookImageLinks          `json:"imageLinks"`
	Language            string                              `json:"language"`
	PreviewLink         string                              `json:"previewLink"`
	InfoLink            string                              `json:"infoLink"`
	CanonicalVolumeLink string                              `json:"canonicalVolumeLink"`
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
		var bookReq client.GoogleBookRequest
		err := json.NewDecoder(r.Body).Decode(&bookReq)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// Fetch data from external API
		books, err := bookClient.ByAuthor(context.Background(), bookReq)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// No results
		if len(books.Items) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Format response
		var bookResp AuthorResponse
		bookResp.Author = bookReq.Author
		bookResp.TotalItems = books.TotalItems
		for _, book := range books.Items {
			var br BookResponse
			br.fromVolumeInfo(book.VolumeInfo)
			bookResp.Books = append(bookResp.Books, br)
		}

		// write de jaysawn
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(bookResp); err != nil {
			slog.Error(err.Error())
			http.Error(w, "", http.StatusInternalServerError)
		}
	}
}

func queryByTitle(bookClient client.BookClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse the request body
		var bookReq client.GoogleBookRequest
		err := json.NewDecoder(r.Body).Decode(&bookReq)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "", http.StatusBadRequest)
			return
		}

		// Fetch data from external API
		books, err := bookClient.ByTitle(context.Background(), bookReq)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		// No results
		if len(books.Items) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Format Response
		var resp TitleResponse
		resp.Title = bookReq.Title
		resp.TotalItems = books.TotalItems
		for _, book := range books.Items {
			var br BookResponse
			br.fromVolumeInfo(book.VolumeInfo)
			resp.Books = append(resp.Books, br)
		}

		// Gift the findings to our user, but this is an internal api so gift to me
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(resp); err != nil {
			slog.Error(err.Error())
			http.Error(w, "", http.StatusInternalServerError)
		}
	}
}
