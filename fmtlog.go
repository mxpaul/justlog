package justlog

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

func timeFormatFuncCommon(buf []byte, t time.Time, format string) []byte {
	return append(buf, t.Format(format)...)
}

// stolen from log package
func itoa(buf []byte, i int, wid int) []byte {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	return append(buf, b[bp:]...)
}

// stolen from log package
func timeFormatFuncDefaultCustomized(buf []byte, t time.Time, _ string) []byte {
	year, month, day := t.Date()
	buf = itoa(buf, year, 4)
	buf = append(buf, '-')
	buf = itoa(buf, int(month), 2)
	buf = append(buf, '-')
	buf = itoa(buf, day, 2)
	buf = append(buf, ' ')

	hour, min, sec := t.Clock()
	buf = itoa(buf, hour, 2)
	buf = append(buf, ':')
	buf = itoa(buf, min, 2)
	buf = append(buf, ':')
	buf = itoa(buf, sec, 2)

	buf = append(buf, '.')
	buf = itoa(buf, t.Nanosecond()/1e3, 6)
	return buf
}

func NewFmtBasedLogger(cfg LoggerConfig) (*FmtBasedLogger, error) {
	logLevel, err := ParseLogLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("ParseLevel error: %w", err)
	}

	logger := &FmtBasedLogger{
		TimeFormat:     DefaultTimeFormat,
		PrevTime:       time.Now(),
		ShowNoTime:     cfg.ShowNoTime,
		Level:          logLevel,
		Out:            os.Stderr,
		timeFormatFunc: timeFormatFuncCommon,
	}

	if cfg.TimeFormat != "" {
		logger.TimeFormat = cfg.TimeFormat
	}

	if logger.TimeFormat == DefaultTimeFormat {
		logger.timeFormatFunc = timeFormatFuncDefaultCustomized
	}

	return logger, nil
}

type FmtBasedLogger struct {
	PrevTime       time.Time
	TimeFormat     string
	ShowNoTime     bool
	Level          Level
	Out            io.Writer
	outMu          sync.Mutex
	timeFormatFunc func([]byte, time.Time, string) []byte
}

func (logger *FmtBasedLogger) WriteMessage(Level Level, Time time.Time, args ...interface{}) {
	if logger.Level > Level {
		return
	}
	msg := logger.MessageBytes(nil, args...)
	buf := make([]byte, 0, len(msg)+45)
	buf = logger.FormatMessage(buf, msg, Level, Time)
	logger.outMu.Lock()
	defer logger.outMu.Unlock()
	logger.Out.Write(buf)
}

func (logger *FmtBasedLogger) WriteMessagef(Level Level, Time time.Time, format string, args ...interface{}) {
	if logger.Level > Level {
		return
	}
	msg := []byte(fmt.Sprintf(format, args...))
	buf := make([]byte, 0, len(msg)+45)
	buf = logger.FormatMessage(buf, msg, Level, Time)

	logger.outMu.Lock()
	defer logger.outMu.Unlock()
	logger.Out.Write(buf)
}

func (logger *FmtBasedLogger) FormatMessage(buf []byte, Message []byte, Level Level, Time time.Time) []byte {
	sinceLastLog := Time.Sub(logger.PrevTime) // FIXME: store time with atomic
	logger.PrevTime = Time

	if !logger.ShowNoTime {
		buf = logger.timeFormatFunc(buf, Time, logger.TimeFormat)
	}

	buf = append(buf, "[+"...)
	buf = append(buf, fmt.Sprintf("%.6f", sinceLastLog.Seconds())...) // FIXME: split for seconds and nanoseconds, use itoa
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

func (logger *FmtBasedLogger) Print(args ...interface{}) {
	logger.WriteMessage(LogLevelInfo, time.Now(), args...)
}

func (logger *FmtBasedLogger) Printf(format string, args ...interface{}) {
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
