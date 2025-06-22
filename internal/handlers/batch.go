package handlers

import (
	"encoding/json"
	"github.com/hotspurs/go-advance-shortener/internal/logger"
	"github.com/hotspurs/go-advance-shortener/internal/rand"
	"io"
	"net/http"

	"github.com/hotspurs/go-advance-shortener/internal/config"
)

type RequestBatchItem struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatchItem struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func BatchHandler(data Storage, config *config.Config, logger *logger.Logger) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		var req []RequestBatchItem

		if err = json.Unmarshal(body, &req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		var res []ResponseBatchItem
		var urls []string
		var shorts []string
		for _, item := range req {
			short := rand.String(8)
			shorts = append(shorts, short)
			urls = append(urls, item.OriginalURL)
			res = append(res, ResponseBatchItem{CorrelationID: item.CorrelationID, ShortURL: config.BaseURL + "/" + short})
		}

		err = data.AddBatch(urls, shorts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

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
