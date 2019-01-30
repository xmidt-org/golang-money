package money

import (
	"context"
	"net/http"
)

// Spanner is an interface which represents different ways spanning occurs.
type Spanner interface {
	ServerDecorator(context.Context, *HTTPSpanner, http.Handler, *http.Request, http.ResponseWriter) http.Handler
}

// MoneyContainer holds the minimum primitives to complete money spans for all servers scytale, talaria, petasos,
// talaria
//
// talaria is a special case that needs use of a channel
type MoneyContainer interface {
	ServerDecorator(context.Context, *HTTPSpanner, http.Handler, *http.Request, http.ResponseWriter) http.Handler
}

// ServerDecorator a function signature for server decorators
type ServerDecorator func(context.Context, *HTTPSpanner, http.Handler, *http.Request, http.ResponseWriter) http.Handler

// HTTPSpanner implements a Spanner and its future possible options.
//
// Future created spanner options go here.
type HTTPSpanner struct {
	Tr1d1um Starter
	Scytale SubTracer
	Talaria MoneyContainer
}

func NewHTTPSpanner(options HTTPSpannerOptions) *HTTPSpanner {
	hs := new(HTTPSpanner)

	options(hs)

	return hs
}

// Start defines the start time of the input span s and returns
// a http tracker which can both start a child span using SubTrace
// as well as mark the end of a span using s
func (hs *HTTPSpanner) Start(ctx context.Context, s *Span) *HTTPTracker {
	return NewHTTPTracker(ctx, s, hs)
}

func (hs *HTTPSpanner) Decorate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if ok := CheckHeaderForMoneyTrace(request.Header); ok {
			switch {
			case hs.Tr1d1um != nil:
				htTracker, err := StarterProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				request = InjectTrackerIntoRequest(request, htTracker)
				next.ServeHTTP(response, request)
			case hs.Scytale != nil:
				htTracker, err := SubTracerProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				request = InjectTrackerIntoRequest(request, htTracker)
				next.ServeHTTP(response, request)
			case hs.Talaria.ServerDecorator != nil:
				// instead of simply decorating the next handler, the next handler is executed in a go routine.
				handler := hs.Talaria.ServerDecorator(request.Context(), hs, next, request, response)
				handler.ServeHTTP(response, request)
			}
		}
	})
}
