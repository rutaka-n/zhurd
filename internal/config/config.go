package config

import (
	"encoding/json"
	"io"
)

// Logger contains all configuration for logger.
type Logger struct {
	Destination string `json:"destination"`
	Level       string `json:"level"`
}

// Server contains all configuration for server.
type Server struct {
	Addr               string `json:"addr"`
	Logger             Logger `json:"logger"`
	GracefulTimeoutSec int    `json:"graceful_timeout_s"`
}

// Config is a high-level struct that contains all configuration.
type Config struct {
	Server Server `json:"server"`
}

// Load configuration from file.
func Load(cfgData io.Reader) (Config, error) {
	cfg := Config{}

	data, err := io.ReadAll(cfgData)
	if err != nil {
		return Config{}, err
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
