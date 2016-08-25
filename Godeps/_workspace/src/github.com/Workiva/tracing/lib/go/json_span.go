package tracing

import (
	opentracing "github.com/opentracing/opentracing-go"
	"time"
)

const SchemaVersion = "0.0.1"

// JSONSpan is the span struct that when turned into json will have the correct fields
type JSONSpan struct {
	// The version of the schema that the json belongs to
	Version string `json:"version"`

	// A probabilistically unique identifier for a [multi-span] trace.
	TraceID uint64 `json:"traceId"`

	// A probabilistically unique identifier for a span.
	SpanID uint64 `json:"spanId"`

	// Whether the trace is sampled.
	Sampled bool `json:"sampled"`

	// The span's associated baggage.
	Baggage map[string]string `json:"baggage"` // initialized on first use

	// The SpanID of this SpanContext's first intra-trace reference (i.e.,
	// "parent"), or 0 if there is no parent.
	ParentSpanID uint64 `json:"parentSpanId"`

	// The name of the "operation" this span is an instance of. (Called a "span
	// name" in some implementations)
	Operation string `json:"operation"`

	// We store <start, duration> rather than <start, end> so that only
	// one of the timestamps has global clock uncertainty issues.
	Start    time.Time     `json:"start"`
	Duration time.Duration `json:"duration"`

	// Essentially an extension mechanism. Can be used for many purposes,
	// not to be enumerated here.
	Tags opentracing.Tags `json:"tags"`

	// The span's "microlog".
	Logs []opentracing.LogData `json:"logs"`
}

func NewJSONSpan(span RawSpan) JSONSpan {
	return JSONSpan{SchemaVersion, span.Context.TraceID, span.Context.SpanID, span.Context.Sampled,
		span.Context.Baggage, span.ParentSpanID, span.Operation, span.Start, span.Duration, span.Tags, span.Logs}
}
