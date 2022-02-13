package justlog

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

type TestLoggerCall struct {
	Method string
	Format string
	Args   []interface{}
}
type TestLoggerCalls []TestLoggerCall

type TestCase_Logger struct {
	Config       *LoggerConfig
	Calls        TestLoggerCalls
	TimeSequence []time.Time
	WantOutput   string
	WantExit     bool
}

// ----------------------------------------------------------------------------
// LogrusBasedLogger logging methods test cases
// ----------------------------------------------------------------------------
func (tc TestCase_Logger) Run(t *testing.T) {

	patches := gomonkey.NewPatches()
	defer patches.Reset()

	if tc.TimeSequence != nil {
		timeSequence := make([]gomonkey.OutputCell, 0, len(tc.TimeSequence))
		for _, timeValue := range tc.TimeSequence {
			timeSequence = append(timeSequence, gomonkey.OutputCell{Values: gomonkey.Params{timeValue}})
		}
		patches.ApplyFuncSeq(time.Now, timeSequence)
	}

	exitCount := 0
	patches.ApplyFunc(os.Exit, func(code int) {
		assert.Equal(t, 1, code)
		exitCount++
	})

	logger, err := NewLogger(*tc.Config)
	assert.NoError(t, err)

	var out strings.Builder
	logger.SetOutput(&out)

	for _, call := range tc.Calls {
		switch call.Method {
		case "Trace":
			logger.Trace(call.Args...)
		case "Debug":
			logger.Debug(call.Args...)
		case "Info":
			logger.Info(call.Args...)
		case "Print":
			logger.Print(call.Args...)
		case "Warn":
			logger.Warn(call.Args...)
		case "Error":
			logger.Error(call.Args...)
		case "Fatal":
			logger.Fatal(call.Args...)
		default:
			assert.Failf(t, "method %q", call.Method)
			return
		}
	}

	assert.Equal(t, tc.WantOutput, out.String())
	if tc.WantExit {
		assert.Equal(t, 1, exitCount, "os.Exit(1) calls")
	}
}

func Test_LogrusBasedLogger_Trace_LevelTrace(t *testing.T) {
	TestCase_Logger{
		Calls: TestLoggerCalls{
			{Method: "Trace", Args: []interface{}{"log message"}},
		},
		Config: &LoggerConfig{
			Level: "trace",
		},
		TimeSequence: []time.Time{
			time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
			time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		},
		WantOutput: "2021-02-01 03:04:05.009000[+2.001000] [TRC] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_Debug_LevelInfo(t *testing.T) {
	TestCase_Logger{
		Calls: TestLoggerCalls{
			{Method: "Debug", Args: []interface{}{"log message"}},
		},
		Config: &LoggerConfig{
			Level: "info",
		},
		WantOutput: "",
	}.Run(t)
}

func Test_LogrusBasedLogger_Debug_LevelDebug(t *testing.T) {
	TestCase_Logger{
		Calls: TestLoggerCalls{
			{Method: "Debug", Args: []interface{}{"log message"}},
		},
		Config: &LoggerConfig{
			Level: "debug",
		},
		TimeSequence: []time.Time{
			time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
			time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		},
		WantOutput: "2021-02-01 03:04:05.009000[+2.001000] [DBG] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_Info_LevelDefault(t *testing.T) {
	TestCase_Logger{
		Calls: TestLoggerCalls{
			{Method: "Info", Args: []interface{}{"log message"}},
		},
		Config: &LoggerConfig{},
		TimeSequence: []time.Time{
			time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
			time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		},
		WantOutput: "2021-02-01 03:04:05.009000[+2.001000] [INF] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_Print_LevelDefault(t *testing.T) {
	TestCase_Logger{
		Calls: TestLoggerCalls{
			{Method: "Print", Args: []interface{}{"log message"}},
		},
		Config: &LoggerConfig{},
		TimeSequence: []time.Time{
			time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
			time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		},
		WantOutput: "2021-02-01 03:04:05.009000[+2.001000] [INF] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_Error_LevelDefault(t *testing.T) {
	TestCase_Logger{
		Calls: TestLoggerCalls{
			{Method: "Error", Args: []interface{}{"log message"}},
		},
		Config: &LoggerConfig{},
		TimeSequence: []time.Time{
			time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
			time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		},
		WantOutput: "2021-02-01 03:04:05.009000[+2.001000] [ERR] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_Fatal_LevelInfo(t *testing.T) {
	TestCase_Logger{
		Calls: TestLoggerCalls{
			{Method: "Fatal", Args: []interface{}{"log message"}},
		},
		Config: &LoggerConfig{
			Level: "info",
		},
		TimeSequence: []time.Time{
			time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
			time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		},
		WantOutput: "2021-02-01 03:04:05.009000[+2.001000] [ERR][FATAL] log message\n",
		WantExit:   true,
	}.Run(t)
}

func Test_LogrusBasedLogger_Info_ShowNoTime(t *testing.T) {
	TestCase_Logger{
		Calls: TestLoggerCalls{
			{Method: "Info", Args: []interface{}{"log message"}},
		},
		Config: &LoggerConfig{
			Level:      "debug",
			ShowNoTime: true,
		},
		TimeSequence: []time.Time{
			time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
			time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		},
		WantOutput: "[+2.001000] [INF] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_DebugInfo_Sequence(t *testing.T) {
	TestCase_Logger{
		Calls: TestLoggerCalls{
			{Method: "Debug", Args: []interface{}{"msg1"}},
			{Method: "Info", Args: []interface{}{"msg2"}},
		},
		Config: &LoggerConfig{
			Level: "debug",
		},
		TimeSequence: []time.Time{
			time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
			time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
			time.Date(2021, time.Month(2), 1, 3, 4, 6, 9500000, time.UTC),
		},
		WantOutput: "2021-02-01 03:04:05.009000[+2.001000] [DBG] msg1\n2021-02-01 03:04:06.009500[+1.000500] [INF] msg2\n",
	}.Run(t)
}
