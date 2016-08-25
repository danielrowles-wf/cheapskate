package rest

import (
	opentracing "github.com/opentracing/opentracing-go"
	"net/http"
)

func NewServerTracingMiddleware(next http.Handler) http.Handler {
	tracer := opentracing.GlobalTracer()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serviceName := r.URL.Path
		span, _ := NewSpan(tracer, serviceName, r.Header, false)

		if span != nil {
			defer CloseSpan(span)
		}

		next.ServeHTTP(w, r)
	})
}
