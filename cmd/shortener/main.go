package main

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/hotspurs/go-advance-shortener/internal/compress"
	"github.com/hotspurs/go-advance-shortener/internal/config"
	"github.com/hotspurs/go-advance-shortener/internal/handlers"
	logger "github.com/hotspurs/go-advance-shortener/internal/logger"
	"github.com/hotspurs/go-advance-shortener/internal/storage"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func main() {
	cfg := config.Init()
	r := chi.NewRouter()
	data, err := storage.NewFileStorage(cfg.FileStoragePath)

	if err != nil {
		panic(err)
	}

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.Debug)
	sugar := log.Sugar
	defer log.Sync()

	sugar.Infof("Initialize")

	r.Method("POST", "/", compress.WithGzip(logger.WithLogging(handlers.GenerateHandler(data, cfg), log)))
	r.Method("POST", "/api/shorten", compress.WithGzip(logger.WithLogging(handlers.ShortenHandler(data, cfg), log)))
	r.Method("GET", "/{link}", logger.WithLogging(handlers.GetHandler(data), log))
	r.Method("GET", "/ping", handlers.PingHandler(db))

	sugar.Infof("Server is listen on port %s", cfg.Address)
	http.ListenAndServe(cfg.Address, r)
}
