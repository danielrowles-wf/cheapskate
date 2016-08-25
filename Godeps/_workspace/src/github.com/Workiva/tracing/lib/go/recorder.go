package tracing

import (
	"fmt"
	"reflect"
	"sync"
)

// A SpanRecorder handles all of the `RawSpan` data generated via an
// associated `Tracer` (see `NewStandardTracer`) instance. It also names
// the containing process and provides access to a straightforward tag map.
type SpanRecorder interface {
	// Implementations must determine whether and where to store `span`.
	RecordSpan(span RawSpan)
}

// InMemorySpanRecorder is a simple thread-safe implementation of
// SpanRecorder that stores all reported spans in memory, accessible
// via reporter.GetSpans(). It is primarily intended for testing purposes.
type InMemorySpanRecorder struct {
	sync.RWMutex
	spans []RawSpan
}

// NewInMemoryRecorder creates new InMemorySpanRecorder
func NewInMemoryRecorder() *InMemorySpanRecorder {
	return new(InMemorySpanRecorder)
}

// RecordSpan implements the respective method of SpanRecorder.
func (r *InMemorySpanRecorder) RecordSpan(span RawSpan) {
	r.Lock()
	defer r.Unlock()
	r.spans = append(r.spans, span)
}

// GetSpans returns a copy of the array of spans accumulated so far.
func (r *InMemorySpanRecorder) GetSpans() []RawSpan {
	r.RLock()
	defer r.RUnlock()
	spans := make([]RawSpan, len(r.spans))
	copy(spans, r.spans)
	return spans
}

// GetSampledSpans returns a slice of spans accumulated so far which were sampled.
func (r *InMemorySpanRecorder) GetSampledSpans() []RawSpan {
	r.RLock()
	defer r.RUnlock()
	spans := make([]RawSpan, 0, len(r.spans))
	for _, span := range r.spans {
		if span.Context.Sampled {
			spans = append(spans, span)
		}
	}
	return spans
}

// Reset clears the internal array of spans.
func (r *InMemorySpanRecorder) Reset() {
	r.Lock()
	defer r.Unlock()
	r.spans = nil
}

// TrivialRecorder implements the basictracer.Recorder interface.
type TrivialRecorder struct {
	processName string
	tags        map[string]string
}

// NewTrivialRecorder returns a TrivialRecorder for the given `processName`.
func NewTrivialRecorder(processName string) *TrivialRecorder {
	return &TrivialRecorder{
		processName: processName,
		tags:        make(map[string]string),
	}
}

// ProcessName returns the process name.
func (t *TrivialRecorder) ProcessName() string { return t.processName }

// SetTag sets a tag.
func (t *TrivialRecorder) SetTag(key string, val interface{}) *TrivialRecorder {
	t.tags[key] = fmt.Sprint(val)
	return t
}

// RecordSpan complies with the basictracer.Recorder interface.
func (t *TrivialRecorder) RecordSpan(span RawSpan) {
	fmt.Printf(
		"RecordSpan: %v[%v, %v us] --> %v logs. context: %v; baggage: %v\n",
		span.Operation, span.Start, span.Duration, len(span.Logs),
		span.Context, span.Context.Baggage)
	for i, l := range span.Logs {
		fmt.Printf(
			"    log %v @ %v: %v --> %v\n", i, l.Timestamp, l.Event, reflect.TypeOf(l.Payload))
	}
}
