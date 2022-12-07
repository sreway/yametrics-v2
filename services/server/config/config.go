package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	log "github.com/sreway/yametrics-v2/pkg/tools/logger"

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
	DefaultConfigFile    string
	DefaultDSN           string
	DefaultMigrateURL    = "file://services/server/migrations"
	DefaultTrustedSubnet string
	DefaultUseGRPC       = false
)

type (
	Config struct {
		ConfigFile    string `env:"CONFIG"`
		SecretKey     string
		TrustedSubnet string              `json:"trusted_subnet" env:"TRUSTED_SUBNET"`
		Postgres      PostgresConfig      `json:"postgres"`
		MemoryStorage MemoryStorageConfig `json:"memory_storage"`
		Delivery      DeliveryConfig      `json:"delivery"`
		UseGRPC       bool                `json:"use_grpc" env:"USE_GRPC"`
	}

	DeliveryConfig struct {
		TrustedSubnet *net.IPNet
		Address       string `json:"address" env:"ADDRESS"`
		CryptoKey     string `json:"crypto_key" env:"CRYPTO_KEY"`
		CryptoCrt     string `json:"crypto_cert" env:"CRYPTO_CRT"`
		CompressTypes []string
		CompressLevel int
	}

	MemoryStorageConfig struct {
		StoreFile     string        `json:"store_file" env:"STORE_FILE"`
		Restore       bool          `json:"restore" env:"RESTORE"`
		StoreInterval time.Duration `json:"store_interval" env:"STORE_INTERVAL"`
	}

	PostgresConfig struct {
		DSN        string `json:"database_dsn" env:"DATABASE_DSN"`
		MigrateURL string
	}
)

func New() (*Config, error) {
	cfg := Config{}
	cfg.Delivery.Address = DefaultAddress
	cfg.Delivery.CompressLevel = DefaultCompressLevel
	cfg.Delivery.CompressTypes = DefaultCompressTypes
	cfg.Delivery.CryptoKey = DefaultCryptoKey
	cfg.Delivery.CryptoCrt = DefaultCryptoCrt
	cfg.TrustedSubnet = DefaultTrustedSubnet
	cfg.MemoryStorage.StoreInterval = DefaultStoreInterval
	cfg.MemoryStorage.Restore = DefaultRestore
	cfg.MemoryStorage.StoreFile = DefaultStoreFile
	cfg.Postgres.DSN = DefaultDSN
	cfg.Postgres.MigrateURL = DefaultMigrateURL
	cfg.SecretKey = DefaultKey
	cfg.ConfigFile = DefaultConfigFile
	cfg.UseGRPC = DefaultUseGRPC

	if cfg.ConfigFile != "" {
		f, err := os.Open(cfg.ConfigFile)
		if err != nil {
			return nil, NewConfigError(err)
		}
		defer f.Close()

		if err = json.NewDecoder(f).Decode(&cfg); err != nil {
			return nil, NewConfigError(err)
		}
		log.Info("success load json config")
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, NewConfigError(err)
	}

	_, port, err := net.SplitHostPort(cfg.Delivery.Address)
	if err != nil {
		return nil, NewConfigError(fmt.Errorf("invalid address %s", cfg.Delivery.Address))
	}

	_, err = strconv.Atoi(port)
	if err != nil {
		return nil, NewConfigError(fmt.Errorf("invalid port %s", port))
	}

	if cfg.TrustedSubnet != "" {
		if _, cfg.Delivery.TrustedSubnet, err = net.ParseCIDR(cfg.TrustedSubnet); err != nil {
			return nil, NewConfigError(fmt.Errorf("incorrect trusted subnet %s", cfg.TrustedSubnet))
		}
	}
	return &cfg, nil
}
