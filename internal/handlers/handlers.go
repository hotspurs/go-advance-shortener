package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/hotspurs/go-advance-shortener/internal/compress"
	"github.com/hotspurs/go-advance-shortener/internal/config"
	"github.com/hotspurs/go-advance-shortener/internal/rand"
	"net/http"
	"strings"
)

type Storage interface {
	Add(key string, value string)
	Get(key string) string
}

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func GetHandler(data Storage) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		short := strings.TrimPrefix(r.URL.Path, "/")
		w.Header().Add("Location", data.Get(short))
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
}

func GenerateHandler(data Storage, config *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var encoding = r.Header.Get("Content-Encoding")
		var buf bytes.Buffer
		var body []byte
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if encoding == "gzip" {
			body, err = compress.Decompress(buf.Bytes())
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			body = buf.Bytes()
		}

		short := rand.String(8)
		data.Add(short, string(body))
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(config.BaseURL + "/" + short))
	})
}

func ShortenHandler(data Storage, config *config.Config) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var encoding = r.Header.Get("Content-Encoding")
		var buf bytes.Buffer
		var body []byte
		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if encoding == "gzip" {
			body, err = compress.Decompress(buf.Bytes())
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		} else {
			body = buf.Bytes()
		}

		var req Request

		if err = json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		short := rand.String(8)
		data.Add(short, req.URL)
		var res Response
		res.Result = config.BaseURL + "/" + short
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		resp, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(resp)
	})
}
