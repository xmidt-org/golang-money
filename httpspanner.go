package money

import (
	"context"
	"net/http"
)

// HTTPSpanner implements a Spanner and its future possible options.
//
// Future created spanner options go here.
type HTTPSpanner struct {
	subtracer SubTracer
	starter   Starter
	ender     Ender
	state     bool
}

// Creates new http Spanner by extracting spanner off of request.
func NewHTTPSpanner(options HTTPSpannerOptions) *HTTPSpanner {
	hs := new(HTTPSpanner)

	return hs
}

// Start defines the start time of the input span s and returns
// a http tracker which can both start a child span using SubTrace
// as well as mark the end of a span using s
func (hs *HTTPSpanner) Start(ctx context.Context, s *Span) *HTTPTracker {
	return NewHTTPTracker(ctx, s, hs)
}

// Decorate provides an Alice-style decorator for handlers that wish to use money
func (hs *HTTPSpanner) Decorate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if ok := checkHeaderForMoneyTrace(request.Header); ok {
			switch {
			case hs.subtracer != nil:
				htTracker, err := SubTracerProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				request = InjectTracker(request, htTracker)

				next.ServeHTTP(response, request)
			case hs.starter != nil:
				htTracker, err := StarterProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				request = InjectTracker(request, htTracker)

				next.ServeHTTP(response, request)
			case hs.ender != nil:
				htTracker, err := EnderProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				request = InjectTracker(request, htTracker)

				next.ServeHTTP(response, request)
			case hs.state:
				next.ServeHTTP(response, request)
			}
		} else {
			next.ServeHTTP(response, request)
		}
	})
}
