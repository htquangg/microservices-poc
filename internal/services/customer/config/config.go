package config

import (
	"os"

	"github.com/htquangg/microservices-poc/internal/config"
	"github.com/htquangg/microservices-poc/internal/kafka"
	"github.com/htquangg/microservices-poc/pkg/constants"
	"github.com/htquangg/microservices-poc/pkg/database"
	"github.com/htquangg/microservices-poc/pkg/discovery/consul"
	"github.com/htquangg/microservices-poc/pkg/rpc"
	"github.com/htquangg/microservices-poc/pkg/web"
)

type (
	Config struct {
		Name        string `mapstructure:"name,omitempty"`
		Environment string `mapstructure:"environment,omitempty"`

		Web    *web.Config      `mapstructure:"web,omitempty"`
		Rpc    *rpc.Config      `mapstructure:"rpc,omitempty"`
		Mysql  *database.Config `mapstructure:"mysql,omitempty"`
		Consul *consul.Config   `mapstructure:"consul"`
		Kafka  *kafka.Config    `mapstructure:"kafka"`
	}
)

func InitConfig() (*Config, error) {
	cfg := &Config{}
	_, err := config.LoadConfig(cfg)

	grpcRegistry := os.Getenv(constants.GrpcRegistry)
	if grpcRegistry != "" {
		cfg.Rpc.Registry = grpcRegistry
	}

	httpRegistry := os.Getenv(constants.HttpRegistry)
	if httpRegistry != "" {
		cfg.Web.Registry = httpRegistry
	}

	kafkaBrokers := os.Getenv(constants.KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}

	return cfg, err
}

func (cfg Config) IsDevelopment() bool {
	return cfg.Environment == constants.Dev
}

func (cfg Config) IsProduction() bool {
	return cfg.Environment == constants.Production
}
