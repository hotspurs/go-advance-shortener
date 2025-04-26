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
	envBaseURL := os.Getenv("BASE_URL")

	flag.Parse()

	if envAddress != "" {
		address = &envAddress
	}

	if envBaseURL != "" {
		baseURL = &envBaseURL
	}

	cfg := &Config{
		Address: *address,
		BaseURL: *baseURL,
	}

	return cfg
}
