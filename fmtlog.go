package justlog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

func NewFmtBasedLogger(cfg LoggerConfig) (*FmtBasedLogger, error) {
	logLevel, err := ParseLogLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("ParseLevel error: %w", err)
	}

	logger := &FmtBasedLogger{
		TimeFormat: "2006-01-02 15:04:05.000000",
		PrevTime:   time.Now(),
		ShowNoTime: cfg.ShowNoTime,
		Level:      logLevel,
		Out:        os.Stderr,
	}

	if cfg.TimeFormat != "" {
		logger.TimeFormat = cfg.TimeFormat
	}

	return logger, nil
}

type FmtBasedLogger struct {
	PrevTime   time.Time
	TimeFormat string
	ShowNoTime bool
	Level      Level
	Out        io.Writer
	outMu      sync.Mutex
}

func (logger *FmtBasedLogger) WriteMessage(Level Level, Time time.Time, args ...interface{}) {
	if logger.Level > Level {
		return
	}
	buf := make([]byte, 0, 200)
	msg := logger.MessageBytes(nil, args...)
	buf = logger.FormatMessage(buf, msg, Level, Time)
	logger.outMu.Lock()
	logger.Out.Write(buf)
	logger.outMu.Unlock()
}

func (logger *FmtBasedLogger) WriteMessagef(Level Level, Time time.Time, format string, args ...interface{}) {
	if logger.Level > Level {
		return
	}
	buf := make([]byte, 0, 200)
	msg := []byte(fmt.Sprintf(format, args...))
	buf = logger.FormatMessage(buf, msg, Level, Time)

	logger.outMu.Lock()
	logger.Out.Write(buf)
	logger.outMu.Unlock()
}

func (logger *FmtBasedLogger) FormatMessage(buf []byte, Message []byte, Level Level, Time time.Time) []byte {
	sinceLastLog := Time.Sub(logger.PrevTime) // FIXME: store time with atomic
	logger.PrevTime = Time

	if !logger.ShowNoTime {
		buf = append(buf, Time.Format(logger.TimeFormat)...)
	}

	buf = append(buf, "[+"...)
	buf = append(buf, fmt.Sprintf("%.6f", sinceLastLog.Seconds())...)
	buf = append(buf, ']')
	buf = append(buf, ' ')
	buf = append(buf, logLevelStringLocal(Level)...)
	buf = append(buf, ' ')
	buf = append(buf, Message...)
	buf = append(buf, '\n')
	return buf
}

func (logger *FmtBasedLogger) MessageBytes(buf []byte, args ...interface{}) []byte {
	for i := 0; i < len(args); i++ {
		switch arg := args[i].(type) {
		case string:
			buf = append(buf, arg...)
		case []byte:
			buf = append(buf, arg...)
		default:
			buf = append(buf, fmt.Sprintf("%+v", arg)...)
		}
	}
	return buf
}

func (logger *FmtBasedLogger) SetOutput(out io.Writer) {
	logger.Out = out
}

func (logger *FmtBasedLogger) Trace(args ...interface{}) {
	logger.WriteMessage(LogLevelTrace, time.Now(), args...)
}

func (logger *FmtBasedLogger) Tracef(format string, args ...interface{}) {
	logger.WriteMessagef(LogLevelTrace, time.Now(), format, args...)
}

func (logger *FmtBasedLogger) Debug(args ...interface{}) {
	logger.WriteMessage(LogLevelDebug, time.Now(), args...)
}

func (logger *FmtBasedLogger) Debugf(format string, args ...interface{}) {
	logger.WriteMessagef(LogLevelDebug, time.Now(), format, args...)
}

func (logger *FmtBasedLogger) Info(args ...interface{}) {
	logger.WriteMessage(LogLevelInfo, time.Now(), args...)
}

func (logger *FmtBasedLogger) Infof(format string, args ...interface{}) {
	logger.WriteMessagef(LogLevelInfo, time.Now(), format, args...)
}

func (logger *FmtBasedLogger) Warn(args ...interface{}) {
	logger.WriteMessage(LogLevelWarn, time.Now(), args...)
}

func (logger *FmtBasedLogger) Warnf(format string, args ...interface{}) {
	logger.WriteMessagef(LogLevelWarn, time.Now(), format, args...)
}

func (logger *FmtBasedLogger) Error(args ...interface{}) {
	logger.WriteMessage(LogLevelError, time.Now(), args...)
}

func (logger *FmtBasedLogger) Errorf(format string, args ...interface{}) {
	logger.WriteMessagef(LogLevelError, time.Now(), format, args...)
}

func (logger *FmtBasedLogger) Fatal(args ...interface{}) {
	logger.WriteMessage(LogLevelFatal, time.Now(), args...)
	os.Exit(1)
}

func (logger *FmtBasedLogger) Fatalf(format string, args ...interface{}) {
	logger.WriteMessagef(LogLevelFatal, time.Now(), format, args...)
	os.Exit(1)
}
