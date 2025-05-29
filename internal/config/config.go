package config

import (
	"flag"
	"os"
	"sync"
)

type Config struct {
	Address string
	BaseURL string
	Debug   bool
}

var (
	cfg  *Config
	once sync.Once
)

func Init() *Config {
	once.Do(func() {
		address := flag.String("a", "localhost:8080", "HTTP server address")
		baseURL := flag.String("b", "http://localhost:8080", "Base URL")

		envAddress := os.Getenv("SERVER_ADDRESS")
		envBaseURL := os.Getenv("BASE_URL")

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

		cfg = &Config{
			Address: *address,
			BaseURL: *baseURL,
			Debug:   debug,
		}
	})

	return cfg
}
