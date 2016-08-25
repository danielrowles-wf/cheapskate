package rest

import (
	"net/http"
	neturl "net/url"

	"github.com/Workiva/go-rest/rest"
	opentracing "github.com/opentracing/opentracing-go"
)

func NewRestTracingClient(next rest.InvocationHandler) rest.InvocationHandler {
	tracer := opentracing.GlobalTracer()

	return rest.InvocationHandler(func(c *http.Client, method string, url string, body interface{}, header http.Header) (*rest.Response, error) {
		urlObj, _ := neturl.Parse(url)
		span, _ := NewSpan(tracer, urlObj.Path, header, true)

		if span != nil {
			defer CloseSpan(span)
		}

		return next(c, method, url, body, header)
	})
}
