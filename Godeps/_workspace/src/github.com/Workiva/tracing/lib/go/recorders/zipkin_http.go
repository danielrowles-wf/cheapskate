package recorders

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	tracing "github.com/Workiva/tracing/lib/go"
)

//************** Annotation.value **************
/**
 * The client sent ("cs") a request to a server. There is only one send per
 * span. For example, if there's a transport error, each attempt can be logged
 * as a WIRE_SEND annotation.
 *
 * If chunking is involved, each chunk could be logged as a separate
 * CLIENT_SEND_FRAGMENT in the same span.
 *
 * Annotation.host is not the server. It is the host which logged the send
 * event, almost always the client. When logging CLIENT_SEND, instrumentation
 * should also log the SERVER_ADDR.
 */
const CLIENT_SEND = "cs"

/**
 * The client received ("cr") a response from a server. There is only one
 * receive per span. For example, if duplicate responses were received, each
 * can be logged as a WIRE_RECV annotation.
 *
 * If chunking is involved, each chunk could be logged as a separate
 * CLIENT_RECV_FRAGMENT in the same span.
 *
 * Annotation.host is not the server. It is the host which logged the receive
 * event, almost always the client. The actual endpoint of the server is
 * recorded separately as SERVER_ADDR when CLIENT_SEND is logged.
 */
const CLIENT_RECV = "cr"

/**
 * The server sent ("ss") a response to a client. There is only one response
 * per span. If there's a transport error, each attempt can be logged as a
 * WIRE_SEND annotation.
 *
 * Typically, a trace ends with a server send, so the last timestamp of a trace
 * is often the timestamp of the root span's server send.
 *
 * If chunking is involved, each chunk could be logged as a separate
 * SERVER_SEND_FRAGMENT in the same span.
 *
 * Annotation.host is not the client. It is the host which logged the send
 * event, almost always the server. The actual endpoint of the client is
 * recorded separately as CLIENT_ADDR when SERVER_RECV is logged.
 */
const SERVER_SEND = "ss"

/**
 * The server received ("sr") a request from a client. There is only one
 * request per span.  For example, if duplicate responses were received, each
 * can be logged as a WIRE_RECV annotation.
 *
 * Typically, a trace starts with a server receive, so the first timestamp of a
 * trace is often the timestamp of the root span's server receive.
 *
 * If chunking is involved, each chunk could be logged as a separate
 * SERVER_RECV_FRAGMENT in the same span.
 *
 * Annotation.host is not the client. It is the host which logged the receive
 * event, almost always the server. When logging SERVER_RECV, instrumentation
 * should also log the CLIENT_ADDR.
 */
const SERVER_RECV = "sr"

/**
 * Optionally logs an attempt to send a message on the wire. Multiple wire send
 * events could indicate network retries. A lag between client or server send
 * and wire send might indicate queuing or processing delay.
 */
const WIRE_SEND = "ws"

/**
 * Optionally logs an attempt to receive a message from the wire. Multiple wire
 * receive events could indicate network retries. A lag between wire receive
 * and client or server receive might indicate queuing or processing delay.
 */
const WIRE_RECV = "wr"

/**
 * Optionally logs progress of a (CLIENT_SEND, WIRE_SEND). For example, this
 * could be one chunk in a chunked request.
 */
const CLIENT_SEND_FRAGMENT = "csf"

/**
 * Optionally logs progress of a (CLIENT_RECV, WIRE_RECV). For example, this
 * could be one chunk in a chunked response.
 */
const CLIENT_RECV_FRAGMENT = "crf"

/**
 * Optionally logs progress of a (SERVER_SEND, WIRE_SEND). For example, this
 * could be one chunk in a chunked response.
 */
const SERVER_SEND_FRAGMENT = "ssf"

/**
 * Optionally logs progress of a (SERVER_RECV, WIRE_RECV). For example, this
 * could be one chunk in a chunked request.
 */
const SERVER_RECV_FRAGMENT = "srf"

//***** BinaryAnnotation.key ******
/**
 * The domain portion of the URL or host header. Ex. "mybucket.s3.amazonaws.com"
 *
 * Used to filter by host as opposed to ip address.
 */
const HTTP_HOST = "http.host"

/**
 * The HTTP method, or verb, such as "GET" or "POST".
 *
 * Used to filter against an http route.
 */
const HTTP_METHOD = "http.method"

type Endpoint struct {
	// IPv4 host address packed into 4 bytes.
	// Ex for the ip 1.2.3.4, it would be (1 << 24) | (2 << 16) | (3 << 8) | 4
	Ipv4        string `json:"ipv4"`
	Port        int16  `json:"port"`
	ServiceName string `json:"serviceName"`
}

type Annotation struct {
	Timestamp int64    `json:"timestamp"`
	Value     string   `json:"value"`
	Endpoint  Endpoint `json:"endpoint,omitempty"` // This is the system that created the send event
}

