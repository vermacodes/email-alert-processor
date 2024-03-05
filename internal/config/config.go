package config

import (
	"errors"
	"os"
)

type Config struct {
	APIKey string
	Port   string
}

func NewConfig() (*Config, error) {

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, errors.New("API_KEY must be set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}

	return &Config{
		APIKey: apiKey,
		Port:   port,
	}, nil
}
