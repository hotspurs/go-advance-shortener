package main

import (
	"net/http"
	"regexp"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" && r.Method == "POST" {
			w.Header().Add("Content-Type", "text/plain")
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("http://localhost:8080/EwHXdJfB"))
			return
		}

		re, err := regexp.MatchString("^/([a-zA-Z]+$)", r.URL.Path)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if re && r.Method == "GET" {
			w.Header().Add("Location", "https://practicum.yandex.ru/")
			w.WriteHeader(http.StatusTemporaryRedirect)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	})
	http.ListenAndServe("localhost:8080", mux)
}
