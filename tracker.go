package money

import (
	"context"
	"net/http"
	"os"
	"sync"
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

//Transactor is an HTTP transactor type
type Transactor func(*http.Request) (*http.Response, error)

// Tracker is the management interface for an active span.  It can be used to create
// child spans and to mark the current span as finished.
type Tracker interface {
	Spanner

	// Finish concludes this span with the given result
	Finish(Result)

	// String provides the representation of the managed span
	String() string

	DecorateTransactor(Transactor, ...SpanForwardingOptions) Transactor

	Spans() []string
}

//SpanForwardingOptions allows decoding of generic span types into that of golang money
//Today, spans should be extracted from an HTTP header and returned as a list of string-encode
//golang-money spans
type SpanForwardingOptions func(http.Header) []string

//HTTPTracker is the management type for child spans
type HTTPTracker struct {
	Spanner
	m    *sync.RWMutex
	span Span

	//spans contains the string-encoded value of all spans created under this tracker
	spans []string

	done bool //indicates whether the span associated with this tracker is finished
}

//DecorateTransactor configures a transactor to both
//inject Money Trace Context into outgoing requests
//and extract Money Spans from their responses (if any)
func (t *HTTPTracker) DecorateTransactor(transactor Transactor, options ...SpanForwardingOptions) Transactor {
	return func(r *http.Request) (resp *http.Response, e error) {
		t.m.RLock()
		r.Header.Add(MoneyHeader, EncodeTraceContext(t.span.TC))
		t.m.RUnlock()

		if resp, e = transactor(r); e == nil {
			t.m.Lock()
			defer t.m.Unlock()

			//the default behavior is always run
			for k, vs := range resp.Header {
				if k == MoneySpansHeader {
					for _, v := range vs {
						t.spans = append(t.spans, v)
					}
				}
			}

			//options allow converting different span types into money-compatible ones
			for _, o := range options {
				for _, span := range o(resp.Header) {
					t.spans = append(t.spans, span)
				}
			}
		}
		return
	}
}

//Start defines the money trace context for span s based
//on the underlying HTTPTracker span before delegating the
//start process to the Spanner
//if such underlying span has already finished, the returned
//tracker is nil
func (t *HTTPTracker) Start(ctx context.Context, s Span) (tracker Tracker) {
	t.m.RLock()
	defer t.m.RUnlock()

	if !t.done {
		s.TC = SubTrace(t.span.TC)
		tracker = t.Spanner.Start(ctx, s)
	}

	return
}

//Finish is an idempotent operation that marks the end of the underlying HTTPTracker span
func (t *HTTPTracker) Finish(r Result) {
	t.m.Lock()
	defer t.m.Unlock()

	if !t.done {
		t.span.Duration = time.Since(t.span.StartTime)
		t.span.Host, _ = os.Hostname()
		t.span.Name = r.Name
		t.span.AppName = r.AppName
		t.span.Code = r.Code
		t.span.Err = r.Err
		t.span.Success = r.Success

		t.spans = append(t.spans, t.span.String())

		t.done = true
	}
}

//String returns the string representation of the span associated with this
//HTTPTrackertracker
func (t *HTTPTracker) String() (v string) {
	t.m.RLock()
	defer t.m.RUnlock()

	v = t.span.String()
	return
}

func (t *HTTPTracker) Spans() (spans []string) {
	t.m.RLock()
	defer t.m.RUnlock()

	spans = make([]string, len(t.spans))
	copy(spans, t.spans)
	return
}

//TrackerFromContext extracts a tracker contained in the given context, if any
func TrackerFromContext(ctx context.Context) (t Tracker, ok bool) {
	t, ok = ctx.Value(contextKeyTracker).(Tracker)
	return
}
