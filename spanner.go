package money

import (
	"context"
	"net/http"
	"time"
)

// Spanner acts as the factory for spans for all downstream code.
type Spanner interface {
	Start(context.Context, Span) Tracker
}

//SpanDecoder decodes an X-Money span off a request
type SpanDecoder func(*http.Request) (Span, error)

// HTTPSpanner implements Spanner and is the root factory
// for HTTP spans
type HTTPSpanner struct {
	sd SpanDecoder
}

//Start defines the start time of the input span s and returns
//a tracker object which can both start a child span for s as
//well as mark the end of span s
func (hs *HTTPSpanner) Start(ctx context.Context, s Span) Tracker {
	s.StartTime = time.Now()

	return &HTTPTracker{
		span:    s,
		Spanner: hs,
	}
}

//Decorate provides an Alice-style decorator for handlers
//that wish to use money
func (hs *HTTPSpanner) Decorate(next http.Handler) http.Handler {
	if hs == nil {
		// allow DI of nil values to shut off money trace
		return next
	}

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if span, err := hs.sd(request); err == nil {
			tracker := hs.Start(request.Context(), span)

			ctx := context.WithValue(request.Context(), contextKeyTracker, tracker)

			s := simpleResponseWriter{
				code:           http.StatusOK,
				ResponseWriter: response,
			}

			next.ServeHTTP(
				s,
				request.WithContext(ctx))

			//TODO: there is work to be done to capture information on the span that wraps the entire
			//ServeHTTP.
			tracker.Finish(Result{
				Code:    s.code,
				Success: s.code < 400,
			})

		} else {
			next.ServeHTTP(response, request)
		}
	})
}

type HTTPSpannerOptions func(*HTTPSpanner)

func New(options ...HTTPSpannerOptions) (spanner *HTTPSpanner) {
	spanner = new(HTTPSpanner)

	//define the default behavior which is a simple
	//extraction of money trace context off the headers
	spanner.sd = func(r *http.Request) (s Span, err error) {
		var tc *TraceContext
		if tc, err = decodeTraceContext(r.Header.Get(MoneyHeader)); err == nil {
			s = Span{
				TC: tc,
			}
		}
		return
	}

	for _, o := range options {
		o(spanner)
	}
	return
}

// simpleResponseWriter is the core decorated http.ResponseWriter.
type simpleResponseWriter struct {
	http.ResponseWriter
	code int
}

func (rw simpleResponseWriter) WriteHeader(code int) {
	rw.code = code
	rw.WriteHeader(code)
}
