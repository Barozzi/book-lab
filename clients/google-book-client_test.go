package client

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"

	model "example.com/book-learn/models"
	"github.com/stretchr/testify/assert"
)

func mockGetData(data []byte, err error) func(string) (*http.Response, error) {
	return func(url string) (*http.Response, error) {
		if err != nil {
			return nil, err
		}
		var res = http.Response{
			Status: "200 OK",
			Body:   io.NopCloser(bytes.NewReader(data)),
		}
		return &res, nil
	}
}

func TestGoogleBookClient_ByAuthor(t *testing.T) {
	type fields struct {
		GetData  func(url string) (resp *http.Response, err error)
		PactMode bool
	}
	type args struct {
		ctx     context.Context
		request GoogleBookRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.GoogleBookResponse
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				GetData: mockGetData([]byte(
					"{ \"Kind\": \"test-response\", \"TotalItems\": 0, \"Items\": [] }",
				), nil),
				PactMode: false,
			},
			args: args{
				ctx: context.Background(),
				request: GoogleBookRequest{
					Author: "test-author",
				},
			},
			want: model.GoogleBookResponse{
				Kind:       "test-response",
				TotalItems: 0,
				Items:      []model.GoogleBookItem{},
			},
		},
		{
			name: "failure",
			fields: fields{
				GetData:  mockGetData(nil, errors.New("test - author request fails")),
				PactMode: false,
			},
			args: args{
				ctx: context.Background(),
				request: GoogleBookRequest{
					Author: "test-author",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := GoogleBookClient{
				GetData:  tt.fields.GetData,
				PactMode: tt.fields.PactMode,
			}
			got, err := bc.ByAuthor(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleBookClient.ByAuthor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, got.Kind, tt.want.Kind)
			assert.Equal(t, got.Items, tt.want.Items)
			assert.Equal(t, got.TotalItems, tt.want.TotalItems)
		})
	}
}

