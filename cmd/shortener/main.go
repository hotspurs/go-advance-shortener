package main

import (
	"github.com/hotspurs/go-advance-shortener/internal/rand"
	"github.com/hotspurs/go-advance-shortener/internal/storage"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	data := storage.NewMemoryStorage(map[string]string{})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		MainHandler(w, r, data)
	})
	http.ListenAndServe("localhost:8080", mux)
}

func MainHandler(w http.ResponseWriter, r *http.Request, storage storage.Storage) {
	if r.URL.Path == "/" && r.Method == "POST" {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		short := rand.String(8)
		storage.Add(short, string(body))
		w.Header().Add("Content-Type", "text/plain")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + short))
		return
	}

	re, err := regexp.MatchString("^/([a-zA-Z]+$)", r.URL.Path)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if re && r.Method == "GET" {
		short := strings.TrimPrefix(r.URL.Path, "/")
		w.Header().Add("Location", storage.Get(short))
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
}
