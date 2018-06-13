package money

import (
	"context"
	"os"
	"time"
)

type contextKey int

const (
	//contextKeyTracker is the key for child spans management component
	contextKeyTracker contextKey = iota
)

//Header keys
const (
	MoneyHeader      = "X-Money-Trace"
	MoneySpansHeader = "X-Money-Spans"

	//money-trace context keys
	tIDKey = "trace-id"
	pIDKey = "parent-id"
	sIDKey = "span-id"
)

// Tracker is the management interface for an active span.  It can be used to create
// child spans and to mark the current span as finished.
type Tracker interface {
	Spanner

	// Finish concludes this span with the given result
	Finish(Result)

	// String provides the representation of the managed span
	String() string
}

//HTTPTracker is the management type for child spans
type HTTPTracker struct {
	Spanner
	span Span
}

//Start defines the money trace context for span s based
//on the underlying HTTPTracker span before delegating the
//start process to the Spanner
func (t *HTTPTracker) Start(ctx context.Context, s Span) Tracker {
	s.TC = SubTrace(t.span.TC)

	return t.Spanner.Start(ctx, s)
}

//Finish marks the end of the underlying HTTPTracker span
//TODO: define what happens when this is called more than once
func (t *HTTPTracker) Finish(r Result) {
	t.span.Duration = time.Since(t.span.StartTime)
	t.span.Host, _ = os.Hostname()
	t.span.Name = r.Name
	t.span.AppName = r.AppName
	t.span.Code = r.Code
	t.span.Err = r.Err
	t.span.Success = r.Success
}

//String returns the string representation of the span associated with this
//HTTPTrackertracker
func (t *HTTPTracker) String() string {
	return t.span.String()
}

//TrackerFromContext extracts a tracker contained in the given context, if any
func TrackerFromContext(ctx context.Context) (t Tracker, ok bool) {
	t, ok = ctx.Value(contextKeyTracker).(Tracker)
	return
}
