package config

import (
	"fmt"
	"net"
	"strconv"
	"time"

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
)

type (
	Config struct {
		MetricsEnpoint   string
		ServerHTTPScheme string
		Key              string        `env:"KEY"`
		ServerAddress    string        `env:"ADDRESS"`
		PollInterval     time.Duration `env:"POLL_INTERVAL"`
		ReportInterval   time.Duration `env:"REPORT_INTERVAL"`
		ServerPublicKey  string        `env:"CRYPTO_KEY"`
	}
)

func New() (*Config, error) {
	cfg := Config{
		ServerAddress:    DefaultServerAddress,
		ReportInterval:   DefaultReportInterval,
		PollInterval:     DefaultPollInterval,
		Key:              DefaultSecretKey,
		MetricsEnpoint:   DefaultMetricsEnpoint,
		ServerHTTPScheme: DefaultServerHTTPScheme,
		ServerPublicKey:  DefaultServerPublicKey,
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
