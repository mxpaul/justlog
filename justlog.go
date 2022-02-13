package justlog

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Level uint8

const DefaultTimeFormat = "2006-01-02 15:04:05.000000"

const (
	LogLevelInvalid Level = iota
	LogLevelTrace
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
)

var (
	stringLevelTrace = []byte("[TRC]")
	stringLevelDebug = []byte("[DBG]")
	stringLevelInfo  = []byte("[INF]")
	stringLevelWarn  = []byte("[WRN]")
	stringLevelError = []byte("[ERR]")
	stringLevelFatal = []byte("[ERR][FATAL]")
	stringLevelWTF   = []byte("[WTF]")
)

func logLevelStringLocal(lvl Level) []byte {
	switch lvl {
	case LogLevelTrace:
		return stringLevelTrace
	case LogLevelDebug:
		return stringLevelDebug
	case LogLevelInfo:
		return stringLevelInfo
	case LogLevelError:
		return stringLevelError
	case LogLevelFatal:
		return stringLevelFatal
	}
	return stringLevelWTF
}

func ParseLogLevel(strLevel string) (Level, error) {
	switch strLevel {
	case "trace":
		return LogLevelTrace, nil
	case "debug":
		return LogLevelDebug, nil
	case "info", "":
		return LogLevelInfo, nil
	case "warn":
		return LogLevelWarn, nil
	case "error":
		return LogLevelError, nil
	case "fatal":
		return LogLevelFatal, nil
	}
	return LogLevelInvalid, fmt.Errorf("invalid log level value: %q", strLevel)
}

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
	Print(args ...interface{})
	Printf(format string, args ...interface{})
}

//func NewLogger(cfg LoggerConfig) (*LogrusBasedLogger, error) {
//	return NewLogrusLogger(cfg)
//}

func NewLogger(cfg LoggerConfig) (*FmtBasedLogger, error) {
	return NewFmtBasedLogger(cfg)
}
