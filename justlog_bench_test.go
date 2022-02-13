package justlog

import (
	"bytes"
	"log"
	"testing"

	"github.com/sirupsen/logrus"
)

/* 2022-02-13

go test -mod=vendor -bench=. -benchmem
goos: linux
goarch: amd64
pkg: github.com/mxpaul/justlog
cpu: Intel(R) Core(TM) i7-1065G7 CPU @ 1.30GHz
BenchmarkFmtBasedLogger-8                 362210              3140 ns/op            1163 B/op         15 allocs/op
BenchmarkLogrusBasedLogger-8              307906              3787 ns/op            1176 B/op         24 allocs/op
BenchmarkLogrusNativeLogger-8             208219              5546 ns/op            2110 B/op         54 allocs/op
BenchmarkLogrusNativeEntry-8              219136              5412 ns/op            1936 B/op         51 allocs/op
BenchmarkGolangLogStd-8                   733635              1621 ns/op             700 B/op          5 allocs/op
PASS
ok      github.com/mxpaul/justlog       6.056s

*/

func BenchmarkFmtBasedLogger(b *testing.B) {
	logger, err := NewFmtBasedLogger(LoggerConfig{})
	if err != nil {
		b.Fail()
		return
	}

	var out bytes.Buffer
	logger.SetOutput(&out)

	for i := 0; i < b.N; i++ {
		logger.Tracef("format %s", "trace")
		logger.Debugf("format %s", "debug")
		logger.Infof("format %s", "info")
		logger.Warnf("format %s", "warn")
		logger.Errorf("format %s", "error")
	}
}

func BenchmarkLogrusBasedLogger(b *testing.B) {
	logger, err := NewLogrusLogger(LoggerConfig{})
	if err != nil {
		b.Fail()
		return
	}

	var out bytes.Buffer
	logger.SetOutput(&out)

	for i := 0; i < b.N; i++ {
		logger.Tracef("format %s", "trace")
		logger.Debugf("format %s", "debug")
		logger.Infof("format %s", "info")
		logger.Warnf("format %s", "warn")
		logger.Errorf("format %s", "error")
	}
}

func BenchmarkLogrusNativeLogger(b *testing.B) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	var out bytes.Buffer
	logger.SetOutput(&out)

	for i := 0; i < b.N; i++ {
		logger.Tracef("format %s", "trace")
		logger.Debugf("format %s", "debug")
		logger.Infof("format %s", "info")
		logger.Warnf("format %s", "warn")
		logger.Errorf("format %s", "error")
	}
}

func BenchmarkLogrusNativeEntry(b *testing.B) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})

	var out bytes.Buffer
	logger.SetOutput(&out)
	entry := logrus.NewEntry(logger)

	for i := 0; i < b.N; i++ {
		entry.Tracef("format %s", "trace")
		entry.Debugf("format %s", "debug")
		entry.Infof("format %s", "info")
		entry.Warnf("format %s", "warn")
		entry.Errorf("format %s", "error")
	}
}

func BenchmarkGolangLogStd(b *testing.B) {
	var out bytes.Buffer
	logger := log.New(&out, "12345", log.LstdFlags|log.Lmicroseconds)

	for i := 0; i < b.N; i++ {
		logger.Printf("format %s", "trace")
		logger.Printf("format %s", "debug")
		logger.Printf("format %s", "info")
		logger.Printf("format %s", "warn")
		logger.Printf("format %s", "error")
	}
}