type BinaryAnnotation struct {
	// TODO: Add value (binary) and annotation type
	// https://github.com/openzipkin/zipkin/blob/94f357f56dcee5f8f77bcc6860d77334f512a4c2/zipkin-thrift/src/main/thrift/com/twitter/zipkin/zipkinCore.thrift#L283
	Key      string   `json:"key"`
	Endpoint Endpoint `json:"endpoint,omitempty"` // This is the system that created the send event
}

type ZipkinSpan struct {
	Name              string             `json:"name"`                        // RawSpan.Operation
	TraceId           string             `json:"traceId"`                     // RawSpan.TraceID
	Id                string             `json:"id"`                          // RawSpan.SpanID
	ParentId          string             `json:"parentId,omitempty"`          // RawSpan.ParentSpanID
	Timestamp         int64              `json:"timestamp,omitempty"`         // RawSpan.Start (convert from time.Time to epoch)
	Duration          int64              `json:"duration,omitempty"`          // RawSpan.Duration
	Debug             bool               `json:"debug,omitempty"`             // RawSpan.Duration
	Annotations       []Annotation       `json:"annotations,omitempty"`       // RawSpan ???
	BinaryAnnotations []BinaryAnnotation `json:"binaryAnnotations,omitempty"` // RawSpan ???
	// Sampled ??
}

func rawSpanToZipkinSpan(rawSpan tracing.RawSpan) ZipkinSpan {
	return ZipkinSpan{
		Name:      rawSpan.Operation,
		TraceId:   strconv.FormatUint(rawSpan.Context.TraceID, 10),
		Id:        strconv.FormatUint(rawSpan.Context.SpanID, 10),
		Timestamp: int64(rawSpan.Start.Unix()),
		Duration:  int64(rawSpan.Duration),
		Debug:     true,
	}
}

// ZipkinHttpRecorder implements the tracing.Recorder interface.
type ZipkinHttpRecorder struct {
	processName string
	tags        map[string]string
}

// NewZipkinHttpRecorder returns a ZipkinHttpRecorder for the given `processName`.
func NewZipkinHttpRecorder(processName string) *ZipkinHttpRecorder {
	return &ZipkinHttpRecorder{
		processName: processName,
		tags:        make(map[string]string),
	}
}

// ProcessName returns the process name.
func (t *ZipkinHttpRecorder) ProcessName() string { return t.processName }

// SetTag sets a tag.
func (t *ZipkinHttpRecorder) SetTag(key string, val interface{}) *ZipkinHttpRecorder {
	t.tags[key] = fmt.Sprint(val)
	return t
}

// RecordSpan complies with the tracing.Recorder interface.
func (t *ZipkinHttpRecorder) RecordSpan(span tracing.RawSpan) {
	fmt.Println("========================================")

	url := "http://localhost:9411/api/v1/spans"

	zipSpan := rawSpanToZipkinSpan(span)

	if span.ParentSpanID > 0 {
		zipSpan.ParentId = strconv.FormatUint(span.ParentSpanID, 10)
	}

	annotations := []Annotation{}

	actions, ok := span.Tags["Actions"]
	if ok {
		for _, action := range actions.([]tracing.Action) {

			annotations = append(annotations, Annotation{
				Timestamp: action.Timestamp,
				Value:     action.Value,
				Endpoint: Endpoint{
					Ipv4:        action.Endpoint.Ipv4,
					Port:        action.Endpoint.Port,
					ServiceName: action.Endpoint.ServiceName,
				},
			})
		}
	}

	zipSpan.Annotations = annotations

	spans := [1]ZipkinSpan{zipSpan}

	fmt.Println("RAW SPAN: ", span)
	jsonBytes := new(bytes.Buffer)
	json.NewEncoder(jsonBytes).Encode(spans)

	req, err := http.NewRequest("POST", url, jsonBytes)

	jsonStr := jsonBytes.String()
	fmt.Println("REQUEST JSON: ", jsonStr)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.Itoa(jsonBytes.Len()))

	hn, e := os.Hostname()

	if e != nil {
		hn = "UNKNOWN"
	}

	req.Header.Add("Host", hn)

	fmt.Println("REQUEST HEADERS: ", req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("RESPONSE STATUS: ", resp.Status)

	content_type, ok := resp.Header["Content-Type"]

	if ok && content_type[0] == "application/octet-stream" {
		reader := bufio.NewReader(resp.Body)

		for {
			line, err := reader.ReadBytes('\r')
			line = bytes.TrimSpace(line)
			fmt.Println("STREAMING REPONSE: ", string(line))

			if err != nil {
				break
			}
		}
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Println("NON-STREAMING RESPONSE BODY: ", string(body))
	}

	//fmt.Printf(
	//"RecordSpan: %v[%v, %v us] --> %v logs. std context: %v; baggage: %v\n",
	//span.Operation, span.Start, span.Duration, len(span.Logs),
	//span.Context, span.Baggage)
	//for i, l := range span.Logs {
	//fmt.Printf(
	//"    log %v @ %v: %v --> %v\n", i, l.Timestamp, l.Event, reflect.TypeOf(l.Payload))
	//}
	fmt.Println("---------------------------------------")
}
