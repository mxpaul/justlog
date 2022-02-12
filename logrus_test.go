package justlog

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

type TestCase_LogrusBasedLogger struct {
	Level       string
	Method      string
	Args        []interface{}
	PrevLogTime time.Time
	LogTime     time.Time
	WantOutput  string
	WantExit    bool
}

func (tc TestCase_LogrusBasedLogger) Run(t *testing.T) {

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	timeSequence := []gomonkey.OutputCell{
		{Values: gomonkey.Params{tc.PrevLogTime}},
		{Values: gomonkey.Params{tc.LogTime}},
	}

	patches.ApplyFuncSeq(time.Now, timeSequence)

	exitCount := 0
	patches.ApplyFunc(os.Exit, func(code int) {
		assert.Equal(t, 1, code)
		exitCount++
	})

	logConfig := LoggerConfig{
		Level: tc.Level,
	}
	logger, err := NewLogrusLogger(logConfig)
	assert.NoError(t, err)

	var out strings.Builder
	logger.SetOutput(&out)

	switch tc.Method {
	case "Debug":
		logger.Debug(tc.Args...)
	case "Fatal":
		logger.Fatal(tc.Args...)
	default:
		assert.Failf(t, "method %q", tc.Method)
	}

	assert.Equal(t, tc.WantOutput, out.String())
	if tc.WantExit {
		assert.Equal(t, 1, exitCount, "os.Exit(1) calls")
	}
}

func Test_LogrusBasedLogger_Debug_LevelInfo(t *testing.T) {
	TestCase_LogrusBasedLogger{
		Method:     "Debug",
		Level:      "info",
		Args:       []interface{}{"log message"},
		WantOutput: "",
	}.Run(t)
}

func Test_LogrusBasedLogger_Debug_LevelDebug(t *testing.T) {
	TestCase_LogrusBasedLogger{
		Method:      "Debug",
		Level:       "debug",
		Args:        []interface{}{"log message"},
		PrevLogTime: time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
		LogTime:     time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		WantOutput:  "2021-02-01 03:04:05.009000[+2.001000] [DBG] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_Fatal_LevelInfo(t *testing.T) {
	TestCase_LogrusBasedLogger{
		Method:      "Fatal",
		Level:       "info",
		Args:        []interface{}{"log message"},
		PrevLogTime: time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
		LogTime:     time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		WantOutput:  "2021-02-01 03:04:05.009000[+2.001000] [DIE] log message\n",
		WantExit:    true,
	}.Run(t)
}
