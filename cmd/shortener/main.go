package main

import (
	"fmt"
	"github.com/hotspurs/go-advance-shortener/internal/rand"
	"io"
	"net/http"
	"regexp"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	storage := make(map[string]string)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			short := rand.String(8)
			storage[short] = string(body)
			fmt.Println(storage)
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
			w.Header().Add("Location", storage[short])
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	})
	http.ListenAndServe("localhost:8080", mux)
}
