package config

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
)

var (
	DefaultAddress       = "127.0.0.1:8080"
	DefaultStoreInterval = 30 * time.Second
	DefaultRestore       = true
	DefaultStoreFile     = "/tmp/devops-metrics-db.json"
	DefaultKey           string
	DefaultCompressLevel = 5
	DefaultCompressTypes = []string{
		"text/html",
		"text/plain",
		"application/json",
	}
	DefaultDSN        string
	DefaultMigrateURL = "file://services/server/migrations"
)

type (
	Config struct {
		HTTP          HTTPConfig
		MemoryStorage MemoryStorageConfig
		Postgres      PostgresConfig
		SecretKey     string
	}

	HTTPConfig struct {
		Address       string `env:"ADDRESS"`
		CompressLevel int
		CompressTypes []string
	}

	MemoryStorageConfig struct {
		StoreInterval time.Duration `env:"STORE_INTERVAL"`
		StoreFile     string        `env:"STORE_FILE"`
		Restore       bool          `env:"RESTORE"`
	}

	PostgresConfig struct {
		DSN        string `env:"DATABASE_DSN"`
		MigrateURL string
	}
)

func New() (*Config, error) {
	cfg := Config{}
	cfg.HTTP.Address = DefaultAddress
	cfg.HTTP.CompressLevel = DefaultCompressLevel
	cfg.HTTP.CompressTypes = DefaultCompressTypes
	cfg.MemoryStorage.StoreInterval = DefaultStoreInterval
	cfg.MemoryStorage.Restore = DefaultRestore
	cfg.MemoryStorage.StoreFile = DefaultStoreFile
	cfg.Postgres.DSN = DefaultDSN
	cfg.Postgres.MigrateURL = DefaultMigrateURL
	cfg.SecretKey = DefaultKey

	if err := env.Parse(&cfg); err != nil {
		return nil, NewConfigError(err)
	}

	_, port, err := net.SplitHostPort(cfg.HTTP.Address)
	if err != nil {
		return nil, NewConfigError(fmt.Errorf("invalid address %s", cfg.HTTP.Address))
	}

	_, err = strconv.Atoi(port)
	if err != nil {
		return nil, NewConfigError(fmt.Errorf("invalid port %s", port))
	}

	return &cfg, nil
}
