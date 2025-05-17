package handlers

import (
	"github.com/hotspurs/go-advance-shortener/internal/config"
	"github.com/hotspurs/go-advance-shortener/internal/rand"
	"io"
	"net/http"
	"strings"
)

type Storage interface {
	Add(key string, value string)
	Get(key string) string
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
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		short := rand.String(8)
		data.Add(short, string(body))
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(config.BaseURL + "/" + short))
	})
}
