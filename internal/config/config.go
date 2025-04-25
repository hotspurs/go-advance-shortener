package config

import (
	"flag"
)

type Config struct {
	Address string
	BaseURL string
}

func Init() *Config {
	address := flag.String("a", "localhost:8080", "HTTP server address")
	baseURL := flag.String("b", "http://localhost:8080", "Base URL")

	flag.Parse()

	cfg := &Config{
		Address: *address,
		BaseURL: *baseURL,
	}

	return cfg
}
