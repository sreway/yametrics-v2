package config

import (
	"fmt"
	"time"
)

type OptionAgent func(config *Config) error

func WithPollInterval(pollInterval string) OptionAgent {
	return func(cfg *Config) error {
		pollIntervalDuration, err := time.ParseDuration(pollInterval)
		if err != nil {
			return NewConfigError(fmt.Errorf("invalid poll interval %s", pollInterval))
		}
		cfg.PollInterval = pollIntervalDuration
		return nil
	}
}

func WithReportInterval(reportInterval string) OptionAgent {
	return func(cfg *Config) error {
		reportIntervalDuration, err := time.ParseDuration(reportInterval)
		if err != nil {
			return NewConfigError(fmt.Errorf("invalid poll interval %s", reportInterval))
		}

		cfg.ReportInterval = reportIntervalDuration
		return nil
	}
}
