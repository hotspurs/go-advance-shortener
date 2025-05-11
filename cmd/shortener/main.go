package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hotspurs/go-advance-shortener/internal/config"
	"github.com/hotspurs/go-advance-shortener/internal/rand"
	"github.com/hotspurs/go-advance-shortener/internal/storage"
)

type Storage interface {
	Add(key string, value string)
	Get(key string) string
}

func main() {
	cfg := config.Init()
	r := chi.NewRouter()
	data := storage.NewMemoryStorage(map[string]string{})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		GenerateHandler(w, r, data, cfg)
	})

	r.Get("/{link}", func(w http.ResponseWriter, r *http.Request) {
		GetHandler(w, r, data)
	})
	http.ListenAndServe(cfg.Address, r)
}

func GenerateHandler(w http.ResponseWriter, r *http.Request, data Storage, config *config.Config) {
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
}

func GetHandler(w http.ResponseWriter, r *http.Request, data Storage) {
	short := strings.TrimPrefix(r.URL.Path, "/")
	w.Header().Add("Location", data.Get(short))
	w.WriteHeader(http.StatusTemporaryRedirect)
}
