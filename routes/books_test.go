package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	client "example.com/book-learn/clients"
	model "example.com/book-learn/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// MockClient
type MockClient struct {
	Response model.GoogleBookResponse
	Err      error
}

func (cli MockClient) ByAuthor(ctx context.Context, request client.GoogleBookRequest) (model.GoogleBookResponse, error) {
	return cli.Response, cli.Err
}
func (cli MockClient) ByTitle(ctx context.Context, request client.GoogleBookRequest) (model.GoogleBookResponse, error) {
	return cli.Response, cli.Err
}

func setupBooksRouter(response model.GoogleBookResponse, err error) http.Handler {
	r := chi.NewRouter()
	cli := MockClient{
		Response: response,
		Err:      err,
	}
	BooksRouter(r, cli)
	return r
}

func TestBooksRouter(t *testing.T) {
	authorReq := client.GoogleBookRequest{
		Author: "test-author",
	}
	titleReq := client.GoogleBookRequest{
		Title: "test-title",
	}
	testAuthorRequestBody, _ := json.Marshal(authorReq)
	testTitleRequestBody, _ := json.Marshal(titleReq)

	mockItems := []model.GoogleBookItem{
		{
			Kind:     "kind",
			ID:       "id",
			Etag:     "etag",
			SelfLink: "selflink",
		},
	}
	mockEmptyItems := []model.GoogleBookItem{}

	tests := []struct {
		name               string
		method             string
		path               string
		expectedStatus     int
		mockClientResponse model.GoogleBookResponse
		mockClientError    error
		testRequestBody    []byte
	}{
		{
			name:   "POST:/books/author with valid client response",
			method: "POST",
			path:   "/books/author",
			mockClientResponse: model.GoogleBookResponse{
				Kind:       "test",
				TotalItems: 42,
				Items:      mockItems,
			},
			mockClientError: nil,
			expectedStatus:  http.StatusOK,
			testRequestBody: testAuthorRequestBody,
		},
		{
			name:   "POST:/books/author with empty client response",
			method: "POST",
			path:   "/books/author",
			mockClientResponse: model.GoogleBookResponse{
				Kind:       "test",
				TotalItems: 0,
				Items:      mockEmptyItems,
			},
			mockClientError: nil,
			expectedStatus:  http.StatusNoContent,
			testRequestBody: testAuthorRequestBody,
		},
		{
			name:               "POST:/books/author with client error",
			method:             "POST",
			path:               "/books/author",
			mockClientResponse: model.GoogleBookResponse{},
			mockClientError:    errors.New("test-error"),
			expectedStatus:     http.StatusInternalServerError,
			testRequestBody:    testAuthorRequestBody,
		},
		{
			name:   "POST:/books/title with valid client response",
			method: "POST",
			path:   "/books/title",
			mockClientResponse: model.GoogleBookResponse{
				Kind:       "test",
				TotalItems: 42,
				Items:      mockItems,
			},
			mockClientError: nil,
			expectedStatus:  http.StatusOK,
			testRequestBody: testTitleRequestBody,
		},
		{
			name:   "POST:/books/title with empty cient response",
			method: "POST",
			path:   "/books/title",
			mockClientResponse: model.GoogleBookResponse{
				Kind:       "test",
				TotalItems: 0,
				Items:      mockEmptyItems,
			},
			mockClientError: nil,
			expectedStatus:  http.StatusNoContent,
			testRequestBody: testTitleRequestBody,
		},
		{
			name:               "POST:/books/title with client error",
			method:             "POST",
			path:               "/books/title",
			mockClientResponse: model.GoogleBookResponse{},
			mockClientError:    errors.New("test-error"),
			expectedStatus:     http.StatusInternalServerError,
			testRequestBody:    testTitleRequestBody,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, bytes.NewBuffer(tt.testRequestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r := setupBooksRouter(tt.mockClientResponse, tt.mockClientError)
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.NotEmpty(t, w.Body.String())
			}
		})
	}
}
