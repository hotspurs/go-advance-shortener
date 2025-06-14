package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hotspurs/go-advance-shortener/internal/config"
	"github.com/hotspurs/go-advance-shortener/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getStorage() *storage.FileStorage {
	store, _ := storage.NewFileStorage("./storage.json")

	return store
}

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
		data    Storage
		name    string
		want    want
		request request
	}{
		{
			name: "GeneratePositive",
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
			data: getStorage(),
		},
		{
			name: "GenerateNegative_BadBody",
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   errReader{},
			},
			want: want{
				code:        http.StatusInternalServerError,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			data: getStorage(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, test.request.url, test.request.body)
			w := httptest.NewRecorder()
			GenerateHandler(test.data, cfg)(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode, "expected status code %d, got %d", test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err, "unexpected error reading response body: %v", err)
			assert.Contains(t, string(resBody), test.want.response, "expected response to contain %q, got %q", test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"), "expected content type %q, got %q", test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestShortenHandler(t *testing.T) {
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
		data    Storage
		name    string
		want    want
		request request
	}{
		{
			name: "ShortenPositive",
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   bytes.NewReader([]byte("{\"url\": \"https://ya.ru\"}")),
			},
			want: want{
				code:        http.StatusCreated,
				response:    cfg.BaseURL,
				contentType: "application/json",
			},
			data: getStorage(),
		},
		{
			name: "ShortenNegative_BadBody",
			request: request{
				method: http.MethodPost,
				url:    "/",
				body:   errReader{},
			},
			want: want{
				code:        http.StatusBadRequest,
				response:    "",
				contentType: "text/plain; charset=utf-8",
			},
			data: getStorage(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, test.request.url, test.request.body)
			w := httptest.NewRecorder()
			ShortenHandler(test.data, cfg)(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode, "expected status code %d, got %d", test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err, "unexpected error reading response body: %v", err)
			assert.Contains(t, string(resBody), test.want.response, "expected response to contain %q, got %q", test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"), "expected content type %q, got %q", test.want.contentType, res.Header.Get("Content-Type"))
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
		data    Storage
		name    string
		want    want
		request request
	}{
		{
			name: "GetPositive",
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
			data: getStorage(),
		},
		{
			name: "GetNegative_UnknownKey",
			request: request{
				method: http.MethodGet,
				url:    "/unknownkey",
				body:   bytes.NewReader([]byte("")),
			},
			want: want{
				code:        http.StatusTemporaryRedirect,
				response:    "",
				contentType: "",
			},
			data: getStorage(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, test.request.url, test.request.body)
			w := httptest.NewRecorder()
			GetHandler(test.data)(w, request)

			res := w.Result()
			assert.Equal(t, test.want.code, res.StatusCode, "expected status code %d, got %d", test.want.code, res.StatusCode)
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err, "unexpected error reading response body: %v", err)
			assert.Contains(t, string(resBody), test.want.response, "expected response to contain %q, got %q", test.want.response, string(resBody))
			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"), "expected content type %q, got %q", test.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

type errReader struct{}

func (errReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}
