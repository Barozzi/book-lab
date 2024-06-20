package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func setupHealthRouter() http.Handler {
	r := chi.NewRouter()
	HealthRouter(r)
	return r
}

func TestHealthRouter(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Ping",
			method:         "GET",
			path:           "/ping",
			expectedStatus: http.StatusOK,
			expectedBody:   "pong",
		},
		{
			name:           "Info",
			method:         "GET",
			path:           "/info",
			expectedStatus: http.StatusOK,
			expectedBody:   "book-lab-api-v1",
		},
	}

	r := setupHealthRouter()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
