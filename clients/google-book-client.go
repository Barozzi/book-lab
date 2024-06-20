package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"example.com/book-learn/models"
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

type BookClientInterface interface {
	ByAuthor(ctx context.Context, author string, limit int) (model.GoogleBookResponse, error)
	ByTitle(ctx context.Context, title string, limit int) (model.GoogleBookResponse, error)
}

type GoogleBookClient struct {
	PactMode bool
}

func (bc GoogleBookClient) ByAuthor(ctx context.Context, author string, limit int) (model.GoogleBookResponse, error) {
	query := fmt.Sprintf("inauthor:%s", url.QueryEscape(author))
	if bc.PactMode == true {
		fmt.Println("serving pact")
		return authorPact()
	} else {
		return booksRequest(ctx, query, limit)
	}
}

func (bc GoogleBookClient) ByTitle(ctx context.Context, title string, limit int) (model.GoogleBookResponse, error) {
	query := fmt.Sprintf("intitle:%s", url.QueryEscape(title))
	if bc.PactMode == true {
		fmt.Println("serving pact")
		return titlePact()
	} else {
		return booksRequest(ctx, query, limit)
	}
}

func booksRequest(ctx context.Context, query string, limit int) (model.GoogleBookResponse, error) {
	url := "https://www.googleapis.com/books/v1/volumes?q=%s+langRestrict:en&maxResults=%d"

	fullurl := fmt.Sprintf(url, strings.ReplaceAll(query, " ", "+"), limit)
	fmt.Println(fullurl)
	res, err := http.Get(fullurl)
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
	pact, err := os.ReadFile("./clients/pacts/google-author-response.json")
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
	pact, err := os.ReadFile("./clients/pacts/google-title-response.json")
	if err != nil {
		return model.GoogleBookResponse{}, err
	}
	fmt.Println(string(pact)) // TODO DEBUG
	var books model.GoogleBookResponse
	err = json.Unmarshal(pact, &books)
	if err != nil {
		return model.GoogleBookResponse{}, err
	}
	fmt.Println(books) // TODO DEBUG
	return books, nil
}
