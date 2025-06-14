package config

import (
	"flag"
	"os"
	"sync"
)

type Config struct {
	Address         string
	BaseURL         string
	FileStoragePath string
	Debug           bool
	DatabaseDSN     string
}

var (
	cfg  *Config
	once sync.Once
)

func Init() *Config {
	once.Do(func() {
		address := flag.String("a", "localhost:8080", "HTTP server address")
		baseURL := flag.String("b", "http://localhost:8080", "Base URL")
		fileStoragePath := flag.String("f", "./storage.json", "File storage path")
		databaseDSN := flag.String("d", "", "Database DSN")

		envAddress := os.Getenv("SERVER_ADDRESS")
		envBaseURL := os.Getenv("BASE_URL")
		envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
		envDatabaseDSN := os.Getenv("DATABASE_DSN")

		envDebug := os.Getenv("DEBUG")

		var debug bool
		if envDebug != "" {
			debug = true
		}

		flag.Parse()

		if envAddress != "" {
			address = &envAddress
		}

		if envBaseURL != "" {
			baseURL = &envBaseURL
		}

		if envFileStoragePath != "" {
			fileStoragePath = &envFileStoragePath
		}

		if envDatabaseDSN != "" {
			databaseDSN = &envDatabaseDSN
		}

		cfg = &Config{
			Address:         *address,
			BaseURL:         *baseURL,
			FileStoragePath: *fileStoragePath,
			Debug:           debug,
			DatabaseDSN:     *databaseDSN,
		}
	})

	return cfg
}
