package database

import (
	"sync/atomic"

	"github.com/htquangg/microservices-poc/pkg/logger"

	xormlog "xorm.io/xorm/log"
)

type XORMLogBridge struct {
	showSQL atomic.Bool
	log     logger.Logger
}

func NewXORMLogger(log logger.Logger, showSQL bool) xormlog.Logger {
	l := &XORMLogBridge{log: log}
	l.showSQL.Store(showSQL)
	return l
}

func (l *XORMLogBridge) Debug(v ...interface{}) {
	l.log.Debug(v...)
}

func (l *XORMLogBridge) Debugf(format string, v ...interface{}) {
	l.log.Debugf(format, v...)
}

func (l *XORMLogBridge) Error(v ...interface{}) {
	l.log.Error(v...)
}

func (l *XORMLogBridge) Errorf(format string, v ...interface{}) {
	l.log.Errorf(format, v...)
}

func (l *XORMLogBridge) Info(v ...interface{}) {
	l.log.Info(v...)
}

func (l *XORMLogBridge) Infof(format string, v ...interface{}) {
	l.log.Infof(format, v...)
}

func (l *XORMLogBridge) Warn(v ...interface{}) {
	l.log.Warn(v...)
}

func (l *XORMLogBridge) Warnf(format string, v ...interface{}) {
	l.log.Warnf(format, v...)
}

func (l *XORMLogBridge) Level() xormlog.LogLevel {
	switch l.log.Level() {
	case "debug":
		return xormlog.LOG_DEBUG
	case "info":
		return xormlog.LOG_INFO
	case "warn":
		return xormlog.LOG_WARNING
	case "error", "panic", "fatal":
		return xormlog.LOG_ERR
	}
	return xormlog.LOG_UNKNOWN
}

func (*XORMLogBridge) SetLevel(xormlog.LogLevel) {
}

func (l *XORMLogBridge) IsShowSQL() bool {
	return l.showSQL.Load()
}

func (l *XORMLogBridge) ShowSQL(show ...bool) {
	if len(show) == 0 {
		show = []bool{true}
	}
	l.showSQL.Store(show[0])
}
