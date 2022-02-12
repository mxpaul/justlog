package justlog

import (
	"github.com/sirupsen/logrus"
)

func Die(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

type LoggerConfig struct {
	Level      string
	TimeFormat string
	ShowNoTime bool
}

type Logger interface {
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

func NewLogger(cfg LoggerConfig) (Logger, error) {
	return NewLogrusLogger(cfg)
}
