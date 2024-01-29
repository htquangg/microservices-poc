package rpc

import "fmt"

type Config struct {
	Registry string `mapstructure:"registry,omitempty"`
	Host     string `mapstructure:"host,omitempty"`
	Port     int    `mapstructure:"port,omitempty"`
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
