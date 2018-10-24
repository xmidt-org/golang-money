package money

import (
	"context"
	"net/http"
)

// Spanner option processes defines the transactions that occur between an options function (i.e SpanDecoder or SubTracer)
// and ServeHTTP within Decorater. They serve to compact processes so Decorate can be easily read.
//
// For every HTTPSpanner option there should be a process.
type Process func(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, Result)

// SubTracerProcess is SubTracer's derivation path.
func SubTracerProcess(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, Result, error) {
	tracker, err := hs.sb.function(r)
	if err != nil {
		return nil, Result{}, errTrackerHasNotBeenInjected
	}

	tracker, err = tracker.SubTrace(r.Context(), hs)
	if err != nil {
		return nil, Result{}, err
	}

	result, err := tracker.Finish()
	if err != nil {
		return nil, Result{}, err
	}

	htTracker := tracker.HTTPTracker()

	go func(ht *HTTPTracker) {
		hs.sb.htChannel <- htTracker
	}(htTracker)

	return htTracker, result, err
}

// StarterProcess is Starter's derivation path.
func StarterProcess(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, Result, error) {
	span, err := hs.st.function(r)
	if err != nil {
		return nil, Result{}, err
	}

	tracker := hs.Start(r.Context(), span)

	result, err := tracker.Finish()
	if err != nil {
		return nil, Result{}, err
	}

	htTracker := tracker.HTTPTracker()

	go func(ht *HTTPTracker) {
		hs.st.htChannel <- htTracker
	}(htTracker)

	return htTracker, result, nil
}

// EnderProcess is Ender's derivation path.
func EnderProcess(ctx context.Context, hs *HTTPSpanner, r *http.Request) (*HTTPTracker, Result, error) {
	tracker, err := hs.ed.function(r)
	if err != nil {
		return nil, Result{}, errTrackerHasNotBeenInjected
	}

	tracker, err = tracker.SubTrace(r.Context(), hs)
	if err != nil {
		return nil, Result{}, err
	}

	result, err := tracker.Finish()
	if err != nil {
		return nil, Result{}, err
	}

	htTracker := tracker.HTTPTracker()

	return htTracker, result, nil
}
