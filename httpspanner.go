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
//
// Currently this does not take into the account MoneySpan headers.
func (hs *HTTPSpanner) Decorate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if ok := checkHeaderForMoneyTrace(request.Header); ok {
			switch {
			case hs.sb.function != nil:
				htTracker, result, err := SubTracerProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				rw := &simpleResponseWriter{
					code:           http.StatusOK,
					ResponseWriter: response,
				}

				rw.WriteMoneySpansHeader(result)

				next.ServeHTTP(rw, InjectTracker(request, htTracker))
			case hs.st.function != nil:
				htTracker, result, err := StarterProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				rw := simpleResponseWriter{
					code:           http.StatusOK,
					ResponseWriter: response,
				}

				rw.WriteMoneySpansHeader(result)

				next.ServeHTTP(rw, InjectTracker(request, htTracker))
			case hs.ed.function != nil:
				htTracker, result, err := EnderProcess(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				rw := simpleResponseWriter{
					code:           http.StatusOK,
					ResponseWriter: response,
				}

				rw.WriteMoneySpansHeader(result)

				next.ServeHTTP(rw, InjectTracker(request, htTracker))
			case hs.s:
				next.ServeHTTP(response, request)
			}
		} else {
			next.ServeHTTP(response, request)
		}
	})
}
