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
					"{ \"Kind\": \"test-reponse\", \"TotalItems\": 0, \"Items\": [] }",
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
				Kind:       "test-reponse",
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
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoogleBookClient.ByAuthor() = %v, want %v", got, tt.want)
			}
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
