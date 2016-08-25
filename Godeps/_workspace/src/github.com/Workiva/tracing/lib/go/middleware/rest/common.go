package rest

import (
	"net/http"
	"os"
	"runtime"
	"strconv"

	"fmt"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func NewSpan(tracer opentracing.Tracer, serviceName string, header http.Header, isClient bool) (opentracing.Span, error) {
	carrier := opentracing.HTTPHeadersCarrier(header)

	rpcType := ext.SpanKindRPCServerEnum
	if isClient {
		rpcType = ext.SpanKindRPCClientEnum
	}

	ctx, err := tracer.Extract(opentracing.TextMap, carrier)

	if err != nil {
		fmt.Println("ERROR CONFIGURING SPAN: ", err)
	}

	var span opentracing.Span

	if ctx == nil {
		span = tracer.StartSpan(serviceName)
	} else {
		span = opentracing.StartSpan(serviceName, ext.RPCServerOption(ctx))
	}

	span.SetTag(string(ext.SpanKind), rpcType)
	opentracing.GlobalTracer().Inject(span.Context(), opentracing.TextMap, carrier)
	return span, nil
}

func CloseSpan(span opentracing.Span) {
	span.SetTag("page.size", strconv.Itoa(os.Getpagesize()))
	span.SetTag("process.id", strconv.Itoa(os.Getpid()))

	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}
	span.SetTag(string(ext.PeerHostname), hostName)

	span.SetTag("cpu.count", strconv.Itoa(runtime.NumCPU()))
	span.SetTag("goroutine.count", strconv.Itoa(runtime.NumGoroutine()))
	span.SetTag("go.version", runtime.Version())
	span.Finish()
}
