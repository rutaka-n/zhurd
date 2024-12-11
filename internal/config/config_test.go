package config

import (
	"bytes"
	"testing"
)

func TestLoad(t *testing.T) {
	cfgData := bytes.NewBufferString(`
{
  "server": {
    "addr": "localhost:3003",
    "graceful_timeout_s": 5,
    "queue_buffer_size": 64,
    "logger": {
      "destination": "stdout",
      "format": "json",
      "level": "info"
    }
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "zhurd",
    "password": "passwordsecretdb",
    "name": "zhurd",
    "ssl_mode": "disable"
  }
}
    `)

	cfg, err := Load(cfgData)
	if err != nil {
		t.Fatalf("expected no errors, got %v\n", err)
	}
	if cfg.Server.Addr != "localhost:3003" {
		t.Errorf("expected %s, got %s\n", "localhost:3003", cfg.Server.Addr)
	}
	if cfg.Server.QueueBufferSize != 64 {
		t.Errorf("expected %d, got %d\n", 64, cfg.Server.QueueBufferSize)
	}
	if cfg.Server.GracefulTimeoutSec != 5 {
		t.Errorf("expected %d, got %d\n", 5, cfg.Server.GracefulTimeoutSec)
	}
	if cfg.Server.Logger.Destination != "stdout" {
		t.Errorf("expected %s, got %s\n", "stdout", cfg.Server.Logger.Destination)
	}
	if cfg.Server.Logger.Format != "json" {
		t.Errorf("expected %s, got %s\n", "json", cfg.Server.Logger.Format)
	}
	if cfg.Server.Logger.Level != "info" {
		t.Errorf("expected %s, got %s\n", "info", cfg.Server.Logger.Level)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected %s, got %s\n", "localhost", cfg.Database.Host)
	}
	if cfg.Database.Port != 5432 {
		t.Errorf("expected %d, got %d\n", 5432, cfg.Database.Port)
	}
	if cfg.Database.User != "zhurd" {
		t.Errorf("expected %s, got %s\n", "zhurd", cfg.Database.User)
	}
	if cfg.Database.Password != "passwordsecretdb" {
		t.Errorf("expected %s, got %s\n", "localhost", cfg.Database.Password)
	}
	if cfg.Database.Name != "zhurd" {
		t.Errorf("expected %s, got %s\n", "zhurd", cfg.Database.Name)
	}
	if cfg.Database.SSLMode != "disable" {
		t.Errorf("expected %s, got %s\n", "disable", cfg.Database.SSLMode)
	}
	if cfg.Database.ConnectionString() != "host=localhost port=5432 user=zhurd password=passwordsecretdb dbname=zhurd sslmode=disable" {
		t.Errorf("expected %s, got %s\n",
			"host=localhost port=5432 user=zhurd password=passwordsecretdb dbname=zhurd sslmode=disable",
			cfg.Database.ConnectionString(),
		)
	}
}
