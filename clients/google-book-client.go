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
	"regexp"
	"runtime"
	"slices"
	"strings"

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

const DEBUG = false

type GoogleBookRequest struct {
	Title  string
	Author string
	Start  int
	Limit  int
	Pages  int
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
			filterAuthorResults(request)(
				bc.bookRequest(ctx, query, request)))
	}
}

func (bc GoogleBookClient) ByTitle(ctx context.Context, request GoogleBookRequest) (model.GoogleBookResponse, error) {
	if bc.PactMode {
		slog.Info("serving pact")
		return titlePact()
	} else {
		query := fmt.Sprintf("intitle:%s+inauthor:%s", url.QueryEscape(request.Title), url.QueryEscape(request.Author))
		return sortByDescLength(
			filterTitleResults(request)(
				bc.bookRequest(ctx, query, request)))
	}
}

func filterTitleResults(req GoogleBookRequest) func(model.GoogleBookResponse, error) (model.GoogleBookResponse, error) {
	return func(resp model.GoogleBookResponse, err error) (model.GoogleBookResponse, error) {
		if err != nil {
			return resp, err
		}
		filteredBooks := []model.GoogleBookItem{}

		for _, book := range resp.Items {
			if DEBUG {
				fmt.Printf("filterTitleResults: exactTitle %s == %s is %t\n", book.VolumeInfo.Title, req.Title, (normalizeString(book.VolumeInfo.Title) == normalizeString(req.Title)))
				fmt.Printf("filterTitleResults: isEnglish %s == %s is %t\n", book.VolumeInfo.Language, "en", (book.VolumeInfo.Language == "en"))
				fmt.Printf("filterTitleResults: hasDesc %d > %d is %t\n", len(book.VolumeInfo.Description), 0, (len(book.VolumeInfo.Description) > 0))
				fmt.Printf("filterTitleResults: hasImage %d > %d is %t\n", len(book.VolumeInfo.ImageLinks.Thumbnail), 0, (len(book.VolumeInfo.ImageLinks.Thumbnail) > 0))
			}
			if filterExactTitle(book, req.Title) && filterIsEnglish(book) && filterHasDescription(book) && filterHasImage(book) {
				filteredBooks = append(filteredBooks, book)
			}
		}
		if len(filteredBooks) == 0 {
			fmt.Println("Overfiltered: ", req.Title)
			for _, book := range resp.Items {
				if filterCloseTitle(book, req.Title) && filterIsEnglish(book) && filterHasImage(book) {
					filteredBooks = append(filteredBooks, book)
				}
			}
		}

		// return filtered results unless all results were filtered
		resp.Items = filteredBooks
		return resp, err
	}
}

func filterAuthorResults(req GoogleBookRequest) func(model.GoogleBookResponse, error) (model.GoogleBookResponse, error) {
	return func(resp model.GoogleBookResponse, err error) (model.GoogleBookResponse, error) {
		if err != nil {
			return resp, err
		}
		filteredBooks := []model.GoogleBookItem{}

		for _, book := range resp.Items {
			if filterExactAuthor(book, req.Author) && filterIsEnglish(book) && filterHasDescription(book) && filterHasImage(book) {
				filteredBooks = append(filteredBooks, book)
			}
		}

		// return filtered results
		resp.Items = filteredBooks
		return resp, err
	}
}

func filterExactTitle(book model.GoogleBookItem, name string) bool {
	return normalizeString(book.VolumeInfo.Title) == normalizeString(name)
}

func filterCloseTitle(book model.GoogleBookItem, name string) bool {
	return strings.Contains(normalizeString(book.VolumeInfo.Title), normalizeString(name))
}

func filterExactAuthor(book model.GoogleBookItem, name string) bool {
	return slices.Contains(book.VolumeInfo.Authors, name)
}

func filterIsEnglish(book model.GoogleBookItem) bool {
	return book.VolumeInfo.Language == "en"
}

func filterHasDescription(book model.GoogleBookItem) bool {
	return len(book.VolumeInfo.Description) > 1
}

func filterHasImage(book model.GoogleBookItem) bool {
	return len(book.VolumeInfo.ImageLinks.Thumbnail) > 10
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

func sortByDescLength(resp model.GoogleBookResponse, err error) (model.GoogleBookResponse, error) {
	if err != nil {
		return resp, err
	}
	slices.SortFunc(resp.Items, func(a, b model.GoogleBookItem) int {
		return cmp.Compare(len(b.VolumeInfo.Description), len(a.VolumeInfo.Description))
	})
	return resp, err
}

func (bc GoogleBookClient) bookRequest(ctx context.Context, query string, request GoogleBookRequest) (model.GoogleBookResponse, error) {
	fullUrl := buildRequestUrl(query, request)
	slog.Info(fullUrl)

	// Make Request to Google Book API
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
func normalizeString(s string) string {
	s = strings.ToLower(s)
	reg, _ := regexp.Compile("[^a-zA-Z0-9]+")
	s = reg.ReplaceAllString(s, "")
	return s
}

func buildRequestUrl(query string, request GoogleBookRequest) string {
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
	return fullUrl
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
