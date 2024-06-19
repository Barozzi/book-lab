package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

type BookClient struct {
	Context context.Context
}

func (bc BookClient) ByAuthor(ctx context.Context, author string, limit int) (model.Books, error) {
	query := fmt.Sprintf("inauthor:%s", replaceSpaces(author))
	return booksRequest(ctx, query, limit)
}

func (bc BookClient) ByTitle(ctx context.Context, title string, limit int) (model.Books, error) {
	query := fmt.Sprintf("intitle:%s", replaceSpaces(title))
	return booksRequest(ctx, query, limit)
}

func replaceSpaces(query string) string {
	return strings.ReplaceAll(query, " ", "+")
}

func booksRequest(ctx context.Context, query string, limit int) (model.Books, error) {
	url := "https://www.googleapis.com/books/v1/volumes?q=%s+langRestrict:en&maxResults=%d"

	fullurl := fmt.Sprintf(url, strings.ReplaceAll(query, " ", "+"), limit)
	fmt.Println(fullurl)
	res, err := http.Get(fullurl)
	if err != nil {
		return model.Books{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return model.Books{}, err
	}
	var books model.Books
	json.Unmarshal(body, &books)
	return books, nil
}
