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

// SubTracerProcess is SubTracer's derivation path.
func SubTracerProcess(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, error) {
	tracker, err := hs.subtracer(r)
	if err != nil {
		return nil, errTrackerHasNotBeenInjected
	}

	tracker, err = tracker.SubTrace(r.Context(), hs)
	if err != nil {
		return nil, err
	}

	htTracker := tracker.HTTPTracker()

	return htTracker, nil
}

// StarterProcess is Starter's derivation path.
func StarterProcess(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, error) {
	span, err := hs.starter(r)
	if err != nil {
		return nil, err
	}

	tracker := hs.Start(r.Context(), span)

	htTracker := tracker.HTTPTracker()

	return htTracker, nil
}

// EnderProcess is Ender's derivation path.
func EnderProcess(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, error) {
	tracker, err := hs.ender(r)
	if err != nil {
		return nil, errTrackerHasNotBeenInjected
	}

	tracker, err = tracker.SubTrace(r.Context(), hs)
	if err != nil {
		return nil, err
	}

	htTracker := tracker.HTTPTracker()

	return htTracker, nil
}
