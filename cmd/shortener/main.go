package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/hotspurs/go-advance-shortener/internal/compress"
	"github.com/hotspurs/go-advance-shortener/internal/config"
	"github.com/hotspurs/go-advance-shortener/internal/handlers"
	logger "github.com/hotspurs/go-advance-shortener/internal/logger"
	"github.com/hotspurs/go-advance-shortener/internal/storage"
	"net/http"
)

func main() {
	cfg := config.Init()
	r := chi.NewRouter()
	data := storage.NewMemoryStorage(map[string]string{})
	log := logger.New(cfg.Debug)
	sugar := log.Sugar
	defer log.Sync()

	sugar.Infof("Initialize")

	r.Method("POST", "/", logger.WithLogging(handlers.GenerateHandler(data, cfg), log))
	r.Method("POST", "/api/shorten", compress.WithGzip(logger.WithLogging(handlers.ShortenHandler(data, cfg), log)))
	r.Method("GET", "/{link}", compress.WithGzip(logger.WithLogging(handlers.GetHandler(data), log)))

	sugar.Infof("Server is listen on port %s", cfg.Address)
	http.ListenAndServe(cfg.Address, r)
}
