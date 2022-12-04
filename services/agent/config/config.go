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
	DefaultServerAddress    = "127.0.0.1:8080"
	DefaultServerHTTPScheme = "http"
	DefaultReportInterval   = 10 * time.Second
	DefaultPollInterval     = 2 * time.Second
	DefaultSecretKey        = "secret"
	DefaultMetricsEnpoint   = "/updates/"
	DefaultServerPublicKey  string
	DefaultConfigFile       string
	DefaultRealIP           = "127.0.0.1"
)

type (
	Config struct {
		MetricsEnpoint   string
		ServerHTTPScheme string
		ConfigFile       string        `env:"CONFIG"`
		Key              string        `env:"KEY"`
		ServerAddress    string        `json:"address" env:"ADDRESS"`
		ServerPublicKey  string        `json:"crypto_key" env:"CRYPTO_KEY"`
		RealIP           string        `json:"real_ip" env:"REAL_IP"`
		PollInterval     time.Duration `json:"poll_interval" env:"POLL_INTERVAL"`
		ReportInterval   time.Duration `json:"report_interval" env:"REPORT_INTERVAL"`
	}
)

func New() (*Config, error) {
	cfg := Config{
		ReportInterval:   DefaultReportInterval,
		PollInterval:     DefaultPollInterval,
		ConfigFile:       DefaultConfigFile,
		ServerAddress:    DefaultServerAddress,
		Key:              DefaultSecretKey,
		MetricsEnpoint:   DefaultMetricsEnpoint,
		ServerHTTPScheme: DefaultServerHTTPScheme,
		ServerPublicKey:  DefaultServerPublicKey,
		RealIP:           DefaultRealIP,
	}

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

	_, port, err := net.SplitHostPort(cfg.ServerAddress)
	if err != nil {
		return nil, NewConfigError(fmt.Errorf("invalid address %s", cfg.ServerAddress))
	}

	_, err = strconv.Atoi(port)
	if err != nil {
		return nil, NewConfigError(fmt.Errorf("invalid port %s", port))
	}

	if cfg.ServerPublicKey != "" {
		cfg.ServerHTTPScheme = "https"
	}

	return &cfg, nil
}
