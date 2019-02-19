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
	Petasos SubTracer
	Talaria SubTracer
}

func NewHTTPSpanner(options HTTPSpannerOptions) *HTTPSpanner {
	hs := new(HTTPSpanner)
	if options == nil {
		return hs
	}

	options(hs)

	return hs
}

func (hs *HTTPSpanner) Decorate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		if ok := CheckHeaderForMoneyTrace(request.Header); ok {
			switch {
			case hs.Tr1d1um != nil:
				tracker, err := StarterProcessTr1d1um(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				// write the ending span result to headers after all other http.Handlers have finished executing.
				defer func() {
					if err := tracker.Finish(); err == nil {
						maps, err := tracker.SpansMap()
						if err != nil {
							return
						}

						response.Header().Set(MoneySpansHeader, mapsToStringResult(maps))
					}
				}()

				request = InjectTrackerIntoRequest(request, tracker)
				next.ServeHTTP(response, request)
			case hs.Scytale != nil:
				tracker, err := SubTracerProcessScytale(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				request = InjectTrackerIntoRequest(request, tracker)
				next.ServeHTTP(response, request)
			case hs.Petasos != nil:
				tracker, err := SubTracerProcessPetasos(request.Context(), hs, request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				request = InjectTrackerIntoRequest(request, tracker)
				next.ServeHTTP(response, request)
			case hs.Talaria != nil:
				tracker, err := hs.Talaria(request)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				tracker, err = tracker.SubTrace(request.Context(), hs)
				if err != nil {
					next.ServeHTTP(response, request)
				}

				next.ServeHTTP(response, request)
			}
		} else {
			next.ServeHTTP(response, request)
		}
	})
}

// Start defines the start time of the input span s and returns
// a http tracker which can both start a child span using SubTrace
// as well as mark the end of a span using s
func (hs *HTTPSpanner) start(ctx context.Context, s *Span) *HTTPTracker {
	return NewHTTPTracker(ctx, s, hs)
}
