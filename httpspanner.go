package money

import (
	"context"
	"net/http"
)

// HTTPSpanner implements a Spanner and its future possible options.
//
// Future created spanner options go here.
type HTTPSpanner struct {
	sb SubTracerContainer
	st StarterContainer
	ed EnderContainer
	s  bool
}

// Creates new http Spanner by extracting spanner off of request.
// see https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
// for more details.
func NewHTTPSpanner(options ...HTTPSpannerOptions) *HTTPSpanner {
	hs := new(HTTPSpanner)

	for _, o := range options {
		o(hs)
	}

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
			case hs.sb.function != nil:
				htTracker, err := SubTracerProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				request = InjectTracker(request, htTracker)

				next.ServeHTTP(response, request)
			case hs.st.function != nil:
				htTracker, err := StarterProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				request = InjectTracker(request, htTracker)

				next.ServeHTTP(response, request)
			case hs.ed.function != nil:
				htTracker, err := EnderProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				request = InjectTracker(request, htTracker)

				next.ServeHTTP(response, request)
			case hs.s:
				next.ServeHTTP(response, request)
			}
		} else {
			next.ServeHTTP(response, request)
		}
	})
}
