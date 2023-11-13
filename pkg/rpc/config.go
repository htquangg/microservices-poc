package rpc

import "fmt"

type Config struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
