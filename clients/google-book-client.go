package client

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	model "example.com/book-learn/models"
)

// Common Query Parameters:
// q: The search query parameter. You can specify search terms here, including combinations of terms.
// intitle: Restricts the search to books with the specified words in the title.
// inpublisher: Restricts the search to books with the specified words in the publisher's name.
// inauthor: Restricts the search to books with the specified words in the author's name.
// isbn: Restricts the search to books with the specified ISBN.
// lccn: Restricts the search to books with the specified LCCN (Library of Congress Control Number).
// oclc: Restricts the search to books with the specified OCLC number.
// langRestrict: Restricts the search to books with the specified language code (e.g., en for English).
// printType: Restricts the results to books or magazines (values: books, magazines).
// filter: Applies a filter to the results (values: ebooks, free-ebooks, full, paid-ebooks, partial).
// maxResults: Maximum number of results to return (default is 10, maximum is 40).
// startIndex: Index of the first result to return (for pagination).
// orderBy: Specifies how the results should be sorted (values: relevance, newest).

type GoogleBookRequest struct {
	Title  string
	Author string
	Start  int
	Limit  int
}
type BookClientInterface interface {
	ByAuthor(ctx context.Context, request GoogleBookRequest) (model.GoogleBookResponse, error)
	ByTitle(ctx context.Context, request GoogleBookRequest) (model.GoogleBookResponse, error)
}

type GoogleBookClient struct {
	GetData  func(url string) (resp *http.Response, err error)
	PactMode bool
}

func (bc GoogleBookClient) ByAuthor(ctx context.Context, request GoogleBookRequest) (model.GoogleBookResponse, error) {
	if bc.PactMode == true {
		slog.Info("serving pact")
		return authorPact()
	} else {
		query := fmt.Sprintf("inauthor:\"%s\"", url.QueryEscape(request.Author))
		return sortByPublishedDate(
			filterResults(request)(
				bc.bookRequest(ctx, query, request)))
	}
}

func filterResults(req GoogleBookRequest) func(model.GoogleBookResponse, error) (model.GoogleBookResponse, error) {
	return func(resp model.GoogleBookResponse, err error) (model.GoogleBookResponse, error) {
		if err != nil {
			return resp, err
		}
		var filteredBooks []model.GoogleBookItem

		for _, book := range resp.Items {
			if filterExactAuthor(book, req.Author) && filterNonEnglish(book) && filterNoDescription(book) && filterNoImage(book) {
				filteredBooks = append(filteredBooks, book)
			}
		}

		// return filtered results
		resp.Items = filteredBooks
		return resp, err
	}
}

func filterExactAuthor(book model.GoogleBookItem, name string) bool {
	return slices.Contains(book.VolumeInfo.Authors, name)
}

func filterNonEnglish(book model.GoogleBookItem) bool {
	return book.VolumeInfo.Language == "en"
}

func filterNoDescription(book model.GoogleBookItem) bool {
	return len(book.VolumeInfo.Description) > 1
}

func filterNoImage(book model.GoogleBookItem) bool {
	return len(book.VolumeInfo.ImageLinks.Thumbnail) > 10
}

func dedupe(resp model.GoogleBookResponse, err error) (model.GoogleBookResponse, error) {
	if err != nil {
		return resp, err
	}
	return resp, err
}

func sortByPublishedDate(resp model.GoogleBookResponse, err error) (model.GoogleBookResponse, error) {
	if err != nil {
		return resp, err
	}
	slices.SortFunc(resp.Items, func(a, b model.GoogleBookItem) int {
		return cmp.Compare(b.VolumeInfo.PublishedDate, a.VolumeInfo.PublishedDate)
	})
	return resp, err
}

func (bc GoogleBookClient) ByTitle(ctx context.Context, request GoogleBookRequest) (model.GoogleBookResponse, error) {
	if bc.PactMode == true {
		slog.Info("serving pact")
		return titlePact()
	} else {
		query := fmt.Sprintf("intitle:%s", url.QueryEscape(request.Title))
		return bc.bookRequest(ctx, query, request)
	}
}

func (bc GoogleBookClient) bookRequest(ctx context.Context, query string, request GoogleBookRequest) (model.GoogleBookResponse, error) {
	type requestPart struct {
		querystring string
		valid       bool
	}
	queryParts := []requestPart{
		{querystring: fmt.Sprintf("&startIndex=%s", url.QueryEscape(fmt.Sprint(request.Start))), valid: request.Start > 0},
		{querystring: fmt.Sprintf("&maxResults=%s", url.QueryEscape(fmt.Sprint(request.Limit))), valid: request.Limit > 0},
	}

	fullUrl := fmt.Sprintf("https://www.googleapis.com/books/v1/volumes?q=%s", query)
	for _, part := range queryParts {
		if part.valid == true {
			fullUrl = fmt.Sprintf("%s%s", fullUrl, part.querystring)
		}
	}
	slog.Info(fullUrl)
	res, err := bc.GetData(fullUrl)
	if err != nil {
		return model.GoogleBookResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return model.GoogleBookResponse{}, err
	}
	var books model.GoogleBookResponse
	json.Unmarshal(body, &books)
	return books, nil
}

func authorPact() (model.GoogleBookResponse, error) {
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)
	fpath := filepath.Join(basepath, "/pacts/google-author-response.json")
	pact, err := os.ReadFile(fpath)
	if err != nil {
		return model.GoogleBookResponse{}, err
	}
	var books model.GoogleBookResponse
	err = json.Unmarshal(pact, &books)
	if err != nil {
		return model.GoogleBookResponse{}, err
	}
	return books, nil
}

func titlePact() (model.GoogleBookResponse, error) {
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)
	fpath := filepath.Join(basepath, "/pacts/google-title-response.json")
	pact, err := os.ReadFile(fpath)
	if err != nil {
		return model.GoogleBookResponse{}, err
	}
	var books model.GoogleBookResponse
	err = json.Unmarshal(pact, &books)
	if err != nil {
		return model.GoogleBookResponse{}, err
	}
	return books, nil
}
