package recorders

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"

	"google.golang.org/api/cloudtrace/v1"

	tracing "github.com/Workiva/tracing/lib/go"
	"github.com/opentracing/opentracing-go/ext"
)

// GOOGLE TRACE OBJECTS AND API
// https://cloud.google.com/trace/api/

// Trace
// {
//   "projectId": string, -- Project ID of the Cloud project where the trace data is stored.
//   "traceId": string,   -- Globally unique identifier for the trace.  This identifier is
//													 a 128-bit numeric value formatted as a 32-byte hex string.
//   "spans": [					  -- Collection of spans in the trace.
//     {
//       object(TraceSpan)
//     }
//   ],
// }

// TraceSpan (Span)
// {
//   "spanId": string,       -- Identifier for the span. This identifier must be unique within a trace.
//   "kind": enum(SpanKind), -- Distinguishes between spans generated in a particular context. For example,
//															two spans with the same name may be distinguished using RPC_CLIENT and
//															RPC_SERVER to identify queueing latency associated with the span.
//   "name": string,         -- Name of the trace. The trace name is sanitized and displayed in the
//															Stackdriver Trace tool in the Google Cloud Platform Console. The name may be a
//															method name or some other per-call site name. For the same executable and the
//															same call point, a best practice is to use a consistent name, which makes it
//															easier to correlate cross-trace spans.
//   "startTime": string,    --	Start time of the span in nanoseconds from the UNIX epoch.
//															A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example:
//															"2014-10-02T15:01:23.045123456Z".
//   "endTime": string,      -- End time of the span in nanoseconds from the UNIX epoch.
//															A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds. Example:
//															"2014-10-02T15:01:23.045123456Z".
//   "parentSpanId": string, --	ID of the parent span, if any. Optional.
//   "labels": {             -- Collection of labels associated with the span.
//															An object containing a list of "key": value pairs.
//															Example: { "name": "wrench", "mass": "1.3kg", "count": "3" }.
//     string: string,
//     ...
//   },
// }

// SpanKind
// ENUM: Type of span. Can be used to specify additional relationships between spans
// in addition to a parent/child relationship.
//
// SPAN_KIND_UNSPECIFIED	  -- Unspecified.
// RPC_SERVER								-- Indicates that the span covers server-side handling of an
//													   RPC or other remote network request.
// RPC_CLIENT								-- Indicates that the span covers the client-side wrapper
//														 around an RPC or other remote request.

const SCOPE_URL = "https://www.googleapis.com/auth/cloud-platform"

type GoogleSpanKind string

const (
	SPAN_KIND_UNSPECIFIED GoogleSpanKind = "SPAN_KIND_UNSPECIFIED"
	RPC_SERVER            GoogleSpanKind = "RPC_SERVER"
	RPC_CLIENT            GoogleSpanKind = "RPC_CLIENT"
)

type Labels map[string]string

// GoogleHttpTraceRecorder implements the tracing.Recorder interface.
type GoogleHttpTraceRecorder struct {
	processName string
	tags        map[string]string
	service     *cloudtrace.ProjectsService
	projectId   string
}

// NewGoogleHttpTraceRecorder returns a GoogleHttpTraceRecorder for the given `processName`.
func NewGoogleJWTHttpTraceRecorder(processName, projectId, jwtKeyPath string) (*GoogleHttpTraceRecorder, error) {
	config, err := newJWTConfigFromJSONFile(jwtKeyPath)

	if err != nil {
		return nil, err
	}

	client := newOAuthClient(config)

	cloudtraceService, err := cloudtrace.New(client)

	if err != nil {
		return nil, err
	}

	return &GoogleHttpTraceRecorder{
		processName: processName,
		tags:        make(map[string]string),
		service:     cloudtraceService.Projects,
		projectId:   projectId,
	}, nil
}

// NewGoogleHttpTraceRecorder returns a GoogleHttpTraceRecorder for the given `processName`.
func NewGoogleHttpTraceRecorder(processName, projectId string, client *http.Client) *GoogleHttpTraceRecorder {
	cloudtraceService, err := cloudtrace.New(client)

	if err != nil {
		panic(err)
	}

	return &GoogleHttpTraceRecorder{
		processName: processName,
		tags:        make(map[string]string),
		service:     cloudtraceService.Projects,
		projectId:   projectId,
	}
}

