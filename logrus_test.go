package justlog

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type TestCase_LogrusBasedLogger struct {
	Config      *LoggerConfig
	Method      string
	Args        []interface{}
	PrevLogTime time.Time
	LogTime     time.Time
	WantOutput  string
	WantExit    bool
}

// ----------------------------------------------------------------------------
// LogrusBasedLogger logging methods test cases
// ----------------------------------------------------------------------------
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

	logger, err := NewLogrusLogger(*tc.Config)
	assert.NoError(t, err)

	var out strings.Builder
	logger.SetOutput(&out)

	switch tc.Method {
	case "Trace":
		logger.Trace(tc.Args...)
	case "Debug":
		logger.Debug(tc.Args...)
	case "Info":
		logger.Info(tc.Args...)
	case "Warn":
		logger.Warn(tc.Args...)
	case "Error":
		logger.Error(tc.Args...)
	case "Fatal":
		logger.Fatal(tc.Args...)
	default:
		assert.Failf(t, "method %q", tc.Method)
		return
	}

	assert.Equal(t, tc.WantOutput, out.String())
	if tc.WantExit {
		assert.Equal(t, 1, exitCount, "os.Exit(1) calls")
	}
}

func Test_LogrusBasedLogger_Trace_LevelTrace(t *testing.T) {
	TestCase_LogrusBasedLogger{
		Method: "Trace",
		Config: &LoggerConfig{
			Level: "trace",
		},
		Args:        []interface{}{"log message"},
		PrevLogTime: time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
		LogTime:     time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		WantOutput:  "2021-02-01 03:04:05.009000[+2.001000] [TRC] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_Debug_LevelInfo(t *testing.T) {
	TestCase_LogrusBasedLogger{
		Method: "Debug",
		Config: &LoggerConfig{
			Level: "info",
		},
		Args:       []interface{}{"log message"},
		WantOutput: "",
	}.Run(t)
}

func Test_LogrusBasedLogger_Debug_LevelDebug(t *testing.T) {
	TestCase_LogrusBasedLogger{
		Method: "Debug",
		Config: &LoggerConfig{
			Level: "debug",
		},
		Args:        []interface{}{"log message"},
		PrevLogTime: time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
		LogTime:     time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		WantOutput:  "2021-02-01 03:04:05.009000[+2.001000] [DBG] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_Info_LevelDefault(t *testing.T) {
	TestCase_LogrusBasedLogger{
		Method:      "Info",
		Config:      &LoggerConfig{},
		Args:        []interface{}{"log message"},
		PrevLogTime: time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
		LogTime:     time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		WantOutput:  "2021-02-01 03:04:05.009000[+2.001000] [INF] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_Error_LevelDefault(t *testing.T) {
	TestCase_LogrusBasedLogger{
		Method:      "Error",
		Config:      &LoggerConfig{},
		Args:        []interface{}{"log message"},
		PrevLogTime: time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
		LogTime:     time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		WantOutput:  "2021-02-01 03:04:05.009000[+2.001000] [ERR] log message\n",
	}.Run(t)
}

func Test_LogrusBasedLogger_Fatal_LevelInfo(t *testing.T) {
	TestCase_LogrusBasedLogger{
		Method: "Fatal",
		Config: &LoggerConfig{
			Level: "info",
		},
		Args:        []interface{}{"log message"},
		PrevLogTime: time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
		LogTime:     time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		WantOutput:  "2021-02-01 03:04:05.009000[+2.001000] [ERR][FATAL] log message\n",
		WantExit:    true,
	}.Run(t)
}

func Test_LogrusBasedLogger_Info_ShowNoTime(t *testing.T) {
	TestCase_LogrusBasedLogger{
		Method: "Info",
		Config: &LoggerConfig{
			Level:      "debug",
			ShowNoTime: true,
		},
		Args:        []interface{}{"log message"},
		PrevLogTime: time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
		LogTime:     time.Date(2021, time.Month(2), 1, 3, 4, 5, 9000000, time.UTC),
		WantOutput:  "[+2.001000] [INF] log message\n",
	}.Run(t)
}

// ----------------------------------------------------------------------------
// LogrusFormatter Format test cases
// ----------------------------------------------------------------------------
type logrusFormatCall struct {
	Entry          *logrus.Entry
	WantBytes      string
	WantErrorMatch []string
}

type TestCase_LogrusFormatter_Format struct {
	Call     []logrusFormatCall
	PrevTime time.Time
	Config   *LoggerConfig
}

func (tc TestCase_LogrusFormatter_Format) Run(t *testing.T) {
	formatter := NewLogrusFormatter(tc.Config)

	if !tc.PrevTime.IsZero() {
		formatter.PrevTime = tc.PrevTime
	}

	for _, call := range tc.Call {

		gotBytes, gotError := formatter.Format(call.Entry)

		if call.WantErrorMatch == nil {
			assert.Equal(t, call.WantBytes, string(gotBytes))
			assert.NoError(t, gotError)
		} else if assert.Error(t, gotError) {
			for _, wantMatchString := range call.WantErrorMatch {
				assert.Contains(t, gotError.Error(), wantMatchString)
			}
		}

	}
}

func Test_LogrusFormatter_Format_SavePrevTime(t *testing.T) {
	TestCase_LogrusFormatter_Format{
		PrevTime: time.Date(2021, time.Month(2), 1, 3, 4, 3, 8000000, time.UTC),
		Call: []logrusFormatCall{
			{
				Entry: &logrus.Entry{
					Time:    time.Date(2021, time.Month(2), 1, 3, 4, 4, 8000000, time.UTC),
					Message: "MSG!",
					Level:   logrus.DebugLevel,
				},
				WantBytes: "2021-02-01 03:04:04.008000[+1.000000] [DBG] MSG!\n",
			},
			{
				Entry: &logrus.Entry{
					Time:    time.Date(2021, time.Month(2), 1, 3, 4, 5, 8000000, time.UTC),
					Message: "MSG!",
					Level:   logrus.InfoLevel,
				},
				WantBytes: "2021-02-01 03:04:05.008000[+1.000000] [INF] MSG!\n",
			},
		},
	}.Run(t)
}

func Test_LogrusFormatter_Format_CustomTimeFormat(t *testing.T) {
	TestCase_LogrusFormatter_Format{
		Config: &LoggerConfig{
			TimeFormat: "2006 .000000",
		},
		PrevTime: time.Date(2020, time.Month(5), 6, 1, 2, 3, 7890000, time.UTC),
		Call: []logrusFormatCall{
			{
				Entry: &logrus.Entry{
					Time:    time.Date(2020, time.Month(5), 6, 1, 2, 4, 7900000, time.UTC),
					Message: "MSG!",
					Level:   logrus.DebugLevel,
				},
				WantBytes: "2020 .007900[+1.000010] [DBG] MSG!\n",
			},
		},
	}.Run(t)
}

func Test_LogrusFormatter_Format_OmitTime(t *testing.T) {
	TestCase_LogrusFormatter_Format{
		Config: &LoggerConfig{
			ShowNoTime: true,
		},
		PrevTime: time.Date(2020, time.Month(5), 6, 1, 2, 3, 7890000, time.UTC),
		Call: []logrusFormatCall{
			{
				Entry: &logrus.Entry{
					Time:    time.Date(2020, time.Month(5), 6, 1, 2, 4, 7900000, time.UTC),
					Message: "MSG!",
					Level:   logrus.DebugLevel,
				},
				WantBytes: "[+1.000010] [DBG] MSG!\n",
			},
		},
	}.Run(t)
}
