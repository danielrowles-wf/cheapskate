package messaging

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/Workiva/frugal/lib/go"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type FrugalTextMapCarrier frugal.FContext

const (
	cid           = "_cid"
	correlationId = "correlation.id"
	prefix        = "ot-tracer-"
)

func (c FrugalTextMapCarrier) Set(key, val string) {
	fc := frugal.FContext(c)
	fc.AddRequestHeader(key, val)
}

func (c FrugalTextMapCarrier) ForeachKey(handler func(key, val string) error) error {
	fc := frugal.FContext(c)

	for k, v := range fc.RequestHeaders() {
		if err := handler(k, v); err != nil {
			return err
		}
	}

	return nil
}

// This is the server side and would need to pull info out of headers, etc
func NewServerTracingMiddleware() frugal.ServiceMiddleware {
	return NewTracingMiddleware(false)
}

func NewClientTracingMiddleware() frugal.ServiceMiddleware {
	return NewTracingMiddleware(true)
}

func NewTracingMiddleware(isClient bool) frugal.ServiceMiddleware {
	tracer := opentracing.GlobalTracer()

	return func(next frugal.InvocationHandler) frugal.InvocationHandler {
		return func(service reflect.Value, method reflect.Method, args frugal.Arguments) frugal.Results {
			serviceName := makeServiceName(service.Type().String(), method.Name)

			span, _ := newSpan(tracer, serviceName, args.Context(), isClient)

			addMethodAttributesToSpan(method, span)

			if span != nil {
				defer closeSpan(span)
			}

			return next(service, method, args)
		}
	}
}

func addMethodAttributesToSpan(method reflect.Method, span opentracing.Span) {
	span.SetTag("source.location", method.PkgPath)
	span.SetTag("method.type", method.Type)
}

func makeServiceName(serviceType string, methodName string) string {
	// TODO: Come up with a better naming convention. Right now client and server
	// will not match. Which may not be an issue. We for sure want to hook into whatever
	// service discovery standards we come up with.
	return strings.Replace(
		fmt.Sprintf("%s.%s", serviceType, methodName), "*", "", 1)
}

func newSpan(tracer opentracing.Tracer, serviceName string, fctx *frugal.FContext, client bool) (opentracing.Span, error) {
	carrier := FrugalTextMapCarrier(*fctx)

	rpcType := ext.SpanKindRPCClientEnum
	if !client {
		rpcType = ext.SpanKindRPCServerEnum
	}

	ctx, _ := tracer.Extract(opentracing.TextMap, carrier)

	var span opentracing.Span
	if ctx == nil {
		span = tracer.StartSpan(serviceName)
	} else {
		span = opentracing.StartSpan(serviceName, ext.RPCServerOption(ctx))
	}

	span.SetTag(string(ext.SpanKind), rpcType)
	span.SetTag(correlationId, fctx.CorrelationID())
	span.Tracer().Inject(span.Context(), opentracing.TextMap, carrier)

	return span, nil
}

func closeSpan(span opentracing.Span) {
	span.SetTag("page.size", strconv.Itoa(os.Getpagesize()))
	span.SetTag("process.id", strconv.Itoa(os.Getpid()))

	hostName, err := os.Hostname()
	if err != nil {
		hostName = "unknown"
	}
	span.SetTag("peer.hostname", hostName)

	span.SetTag("cpu.count", strconv.Itoa(runtime.NumCPU()))
	span.SetTag("goroutine.count", strconv.Itoa(runtime.NumGoroutine()))
	span.SetTag("go.version", runtime.Version())

	span.Finish()
}