func newJWTConfigFromJSONFile(keyFilePath string) (*jwt.Config, error) {
	scope := SCOPE_URL
	keyFile, err := ioutil.ReadFile(keyFilePath)

	if err != nil {
		panic(err)
	}

	return google.JWTConfigFromJSON(keyFile, scope)
}

func newOAuthClient(config *jwt.Config) *http.Client {
	return config.Client(oauth2.NoContext)
}

// TIMES
// Start time of the span in nanoseconds from the UNIX epoch.
// A timestamp in RFC3339 UTC "Zulu" format, accurate to nanoseconds.
// Example: "2014-10-02T15:01:23.045123456Z".

func (t *GoogleHttpTraceRecorder) rawSpanToGoogleTraceSpan(rawSpan tracing.RawSpan) *cloudtrace.TraceSpan {
	// Doesn't handle Kind or Labels yet

	labels := make(Labels)

	span := &cloudtrace.TraceSpan{
		Name:      rawSpan.Operation,
		SpanId:    uint64(rawSpan.Context.SpanID),
		StartTime: rawSpan.Start.Format("2006-01-02T15:04:05.000000000Z"),
		EndTime:   rawSpan.Start.Add(rawSpan.Duration).Format("2006-01-02T15:04:05.000000000Z"),
		Labels:    labels,
	}

	for k, v := range rawSpan.Tags {
		if k == string(ext.SpanKind) {
			switch v {
			case ext.SpanKindRPCClientEnum:
				span.Kind = string(RPC_CLIENT)
			case ext.SpanKindRPCServerEnum:
				span.Kind = string(RPC_SERVER)
			default:
				span.Kind = string(SPAN_KIND_UNSPECIFIED)
			}
		} else if k == "Actions" {
			for _, action := range v.([]tracing.Action) {
				labels["Ipv4"] = action.Endpoint.Ipv4
				labels["Port"] = string(action.Endpoint.Port)
			}
		} else {
			labels[k] = fmt.Sprintf("%+v", v)
		}
	}

	span.Labels = labels

	if rawSpan.ParentSpanID > 0 {
		span.ParentSpanId = uint64(rawSpan.ParentSpanID)
	}

	return span
}

// ProcessName returns the process name.
func (t *GoogleHttpTraceRecorder) ProcessName() string { return t.processName }

// SetTag sets a tag.
func (t *GoogleHttpTraceRecorder) SetTag(key string, val interface{}) *GoogleHttpTraceRecorder {
	t.tags[key] = fmt.Sprint(val)
	return t
}

// RecordSpan complies with the tracing.Recorder interface.
func (t *GoogleHttpTraceRecorder) RecordSpan(span tracing.RawSpan) {
	traceId := strings.Repeat(strings.ToLower(strconv.FormatUint(span.Context.TraceID, 16)), 2)

	traceSpan := t.rawSpanToGoogleTraceSpan(span)

	trace, err := t.GetExistingTrace(traceId)

	if err != nil {
		fmt.Println(err)
	}

	if trace == nil || trace.TraceId == "" || len(trace.Spans) == 0 {
		trace = &cloudtrace.Trace{
			ProjectId: t.projectId,
			TraceId:   traceId,
			Spans:     []*cloudtrace.TraceSpan{},
		}
	}

	log.Println("========== TRACE: ", trace.TraceId, span.Operation)

	trace.Spans = append(trace.Spans, traceSpan)

	traces := &cloudtrace.Traces{Traces: []*cloudtrace.Trace{trace}}

	_, err = t.service.PatchTraces(t.projectId, traces).Do()

	if err != nil {
		log.Println(err)
	}
}

func (t *GoogleHttpTraceRecorder) GetExistingTrace(traceId string) (*cloudtrace.Trace, error) {
	service := t.service.Traces

	return service.Get(t.projectId, traceId).Do()
}
