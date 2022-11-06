package config

import "fmt"

type ErrConfig struct {
	error error
}

func NewConfigError(err error) error {
	return &ErrConfig{
		error: err,
	}
}

func (c *ErrConfig) Error() string {
	return fmt.Sprintf("Config_Error: %s", c.error)
}
