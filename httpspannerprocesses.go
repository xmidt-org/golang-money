package money

import (
	"context"
	"net/http"
)

// Spanner option processes defines the transactions that occur between an options function (i.e SpanDecoder or SubTracer)
// and ServeHTTP within Decorater. They serve to compact processes so Decorate can be easily read.
//
// For every HTTPSpanner option there should be a process.
type Process func(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, error)

// StarterProcess is Starter's derivation path.
//
// Tr1d1um by default it needs to start the span due to where customer's request enter.
func StarterProcessTr1d1um(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, error) {
	span, err := hs.Tr1d1um(r)
	if err != nil {
		return nil, err
	}

	tracker := hs.start(r.Context(), span)
	return tracker.HTTPTracker(), nil
}

// SubTracerProcess is SubTracer's derivation path.
//
// Downstream nodes who subtrace include: scytale, petasos, talaria, and parados.
func SubTracerProcessScytale(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, error) {
	tracker, err := hs.Scytale(r)
	if err != nil {
		return nil, errTrackerHasNotBeenInjected
	}

	tracker, err = tracker.SubTrace(r.Context(), hs)
	if err != nil {
		return nil, err
	}

	return tracker.HTTPTracker(), nil
}

func SubTracerProcessPetasos(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, error) {
	tracker, err := hs.Scytale(r)
	if err != nil {
		return nil, errTrackerHasNotBeenInjected
	}

	tracker, err = tracker.SubTrace(r.Context(), hs)
	if err != nil {
		return nil, err
	}

	return tracker.HTTPTracker(), nil
}
