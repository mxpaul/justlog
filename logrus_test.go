package justlog

import (
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

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