func TestGoogleBookClient_ByTitle(t *testing.T) {
	type fields struct {
		GetData  func(url string) (resp *http.Response, err error)
		PactMode bool
	}
	type args struct {
		ctx     context.Context
		request GoogleBookRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    model.GoogleBookResponse
		wantErr bool
	}{
		{
			name: "makes a book request",
			fields: fields{
				GetData: mockGetData([]byte(
					"{ \"Kind\": \"test-reponse\", \"TotalItems\": 0, \"Items\": [] }",
				), nil),
				PactMode: false,
			},
			args: args{
				ctx: context.Background(),
				request: GoogleBookRequest{
					Author: "test-title",
				},
			},
			want: model.GoogleBookResponse{
				Kind:       "test-reponse",
				TotalItems: 0,
				Items:      []model.GoogleBookItem{},
			},
		},
		{
			name: "failure",
			fields: fields{
				GetData:  mockGetData(nil, errors.New("test - title request fails")),
				PactMode: false,
			},
			args: args{
				ctx: context.Background(),
				request: GoogleBookRequest{
					Author: "test-title",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bc := GoogleBookClient{
				GetData:  tt.fields.GetData,
				PactMode: tt.fields.PactMode,
			}
			got, err := bc.ByTitle(tt.args.ctx, tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("GoogleBookClient.ByTitle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleBookClient.ByTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterTitleResults(t *testing.T) {
	type args struct {
		req  GoogleBookRequest
		resp model.GoogleBookResponse
		err  error
	}
	tests := []struct {
		name string
		args args
		want struct {
			count int
			err   error
		}
	}{
		{
			name: "with error",
			args: args{
				req: GoogleBookRequest{},
				err: errors.New("test-error"),
			},
			want: struct {
				count int
				err   error
			}{
				count: 0,
				err:   errors.New("test-error"),
			},
		},
		{
			name: "with empty results",
			args: args{
				req: GoogleBookRequest{
					Title: "test-title",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   0,
					HasMorePages: false,
					Items:        []model.GoogleBookItem{},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 0,
				err:   nil,
			},
		},
		{
			name: "with results and nothing to filter",
			args: args{
				req: GoogleBookRequest{
					Title: "test-title",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   2,
					HasMorePages: false,
					Items: []model.GoogleBookItem{
						{
							Kind: "Book",
							ID:   "1",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title: "test-title",
							},
						},
						{
							Kind: "Book",
							ID:   "2",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title: "test-title",
							},
						},
					},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 2,
				err:   nil,
			},
		},
		{
			name: "with results and filtered title",
			args: args{
				req: GoogleBookRequest{
					Title: "test-title",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   2,
					HasMorePages: false,
					Items: []model.GoogleBookItem{
						{
							Kind: "Book",
							ID:   "1",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title: "NOT-test-title",
							},
						},
						{
							Kind: "Book",
							ID:   "2",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title: "test-title",
							},
						},
					},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 1,
				err:   nil,
			},
		},
		{
			name: "with results and filtered non-english",
			args: args{
				req: GoogleBookRequest{
					Title: "test-title",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   2,
					HasMorePages: false,
					Items: []model.GoogleBookItem{
						{
							Kind: "Book",
							ID:   "1",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "NOT-en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title: "test-title",
							},
						},
						{
							Kind: "Book",
							ID:   "2",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title: "test-title",
							},
						},
					},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 1,
				err:   nil,
			},
		},
		{
			name: "with results and filtered no-description",
			args: args{
				req: GoogleBookRequest{
					Title: "test-title",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   2,
					HasMorePages: false,
					Items: []model.GoogleBookItem{
						{
							Kind: "Book",
							ID:   "1",
							VolumeInfo: model.GoogleBookVolumeInfo{
								// Description: nil, -- no description
								Language: "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title: "test-title",
							},
						},
						{
							Kind: "Book",
							ID:   "2",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title: "test-title",
							},
						},
					},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 1,
				err:   nil,
			},
		},
		{
			name: "with results and filtered no-image",
			args: args{
				req: GoogleBookRequest{
					Title: "test-title",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   2,
					HasMorePages: false,
					Items: []model.GoogleBookItem{
						{
							Kind: "Book",
							ID:   "1",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								Title:       "test-title",
							},
						},
						{
							Kind: "Book",
							ID:   "2",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title: "test-title",
							},
						},
					},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 1,
				err:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterTitleResults(tt.args.req)(tt.args.resp, tt.args.err)
			if tt.want.err != nil && err == nil {
				t.Errorf("filterTitleResults() = expected error but got nil")
			} else {
				if !reflect.DeepEqual(len(got.Items), tt.want.count) {
					t.Errorf("filterTitleResults() = %v, want %v", len(got.Items), tt.want.count)
				}
			}
		})
	}
}

func Test_filterAuthorResults(t *testing.T) {
	type args struct {
		req  GoogleBookRequest
		resp model.GoogleBookResponse
		err  error
	}
	tests := []struct {
		name string
		args args
		want struct {
			count int
			err   error
		}
	}{
		{
			name: "with error",
			args: args{
				req: GoogleBookRequest{},
				err: errors.New("test-error"),
			},
			want: struct {
				count int
				err   error
			}{
				count: 0,
				err:   errors.New("test-error"),
			},
		},
		{
			name: "with empty results",
			args: args{
				req: GoogleBookRequest{
					Author: "Test-Author",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   0,
					HasMorePages: false,
					Items:        []model.GoogleBookItem{},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 0,
				err:   nil,
			},
		},
		{
			name: "with results and nothing to filter",
			args: args{
				req: GoogleBookRequest{
					Author: "Test-Author",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   2,
					HasMorePages: false,
					Items: []model.GoogleBookItem{
						{
							Kind: "Book",
							ID:   "1",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "http://example.com/has-thumbnail",
								},
								Title:   "test-title",
								Authors: []string{"Test-Author"},
							},
						},
						{
							Kind: "Book",
							ID:   "2",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "http://example.com/has-thumbnail",
								},
								Title:   "test-title",
								Authors: []string{"Test-Author"},
							},
						},
					},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 2,
				err:   nil,
			},
		},
		{
			name: "with results and filtered author",
			args: args{
				req: GoogleBookRequest{
					Author: "Test-Author",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   2,
					HasMorePages: false,
					Items: []model.GoogleBookItem{
						{
							Kind: "Book",
							ID:   "1",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title:   "test-title",
								Authors: []string{"not-the-author-you-are-looking-for"},
							},
						},
						{
							Kind: "Book",
							ID:   "2",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title:   "test-title",
								Authors: []string{"Test-Author"},
							},
						},
					},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 1,
				err:   nil,
			},
		},
		{
			name: "with results and filtered non-english",
			args: args{
				req: GoogleBookRequest{
					Author: "Test-Author",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   2,
					HasMorePages: false,
					Items: []model.GoogleBookItem{
						{
							Kind: "Book",
							ID:   "1",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "NOT-en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title:   "test-title",
								Authors: []string{"Test-Author"},
							},
						},
						{
							Kind: "Book",
							ID:   "2",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title:   "test-title",
								Authors: []string{"Test-Author"},
							},
						},
					},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 1,
				err:   nil,
			},
		},
		{
			name: "with results and filtered no-description",
			args: args{
				req: GoogleBookRequest{
					Author: "Test-Author",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   2,
					HasMorePages: false,
					Items: []model.GoogleBookItem{
						{
							Kind: "Book",
							ID:   "1",
							VolumeInfo: model.GoogleBookVolumeInfo{
								// Description: nil, -- no description
								Language: "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title:   "test-title",
								Authors: []string{"Test-Author"},
							},
						},
						{
							Kind: "Book",
							ID:   "2",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title:   "test-title",
								Authors: []string{"Test-Author"},
							},
						},
					},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 1,
				err:   nil,
			},
		},
		{
			name: "with results and filtered no-image",
			args: args{
				req: GoogleBookRequest{
					Author: "Test-Author",
				},
				resp: model.GoogleBookResponse{
					Kind:         "books#volumes",
					TotalItems:   2,
					HasMorePages: false,
					Items: []model.GoogleBookItem{
						{
							Kind: "Book",
							ID:   "1",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								Title:       "test-title",
								Authors:     []string{"Test-Author"},
							},
						},
						{
							Kind: "Book",
							ID:   "2",
							VolumeInfo: model.GoogleBookVolumeInfo{
								Description: "has-description",
								Language:    "en",
								ImageLinks: model.GoogleBookImageLinks{
									Thumbnail: "has-thumbnail",
								},
								Title:   "test-title",
								Authors: []string{"Test-Author"},
							},
						},
					},
				},
				err: nil,
			},
			want: struct {
				count int
				err   error
			}{
				count: 1,
				err:   nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := filterAuthorResults(tt.args.req)(tt.args.resp, tt.args.err)
			if tt.want.err != nil && err == nil {
				t.Errorf("filterResults() = expected error but got nil")
			} else {
				if !reflect.DeepEqual(len(got.Items), tt.want.count) {
					t.Errorf("filterResults() = %v, want %v", got, tt.want.count)
				}
			}
		})
	}
}

