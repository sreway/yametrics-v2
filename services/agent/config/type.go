package config

import (
	"encoding/json"
	"fmt"
	"time"
)

func (c *Config) UnmarshalJSON(data []byte) error {
	var err error
	type AliasType Config
	aliasValue := &struct {
		*AliasType
		PollInterval   string `json:"poll_interval"`
		ReportInterval string `json:"report_interval"`
	}{
		AliasType: (*AliasType)(c),
	}

	if err = json.Unmarshal(data, aliasValue); err != nil {
		return err
	}

	if aliasValue.PollInterval != "" {
		if c.PollInterval, err = time.ParseDuration(aliasValue.PollInterval); err != nil {
			return fmt.Errorf("failed to parse '%s' to time.Duration: %w", aliasValue.PollInterval, err)
		}
	}

	if aliasValue.ReportInterval != "" {
		if c.ReportInterval, err = time.ParseDuration(aliasValue.ReportInterval); err != nil {
			return fmt.Errorf("failed to parse '%s' to time.Duration: %w", aliasValue.ReportInterval, err)
		}
	}

	return nil
}
