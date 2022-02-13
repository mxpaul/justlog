package justlog

import (
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

func NewLogrusLogger(cfg LoggerConfig) (*LogrusBasedLogger, error) {

	logger := &LogrusBasedLogger{
		Log: logrus.New(),
	}

	logLevelText := cfg.Level
	if logLevelText == "" {
		logLevelText = "info"
	}
	logLevel, err := logrus.ParseLevel(logLevelText)
	if err != nil {
		return nil, fmt.Errorf("parse loglevel %q: %w", logLevelText, err)
	}
	logger.Log.SetLevel(logLevel)

	fmtr := NewLogrusFormatter(&cfg)
	logger.Log.SetFormatter(fmtr)

	logger.LogEntry = logrus.NewEntry(logger.Log)

	return logger, nil
}

type LogrusBasedLogger struct {
	Log      *logrus.Logger
	LogEntry *logrus.Entry
}

func (logger *LogrusBasedLogger) SetOutput(out io.Writer) {
	logger.Log.SetOutput(out)
}

func (logger *LogrusBasedLogger) Trace(args ...interface{}) {
	logger.LogEntry.Trace(args...)
}

func (logger *LogrusBasedLogger) Tracef(format string, args ...interface{}) {
	logger.LogEntry.Tracef(format, args...)
}

func (logger *LogrusBasedLogger) Debug(args ...interface{}) {
	logger.LogEntry.Debug(args...)
}

func (logger *LogrusBasedLogger) Debugf(format string, args ...interface{}) {
	logger.LogEntry.Debugf(format, args...)
}

func (logger *LogrusBasedLogger) Info(args ...interface{}) {
	logger.LogEntry.Info(args...)
}

func (logger *LogrusBasedLogger) Infof(format string, args ...interface{}) {
	logger.LogEntry.Infof(format, args...)
}

func (logger *LogrusBasedLogger) Print(args ...interface{}) {
	logger.LogEntry.Info(args...)
}

func (logger *LogrusBasedLogger) Printf(format string, args ...interface{}) {
	logger.LogEntry.Infof(format, args...)
}

func (logger *LogrusBasedLogger) Warn(args ...interface{}) {
	logger.LogEntry.Warn(args...)
}

func (logger *LogrusBasedLogger) Warnf(format string, args ...interface{}) {
	logger.LogEntry.Warnf(format, args...)
}

func (logger *LogrusBasedLogger) Error(args ...interface{}) {
	logger.LogEntry.Error(args...)
}

func (logger *LogrusBasedLogger) Errorf(format string, args ...interface{}) {
	logger.LogEntry.Errorf(format, args...)
}

func (logger *LogrusBasedLogger) Fatal(args ...interface{}) {
	logger.LogEntry.Fatal(args...)
}

func (logger *LogrusBasedLogger) Fatalf(format string, args ...interface{}) {
	logger.LogEntry.Fatalf(format, args...)
}

type LogrusFormatter struct {
	PrevTime   time.Time
	TimeFormat string
	ShowNoTime bool
}

func NewLogrusFormatter(cfg *LoggerConfig) *LogrusFormatter {
	f := &LogrusFormatter{
		TimeFormat: "2006-01-02 15:04:05.000000",
		PrevTime:   time.Now(),
	}

	if cfg == nil {
		return f
	}

	if cfg.TimeFormat != "" {
		f.TimeFormat = cfg.TimeFormat
	}
	f.ShowNoTime = cfg.ShowNoTime

	return f
}

func (f *LogrusFormatter) Format(ent *logrus.Entry) ([]byte, error) {
	buf := ent.Buffer
	if buf == nil {
		buf = &bytes.Buffer{}
	}

	sinceLastLog := ent.Time.Sub(f.PrevTime) // FIXME: store time with atomic
	f.PrevTime = ent.Time

	if !f.ShowNoTime {
		buf.WriteString(ent.Time.Format(f.TimeFormat))
	}

	buf.WriteString("[+")
	buf.WriteString(f.durationSecondsString(sinceLastLog))
	buf.WriteRune(']')
	buf.WriteRune(' ')
	buf.Write(logLevelString(ent.Level))
	buf.WriteRune(' ')
	buf.WriteString(ent.Message)
	buf.WriteRune('\n')
	return buf.Bytes(), nil
}

func (f *LogrusFormatter) durationSecondsString(d time.Duration) string {
	return fmt.Sprintf("%.6f", d.Seconds())
}

func logLevelString(lvl logrus.Level) []byte {
	switch lvl {
	case logrus.TraceLevel:
		return stringLevelTrace
	case logrus.DebugLevel:
		return stringLevelDebug
	case logrus.InfoLevel:
		return stringLevelInfo
	case logrus.ErrorLevel:
		return stringLevelError
	case logrus.FatalLevel:
		return stringLevelFatal
	}
	return stringLevelWTF
}