func Test_filterExactAuthor(t *testing.T) {
	type args struct {
		book model.GoogleBookItem
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with one element that matches",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{
						Authors: []string{"Testy Testerson"},
					},
				},
				name: "Testy Testerson",
			},
			want: true,
		},
		{
			name: "with one element that does not match",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{
						Authors: []string{"Testy Testerson"},
					},
				},
				name: "Not Testy Testerson",
			},
			want: false,
		},
		{
			name: "with two elements and one match",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{
						Authors: []string{"Testy Testerson", "Besty Besterson"},
					},
				},
				name: "Testy Testerson",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterExactAuthor(tt.args.book, tt.args.name); got != tt.want {
				t.Errorf("filterExactName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_buildRequestUrl(t *testing.T) {
	type args struct {
		query   string
		request GoogleBookRequest
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "with no query string",
			args: args{
				query: "",
				request: GoogleBookRequest{
					Title: "Test",
				},
			},
			want: "https://www.googleapis.com/books/v1/volumes?q=",
		},
		{
			name: "with a query string",
			args: args{
				query: "Testy-Testersons-Novel",
				request: GoogleBookRequest{
					Title: "Test",
				},
			},
			want: "https://www.googleapis.com/books/v1/volumes?q=Testy-Testersons-Novel",
		},
		{
			name: "with a start",
			args: args{
				query: "Testy-Testersons-Novel",
				request: GoogleBookRequest{
					Title: "Test",
					Start: 42,
				},
			},
			want: "https://www.googleapis.com/books/v1/volumes?q=Testy-Testersons-Novel&startIndex=42",
		},
		{
			name: "with a Limit",
			args: args{
				query: "Testy-Testersons-Novel",
				request: GoogleBookRequest{
					Title: "Test",
					Limit: 42,
				},
			},
			want: "https://www.googleapis.com/books/v1/volumes?q=Testy-Testersons-Novel&maxResults=42",
		},
		{
			name: "with a Start and a Limit",
			args: args{
				query: "Testy-Testersons-Novel",
				request: GoogleBookRequest{
					Title: "Test",
					Start: 24,
					Limit: 42,
				},
			},
			want: "https://www.googleapis.com/books/v1/volumes?q=Testy-Testersons-Novel&startIndex=24&maxResults=42",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildRequestUrl(tt.args.query, tt.args.request); got != tt.want {
				t.Errorf("buildRequestUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterExactTitle(t *testing.T) {
	type args struct {
		book model.GoogleBookItem
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "with blank book",
			args: args{
				book: model.GoogleBookItem{},
				name: "empty book",
			},
			want: false,
		},
		{
			name: "with close match on title",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{
						Title: "The Test Book",
					},
				},
				name: "The Test Book, Part 2",
			},
			want: false,
		},
		{
			name: "with exact match on title",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{
						Title: "The Test Book",
					},
				},
				name: "The Test Book",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterExactTitle(tt.args.book, tt.args.name); got != tt.want {
				t.Errorf("filterExactTitle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterIsEnglish(t *testing.T) {
	type args struct {
		book model.GoogleBookItem
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is English",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{
						Language: "en",
					},
				},
			},
			want: true,
		},
		{
			name: "is not English",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{
						Language: "it",
					},
				},
			},
			want: false,
		},
		{
			name: "no language",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterIsEnglish(tt.args.book); got != tt.want {
				t.Errorf("filterIsEnglish() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterHasDescription(t *testing.T) {
	type args struct {
		book model.GoogleBookItem
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "has a description",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{
						Description: "imadescription",
					},
				},
			},
			want: true,
		},
		{
			name: "no description",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterHasDescription(tt.args.book); got != tt.want {
				t.Errorf("filterHasDescription() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterHasImage(t *testing.T) {
	type args struct {
		book model.GoogleBookItem
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "has an image",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{
						ImageLinks: model.GoogleBookImageLinks{
							Thumbnail: "http://example.com/assets/image.jpg",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "doesn't have an image",
			args: args{
				book: model.GoogleBookItem{
					VolumeInfo: model.GoogleBookVolumeInfo{
						ImageLinks: model.GoogleBookImageLinks{},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterHasImage(tt.args.book); got != tt.want {
				t.Errorf("filterHasImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortByPublishedDate(t *testing.T) {
	book1 := model.GoogleBookItem{
		VolumeInfo: model.GoogleBookVolumeInfo{
			PublishedDate: "2023-01-01",
		},
	}
	book2 := model.GoogleBookItem{
		VolumeInfo: model.GoogleBookVolumeInfo{
			PublishedDate: "2022-01-01",
		},
	}
	book3 := model.GoogleBookItem{
		VolumeInfo: model.GoogleBookVolumeInfo{
			PublishedDate: "2021-01-01",
		},
	}

	type args struct {
		resp model.GoogleBookResponse
		err  error
	}

	tests := []struct {
		name    string
		args    args
		want    model.GoogleBookResponse
		wantErr bool
	}{
		{
			name: "sorts by published date",
			args: args{
				resp: model.GoogleBookResponse{
					Items: []model.GoogleBookItem{book2, book1, book3},
				},
				err: nil,
			},
			want: model.GoogleBookResponse{
				Items: []model.GoogleBookItem{book1, book2, book3},
			},
			wantErr: false,
		},
		{
			name: "with and error param",
			args: args{
				resp: model.GoogleBookResponse{},
				err:  errors.New("test-error"),
			},
			want:    model.GoogleBookResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sortByPublishedDate(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortByPublishedDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortByPublishedDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortByDescLength(t *testing.T) {
	book1 := model.GoogleBookItem{
		VolumeInfo: model.GoogleBookVolumeInfo{
			Description: "The very very longest description of all the books.",
		},
	}
	book2 := model.GoogleBookItem{
		VolumeInfo: model.GoogleBookVolumeInfo{
			Description: "A Description",
		},
	}
	book3 := model.GoogleBookItem{
		VolumeInfo: model.GoogleBookVolumeInfo{
			// No description at all. In fact it's nil !!!
		},
	}

	type args struct {
		resp model.GoogleBookResponse
		err  error
	}

	tests := []struct {
		name    string
		args    args
		want    model.GoogleBookResponse
		wantErr bool
	}{
		{
			name: "sorts by published date",
			args: args{
				resp: model.GoogleBookResponse{
					Items: []model.GoogleBookItem{book2, book1, book3},
				},
				err: nil,
			},
			want: model.GoogleBookResponse{
				Items: []model.GoogleBookItem{book1, book2, book3},
			},
			wantErr: false,
		},
		{
			name: "with and error param",
			args: args{
				resp: model.GoogleBookResponse{},
				err:  errors.New("test-error"),
			},
			want:    model.GoogleBookResponse{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sortByDescLength(tt.args.resp, tt.args.err)
			if (err != nil) != tt.wantErr {
				t.Errorf("sortByDescLength() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortByDescLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
