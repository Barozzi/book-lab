package routes

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	model "example.com/book-learn/models"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// MockClient
type MockClient struct {
	Response model.Books
	Err      error
}

func (cli MockClient) ByAuthor(ctx context.Context, author string, limit int) (model.Books, error) {
	return cli.Response, cli.Err
}
func (cli MockClient) ByTitle(ctx context.Context, author string, limit int) (model.Books, error) {
	return cli.Response, cli.Err
}

// Mock Router
func setupBooksRouter(response model.Books, err error) http.Handler {
	r := chi.NewRouter()
	cli := MockClient{
		Response: response,
		Err:      err,
	}
	BooksRouter(r, cli)
	return r
}

func TestBooksRouter(t *testing.T) {
	booksReq := BooksRequest{
		Author: "test-author",
		Title:  "test-title",
	}
	testRequestBody, _ := json.Marshal(booksReq)

	mockAuthorItems := []model.GoogleBookItem{
		{
			VolumeInfo: model.GoogleBookVolumeInfo{
				Title: "test-title-1",
			},
		},
		{
			VolumeInfo: model.GoogleBookVolumeInfo{
				Title: "test-title-2",
			},
		},
	}

	expectedAuthorResponseBody := "{\"Author\":\"test-author\",\"Books\":[{\"Title\":\"test-title-1\",\"Authors\":null,\"PublishedDate\":\"\",\"Description\":\"\",\"PageCount\":0,\"Categories\":null,\"ContentVersion\":\"\",\"PanelizationSummary\":{\"containsEpubBubbles\":false,\"containsImageBubbles\":false},\"ImageLinks\":{\"smallThumbnail\":\"\",\"thumbnail\":\"\"},\"Language\":\"\",\"PreviewLink\":\"\",\"InfoLink\":\"\",\"CanonicalVolumeLink\":\"\"},{\"Title\":\"test-title-2\",\"Authors\":null,\"PublishedDate\":\"\",\"Description\":\"\",\"PageCount\":0,\"Categories\":null,\"ContentVersion\":\"\",\"PanelizationSummary\":{\"containsEpubBubbles\":false,\"containsImageBubbles\":false},\"ImageLinks\":{\"smallThumbnail\":\"\",\"thumbnail\":\"\"},\"Language\":\"\",\"PreviewLink\":\"\",\"InfoLink\":\"\",\"CanonicalVolumeLink\":\"\"}]}"
	mockClientAuthorResponse := model.Books{
		Kind:       "mock-kind",
		TotalItems: 47,
		Items:      mockAuthorItems,
	}

	expectedTitleResponseBody := "{\"Title\":\"test-title\",\"Books\":[{\"Title\":\"test-title-1\",\"Authors\":[\"test-author\"],\"PublishedDate\":\"\",\"Description\":\"\",\"PageCount\":0,\"Categories\":null,\"ContentVersion\":\"\",\"PanelizationSummary\":{\"containsEpubBubbles\":false,\"containsImageBubbles\":false},\"ImageLinks\":{\"smallThumbnail\":\"\",\"thumbnail\":\"\"},\"Language\":\"\",\"PreviewLink\":\"\",\"InfoLink\":\"\",\"CanonicalVolumeLink\":\"\"}]}"
	mockClientTitleResponse := model.Books{
		Kind:       "mock-kind",
		TotalItems: 47,
		Items: []model.GoogleBookItem{
			{
				VolumeInfo: model.GoogleBookVolumeInfo{
					Authors: []string{"test-author"},
					Title:   "test-title-1",
				},
			},
		},
	}

	tests := []struct {
		name               string
		method             string
		path               string
		expectedStatus     int
		mockClientResponse model.Books
		mockClientError    error
		testRequestBody    []byte
		expectedBody       string
	}{
		{
			name:               "POST:/books/author",
			method:             "POST",
			path:               "/books/author",
			mockClientResponse: mockClientAuthorResponse,
			mockClientError:    nil,
			expectedStatus:     http.StatusOK,
			testRequestBody:    testRequestBody,
			expectedBody:       expectedAuthorResponseBody,
		},
		{
			name:               "POST:/books/title",
			method:             "POST",
			path:               "/books/title",
			mockClientResponse: mockClientTitleResponse,
			mockClientError:    nil,
			expectedStatus:     http.StatusOK,
			testRequestBody:    testRequestBody,
			expectedBody:       expectedTitleResponseBody,
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
			assert.Equal(t, string(tt.expectedBody), w.Body.String())
		})
	}
}
