package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupBooksRouter() http.Handler {
	r := chi.NewRouter()
	BooksRouter(r)
	return r
}

func TestPostBooksAuthor(t *testing.T) {
	r := setupBooksRouter()

	booksReq := BooksRequest{
		Author: "test-author",
	}
	requestBody, _ := json.Marshal(booksReq)

	req, _ := http.NewRequest("POST", "/books/author", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	expectedResponse := AuthorResponse{
		Author: "test-author",
		Books:  nil,
	}
	expectedResponseBody, _ := json.Marshal(expectedResponse)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expectedResponseBody), w.Body.String())
}

func TestPostBooksTitle(t *testing.T) {
	r := setupBooksRouter()

	booksReq := BooksRequest{
		Title: "test-title",
	}
	requestBody, _ := json.Marshal(booksReq)

	req, _ := http.NewRequest("POST", "/books/title", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	expectedResponse := TitleResponse{
		Title: "test-title",
		Books: nil,
	}
	expectedResponseBody, _ := json.Marshal(expectedResponse)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, string(expectedResponseBody), w.Body.String())
}

func Test_setupBooksRouter(t *testing.T) {
	tests := []struct {
		name string
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setupBooksRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setupBooksRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}
