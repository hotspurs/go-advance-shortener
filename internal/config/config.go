package config

import (
	"flag"
	"os"
)

type Config struct {
	Address string
	BaseURL string
}

func Init() *Config {
	address := flag.String("a", "localhost:8080", "HTTP server address")
	baseURL := flag.String("b", "http://localhost:8080", "Base URL")

	envAddress := os.Getenv("SERVER_ADDRESS")
	envBaseUrl := os.Getenv("BASE_URL")

	flag.Parse()

	if envAddress != "" {
		address = &envAddress
	}

	if envBaseUrl != "" {
		address = &envBaseUrl
	}

	cfg := &Config{
		Address: *address,
		BaseURL: *baseURL,
	}

	return cfg
}
