package justlog

import (
	"fmt"
	"io"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	stringLevelTrace = []byte("TRC")
	stringLevelDebug = []byte("DBG")
	stringLevelInfo  = []byte("INF")
	stringLevelWarn  = []byte("WRN")
	stringLevelErr   = []byte("ERR")
	stringLevelFatal = []byte("DIE")
	stringLevelWTF   = []byte("WTF")
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

	fmtr := LogrusFormatter{
		PrevTime: time.Now(),
	}
	logger.Log.SetFormatter(&fmtr)

	return logger, nil
}

type LogrusBasedLogger struct {
	Log *logrus.Logger
}

func (logger *LogrusBasedLogger) SetOutput(out io.Writer) {
	logger.Log.SetOutput(out)
}

func (logger *LogrusBasedLogger) Trace(args ...interface{}) {
	logger.Log.Trace(args...)
}

func (logger *LogrusBasedLogger) Tracef(format string, args ...interface{}) {
	logger.Log.Tracef(format, args...)
}

func (logger *LogrusBasedLogger) Debug(args ...interface{}) {
	logger.Log.Debug(args...)
}

func (logger *LogrusBasedLogger) Debugf(format string, args ...interface{}) {
	logger.Log.Debugf(format, args...)
}

func (logger *LogrusBasedLogger) Info(args ...interface{}) {
	logger.Log.Info(args...)
}

func (logger *LogrusBasedLogger) Infof(format string, args ...interface{}) {
	logger.Log.Infof(format, args...)
}

func (logger *LogrusBasedLogger) Warn(args ...interface{}) {
	logger.Log.Warn(args...)
}

func (logger *LogrusBasedLogger) Warnf(format string, args ...interface{}) {
	logger.Log.Warnf(format, args...)
}

func (logger *LogrusBasedLogger) Error(args ...interface{}) {
	logger.Log.Error(args...)
}

func (logger *LogrusBasedLogger) Errorf(format string, args ...interface{}) {
	logger.Log.Errorf(format, args...)
}

func (logger *LogrusBasedLogger) Fatal(args ...interface{}) {
	logger.Log.Fatal(args...)
}

func (logger *LogrusBasedLogger) Fatalf(format string, args ...interface{}) {
	logger.Log.Fatalf(format, args...)
}

type LogrusFormatter struct {
	PrevTime time.Time
}

func (f *LogrusFormatter) Format(ent *logrus.Entry) ([]byte, error) {
	buf := make([]byte, 0, 0) // FIXME: use ent.Buffer

	sinceLastLog := ent.Time.Sub(f.PrevTime) // FIXME: store time with atomic
	f.PrevTime = ent.Time                    // TODO: cover with tests

	buf = append(buf, ent.Time.Format("2006-01-02 15:04:05.000000")...)
	buf = append(buf, "[+"...)
	buf = append(buf, f.durationSecondsString(sinceLastLog)...)
	buf = append(buf, ']')
	buf = append(buf, ' ')
	buf = append(buf, '[')
	buf = append(buf, logLevelString(ent.Level)...)
	buf = append(buf, ']')
	buf = append(buf, ' ')
	buf = append(buf, ent.Message...)
	buf = append(buf, '\n')
	return buf, nil
}

func (f *LogrusFormatter) durationSecondsString(d time.Duration) string {
	return fmt.Sprintf("%.6f", d.Seconds())
}

func logLevelString(lvl logrus.Level) []byte {
	switch lvl {
	case logrus.DebugLevel:
		return stringLevelDebug
	case logrus.InfoLevel:
		return stringLevelInfo
	case logrus.FatalLevel:
		return stringLevelFatal
	}
	return stringLevelWTF
}
