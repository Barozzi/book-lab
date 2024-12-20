package routes

import (
	"context"
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"sync"

	client "example.com/book-learn/clients"
	model "example.com/book-learn/models"
	"github.com/go-chi/chi/v5"
)

type AuthorResponse struct {
	Author       string         `json:"author"`
	TotalItems   int            `json:"totalItems"`
	Books        []BookResponse `json:"books"`
	HasMorePages bool           `json:"hasMorePages"`
}

type TitleResponse struct {
	Title      string         `json:"title"`
	TotalItems int            `json:"totalItems"`
	Books      []BookResponse `json:"books"`
}

type BookResponse struct {
	Title               string                  `json:"title"`
	Authors             []string                `json:"authors"`
	PublishedDate       string                  `json:"publishedDate"`
	Description         string                  `json:"description"`
	PageCount           int                     `json:"pageCount"`
	Categories          []string                `json:"categories"`
	ContentVersion      string                  `json:"contentVersion"`
	PanelizationSummary BookPanelizationSummary `json:"panelizationSummary"`
	ImageLinks          BookImageLinks          `json:"imageLinks"`
	Language            string                  `json:"language"`
	PreviewLink         string                  `json:"previewLink"`
	InfoLink            string                  `json:"infoLink"`
	CanonicalVolumeLink string                  `json:"canonicalVolumeLink"`
}

type BookPanelizationSummary struct {
	ContainsEpubBubbles  bool `json:"containsEpubBubbles"`
	ContainsImageBubbles bool `json:"containsImageBubbles"`
}

type BookImageLinks struct {
	SmallThumbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
}

func (br *BookResponse) fromVolumeInfo(vi model.GoogleBookVolumeInfo) {
	br.Title = vi.Title
	br.Authors = vi.Authors
	br.PublishedDate = vi.PublishedDate
	br.Description = vi.Description
	br.PageCount = vi.PageCount
	br.Categories = vi.Categories
	br.ContentVersion = vi.ContentVersion
	br.PanelizationSummary = BookPanelizationSummary{
		ContainsEpubBubbles:  vi.PanelizationSummary.ContainsEpubBubbles,
		ContainsImageBubbles: vi.PanelizationSummary.ContainsImageBubbles,
	}
	br.ImageLinks = BookImageLinks{
		SmallThumbnail: vi.ImageLinks.SmallThumbnail,
		Thumbnail:      vi.ImageLinks.Thumbnail,
	}
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
		slog.Info("BookRequest:", "Author", bookReq.Author, "Start", strconv.Itoa(bookReq.Start), "limit", strconv.Itoa(bookReq.Limit), "Pages", strconv.Itoa(bookReq.Pages))

		results := make(chan model.GoogleBookResponse, bookReq.Pages+1)
		var wg sync.WaitGroup

		fetch := func(start int) {
			defer wg.Done()
			// Fetch data from external API
			req := client.GoogleBookRequest{
				Title:  bookReq.Title,
				Author: bookReq.Author,
				Start:  start,
				Limit:  bookReq.Limit,
				Pages:  0,
			}
			books, err := bookClient.ByAuthor(context.Background(), req)
			if err != nil {
				slog.Error(err.Error())
				http.Error(w, "", http.StatusInternalServerError)
				return
			}
			slog.Info(req.Author, "Start", strconv.Itoa(req.Start), "limit", strconv.Itoa(req.Limit), "Pages", strconv.Itoa(req.Pages))

			results <- books
		}
		wg.Add(1)
		go fetch(bookReq.Start)

		for i := 0; i <= bookReq.Pages; i++ {
			wg.Add(1)
			go fetch(bookReq.Start + i*bookReq.Limit)
		}

		go func() {
			wg.Wait()
			close(results)
		}()

		var books []model.GoogleBookItem
		totalItems := 0

		for result := range results {
			if totalItems == 0 {
				totalItems = result.TotalItems
			}
			books = append(books, result.Items...)
		}

		// No results
		if len(books) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Format response
		var bookResp AuthorResponse
		bookResp.Author = bookReq.Author
		bookResp.TotalItems = totalItems
		bookResp.HasMorePages = int(math.Ceil(float64(totalItems)/float64(bookReq.Limit))-float64(bookReq.Pages)) > 0
		for _, book := range books {
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
