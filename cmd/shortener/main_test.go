package main

import (
	"bytes"
	"github.com/hotspurs/go-advance-shortener/internal/config"
	"github.com/hotspurs/go-advance-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateHandler(t *testing.T) {
	cfg := config.Init()
	type request struct {
		method string
		body   io.Reader
		url    string
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		data    storage.Storage
		name    string
		want    want
		request request
	}{
		{
			name: "generate positive",
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   bytes.NewReader([]byte("https://ya.ru")),
			},
			want: want{
				code:        http.StatusCreated,
				response:    cfg.BaseURL,
				contentType: "text/plain",
			},
			data: storage.NewMemoryStorage(map[string]string{}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, test.request.url, test.request.body)
			w := httptest.NewRecorder()
			GenerateHandler(w, request, test.data, cfg)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Contains(t, string(resBody), test.want.response)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestGetHandler(t *testing.T) {
	type request struct {
		method string
		body   io.Reader
		url    string
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		data    storage.Storage
		name    string
		want    want
		request request
	}{
		{
			name: "get positive",
			request: request{
				method: http.MethodGet,
				url:    "/tdluNOuy",
				body:   bytes.NewReader([]byte("")),
			},
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "",
				contentType: "",
			},
			data: storage.NewMemoryStorage(map[string]string{
				"tdluNOuy": "https://ya.ru",
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, test.request.url, test.request.body)
			w := httptest.NewRecorder()
			GetHandler(w, request, test.data)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Contains(t, string(resBody), test.want.response)
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
