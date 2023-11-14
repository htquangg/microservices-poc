package logger

type Fields map[string]interface{}

type LogType int32

const (
	Zap LogType = 0
)

type LogConfig struct {
	Environment string  `mapstructure:"environment"`
	LogLevel    string  `mapstructure:"level"`
	LogType     LogType `mapstructure:"logType"`
}

type Logger interface {
	Configure(cfg func(internalLog interface{}))
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})
	Debugw(msg string, fields Fields)
	LogType() LogType
	Level() string
	Info(args ...interface{})
	Infof(template string, args ...interface{})
	Infow(msg string, fields Fields)
	Warn(args ...interface{})
	Warnf(template string, args ...interface{})
	WarnMsg(msg string, err error)
	Error(args ...interface{})
	Errorw(msg string, fields Fields)
	Errorf(template string, args ...interface{})
	Err(msg string, err error)
	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})
	Println(args ...interface{})
	Printf(template string, args ...interface{})
	WithName(name string)
}
