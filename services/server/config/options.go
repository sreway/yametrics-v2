package config

import (
	"fmt"
	"net"
	"strconv"
)

type OptionServer func(config *Config) error

func WithAddr(address string) OptionServer {
	return func(cfg *Config) error {
		_, port, err := net.SplitHostPort(address)
		if err != nil {
			return NewConfigError(fmt.Errorf("invalid address %s", address))
		}

		_, err = strconv.Atoi(port)
		if err != nil {
			return NewConfigError(fmt.Errorf("invalid port %s", port))
		}

		cfg.HTTP.Address = address
		return nil
	}
}
