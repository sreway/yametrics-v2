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
	DefaultCryptoKey     string
	DefaultCryptoCrt     string
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
		SecretKey     string
		Postgres      PostgresConfig
		MemoryStorage MemoryStorageConfig
		HTTP          HTTPConfig
	}

	HTTPConfig struct {
		Address       string `env:"ADDRESS"`
		CryptoKey     string `env:"CRYPTO_KEY"`
		CryptoCrt     string `env:"CRYPTO_CRT"`
		CompressTypes []string
		CompressLevel int
	}

	MemoryStorageConfig struct {
		StoreFile     string        `env:"STORE_FILE"`
		Restore       bool          `env:"RESTORE"`
		StoreInterval time.Duration `env:"STORE_INTERVAL"`
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
	cfg.HTTP.CryptoKey = DefaultCryptoKey
	cfg.HTTP.CryptoCrt = DefaultCryptoCrt
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
