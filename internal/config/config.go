package config

import (
	"encoding/json"
	"fmt"
	"io"
)

// Logger contains all configuration for logger.
type Logger struct {
	Destination string `json:"destination"`
	Level       string `json:"level"`
	Format      string `json:"format"`
}

// Server contains all configuration for server.
type Server struct {
	Addr               string `json:"addr"`
	Logger             Logger `json:"logger"`
	GracefulTimeoutSec int    `json:"graceful_timeout_s"`
	QueueBufferSize    int    `json:"queue_buffer_size"`
}

// Databse contains all configuration for database connection.
type Database struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	SSLMode  string `json:"ssl_mode"`
}

func (cfg *Database) ConnectionString() string {
	if cfg.Host == "" ||
		cfg.Port == 0 ||
		cfg.User == "" ||
		cfg.Password == "" ||
		cfg.Name == "" ||
		cfg.SSLMode == "" {
		return ""
	}
	return fmt.Sprintf(
		"host=%s port=%v user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)
}

// Config is a high-level struct that contains all configuration.
type Config struct {
	Server   Server   `json:"server"`
	Database Database `json:"database"`
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
