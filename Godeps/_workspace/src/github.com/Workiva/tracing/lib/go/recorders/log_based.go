package recorders

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/Workiva/tracing/lib/go"
	"os"
)

const logType = "tracing"

// LogBasedRecorder implements the tracing.Recorder interface.
type LogBasedRecorder struct {
	processName string
	tags        map[string]string
	logger      logrus.Logger
}

// NewLogBasedRecorder returns a LogBasedRecorder for the given `processName`.
func NewLogBasedRecorder(processName string) *LogBasedRecorder {
	return &LogBasedRecorder{
		processName: processName,
		tags:        make(map[string]string),
		logger: logrus.Logger{
			Out:       os.Stderr,
			Formatter: new(logrus.JSONFormatter),
			Level:     logrus.DebugLevel,
		},
	}
}

// ProcessName returns the process name.
func (t *LogBasedRecorder) ProcessName() string { return t.processName }

// SetTag sets a tag.
func (t *LogBasedRecorder) SetTag(key string, val interface{}) *LogBasedRecorder {
	t.tags[key] = fmt.Sprint(val)
	return t
}

// RecordSpan complies with the tracing.Recorder interface.
func (t *LogBasedRecorder) RecordSpan(span tracing.RawSpan) {
	span.Tags["end"] = span.Start.Add(span.Duration).Format("2006-01-02T15:04:05.000000000Z")
	t.logger.WithFields(logrus.Fields{
		"type": logType,
		"span": tracing.NewJSONSpan(span)}).Debug()
}
