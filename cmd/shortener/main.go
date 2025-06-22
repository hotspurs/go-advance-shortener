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
	var data handlers.Storage
	var err error

	if cfg.DatabaseDSN != "" {
		db, err := sql.Open("pgx", cfg.DatabaseDSN)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		data = storage.NewDatabaseStorage(db)
		r.Method("GET", "/ping", handlers.PingHandler(db))
	} else if cfg.FileStoragePath != "" {
		data, err = storage.NewFileStorage(cfg.FileStoragePath)

		if err != nil {
			panic(err)
		}
	} else {
		data = storage.NewMemoryStorage(map[string]string{})
	}

	log := logger.New(cfg.Debug)
	sugar := log.Sugar
	defer log.Sync()

	sugar.Infof("Initialize")

	r.Method("POST", "/", compress.WithGzip(logger.WithLogging(handlers.GenerateHandler(data, cfg), log)))
	r.Method("POST", "/api/shorten", compress.WithGzip(logger.WithLogging(handlers.ShortenHandler(data, cfg), log)))
	r.Method("POST", "/api/shorten/batch", logger.WithLogging(handlers.BatchHandler(data, cfg, log), log))
	r.Method("GET", "/{link}", logger.WithLogging(handlers.GetHandler(data), log))

	sugar.Infof("Server is listen on port %s", cfg.Address)
	http.ListenAndServe(cfg.Address, r)
}
