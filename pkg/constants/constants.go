package constants

import "time"

const (
	ConfigPath   = "CONFIG_PATH"
	AppEnv       = "APP_ENV"
	AppRootPath  = "APP_ROOT_PATH"
	HttpRegistry = "HTTP_REGISTRY"
	GrpcRegistry = "GRPC_REGISTRY"
	KafkaBrokers = "KAFKA_BROKERS"

	Yaml = "yaml"

	MaxHeaderBytes = 1 << 20 // 1 MB
	StackSize      = 1 << 10 // 1 KB
	BodyLimit      = "2M"
	GzipLevel      = 5

	ReadTimeout          = 15 * time.Second
	WriteTimeout         = 15 * time.Second
	WaitShutdownDuration = 30 * time.Second

	Dev        = "development"
	Test       = "test"
	Production = "production"
)
