package client

import (
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
		query := fmt.Sprintf("inauthor:%s+langRestrict:en", url.QueryEscape(request.Author))
		return bc.bookRequest(ctx, query, request)
	}
}

func (bc GoogleBookClient) ByTitle(ctx context.Context, request GoogleBookRequest) (model.GoogleBookResponse, error) {
	if bc.PactMode == true {
		slog.Info("serving pact")
		return titlePact()
	} else {
		query := fmt.Sprintf("intitle:%s+langRestrict:en", url.QueryEscape(request.Title))
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
