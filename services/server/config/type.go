package config

import (
	"encoding/json"
	"fmt"
	"time"
)

func (mc *MemoryStorageConfig) UnmarshalJSON(data []byte) error {
	var err error
	type AliasType MemoryStorageConfig
	aliasValue := &struct {
		*AliasType
		StoreInterval string `json:"store_interval"`
	}{
		AliasType: (*AliasType)(mc),
	}

	if err = json.Unmarshal(data, aliasValue); err != nil {
		return err
	}

	if aliasValue.StoreInterval != "" {
		if mc.StoreInterval, err = time.ParseDuration(aliasValue.StoreInterval); err != nil {
			return fmt.Errorf("failed to parse '%s' to time.Duration: %w", aliasValue.StoreInterval, err)
		}
	}

	return nil
}
