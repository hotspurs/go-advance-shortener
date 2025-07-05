package handlers

import (
	"bytes"
	"github.com/hotspurs/go-advance-shortener/internal/config"
	"github.com/hotspurs/go-advance-shortener/internal/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBatchHandler(t *testing.T) {
	cfg := config.Init()
	log := logger.New(cfg.Debug)
	defer log.Sync()

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
			name: "BatchPositive",
			request: request{
				method: http.MethodPost,
				url:    "/api/shorten/batch",
				body:   bytes.NewReader([]byte("[\n    {\n        \"correlation_id\": \"1\",\n        \"original_url\": \"https://google.com\"\n    },\n    {\n        \"correlation_id\": \"2\",\n        \"original_url\": \"https://openai.com\"\n    }\n]")),
			},
			want: want{
				code:        http.StatusCreated,
				response:    cfg.BaseURL,
				contentType: "application/json",
			},
			data: getStorage(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.request.method, test.request.url, test.request.body)
			w := httptest.NewRecorder()
			BatchHandler(test.data, cfg, log)(w, request)

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
