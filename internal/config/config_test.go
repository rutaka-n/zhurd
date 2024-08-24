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
            "logger": {
                "destination": "stdout",
                "level": "INFO"
            }
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
	if cfg.Server.Logger.Destination != "stdout" {
		t.Errorf("expected %s, got %s\n", "stdout", cfg.Server.Logger.Destination)
	}
	if cfg.Server.Logger.Level != "INFO" {
		t.Errorf("expected %s, got %s\n", "INFO", cfg.Server.Logger.Destination)
	}
}
