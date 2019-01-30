package money

import (
	"context"
	"errors"
	"net/http"
	"os"
	"sync"
	"time"
)

// Tracking errors
var (
	errTrackerNotFinished = errors.New("Tracker should have not yet finished")
	// As of now the only method that changes a http trackers field, done, is Finish
	errTrackerAlreadyFinished        = errors.New("Tracker needs to be finished to utilize this function")
	errRequestDoesNotContainTracker  = errors.New("Request does not contain tracker")
	errResponseDoesNotContainTracker = errors.New("Response does not contain tracker")
	errTrackerHasNotBeenInjected     = errors.New("Tracker has not been injected")
)

type contextKey int

const (
	//contextKeyTracker is the key for child spans management component
	contextKeyTracker contextKey = iota
)

// Tracker is the management interface for an active span.
type Tracker interface {
	// Create a new child span given a request's context and prior span's trace context
	SubTrace(context.Context, HTTPSpanner) (*HTTPTracker, error)

	// Finish concludes this span with the given result
	Finish() (Result, error)

	// String provides the string representation of the managed span
	String() (string, error)

	// Map provides the representation the map span representation of the managed span.
	Map() (map[string]string, error)

	// Spans returns a list of string-encoded spans in a map fashion.
	SpansMap() ([]map[string]string, error)

	// Spans returns a list of string-encoded Money spans that have been created under this tracker
	SpansList() ([]string, error)

	// HTTPTracker returns a http tracker object.
	HTTPTracker() *HTTPTracker
}

// HTTPTracker is the management type for child spans
type HTTPTracker struct {
	*HTTPSpanner
	m sync.RWMutex

	span *Span

	// spans contains the string-encoded value of all spans created under this tracker
	// should be modifiable by multiple goroutines
	spansList []string

	// spansMaps contains span maps of all spans created under this tracker
	spansMaps []map[string]string

	// indicates whether the span associated with this tracker is finished
	done bool
}

// NewHTTPTracker defines the start time of the input span s and returns
// a HTTPTracker.  It is utilized by a HTTPTracker Start method.
func NewHTTPTracker(ctx context.Context, s *Span, sp *HTTPSpanner) *HTTPTracker {
	s.StartTime = time.Now()

	return &HTTPTracker{
		span:        s,
		HTTPSpanner: sp,
	}
}

// SubTrace starts a child span from the given span s.   A child span's paramount attribute
// is it's trace context, TC,  due to a span's span-id/SID uniqueness.
func (t *HTTPTracker) SubTrace(ctx context.Context, sp *HTTPSpanner) (*HTTPTracker, error) {
	t.m.RLock()
	defer t.m.RUnlock()

	if !t.done {
		return &HTTPTracker{
			span: &Span{
				TC:        doSubTrace(t.span.TC),
				StartTime: time.Now(),
			},
			HTTPSpanner: sp,
		}, nil
	}

	return nil, errTrackerNotFinished
}

// Finish is an idempotent operation that marks the end of the underlying HTTPTracker by adding appending spans maps and spans list as well as
// returning a Span's contents as a Result object.
func (t *HTTPTracker) Finish() (Result, error) {
	t.m.Lock()
	defer t.m.Unlock()

	if !t.done {
		t.span.Duration = time.Since(t.span.StartTime)
		//	t.span.Code = //TODO get span code
		t.span.Success = t.span.Code < 400

		t.spansList = append(t.spansList, t.span.String())
		t.spansMaps = append(t.spansMaps, t.span.Map())
		t.done = true

		return Result{
			Name:      t.span.Name,
			TC:        encodeTraceContext(t.span.TC),
			AppName:   os.Args[0],
			Code:      t.span.Code,
			Success:   t.span.Success,
			Err:       t.span.Err,
			StartTime: t.span.StartTime,
			Duration:  t.span.Duration,
			Host:      t.span.Host,
		}, nil
	} else {
		return Result{}, errTrackerAlreadyFinished
	}
}

// finish the tracker
// tr1d1um finishes the tracker & loops through list of spanMaps and writes the
// contents of each map in spanMaps to http.ResponseWriter.
// I need a function that turns the contents of a map[string]string to a single concatenated string.

// String returns the string representation of the span associated with this
// HTTPTracker once such span has finished, zero value otherwise
func (t *HTTPTracker) String() (v string, err error) {
	t.m.Lock()
	defer t.m.Unlock()

	if t.done {
		v = t.span.String()
		return v, nil
	}

	return "", errTrackerNotFinished
}

// SpansMap returns the map representation of the span associated with this
// HTTPTracker once such span has finished, zero value otherwise.
func (t *HTTPTracker) Map() (v map[string]string, err error) {
	t.m.Lock()
	defer t.m.Unlock()

	if t.done {
		v = t.span.Map()
		return v, nil
	}

	return nil, errTrackerNotFinished
}

// Spans returns the list of string-encoded spans under this tracker
// once the main span under the tracker is finished, zero value otherwise.
func (t *HTTPTracker) SpansList() (spansList []string, err error) {
	t.m.RLock()
	defer t.m.RUnlock()

	if t.done {
		spansList = make([]string, len(t.spansList))
		copy(spansList, t.spansList)
		return spansList, nil
	}

	return nil, errTrackerNotFinished
}

// SpansMaps returns the list of span map objects under this tracker
// once the main span under the tracker is finished, zero value otherwise.
func (t *HTTPTracker) SpansMap() (spansMap []map[string]string, err error) {
	t.m.RLock()
	defer t.m.RUnlock()

	if t.done {
		spansMaps := make([]map[string]string, len(t.spansMaps))
		copy(spansMaps, t.spansMaps)
		return spansMaps, nil
	}

	return nil, errTrackerNotFinished
}

// storeMoneySpans adds a responses money spans to a HTTPTracker objects spansList
func (t *HTTPTracker) storeMoneySpans(h http.Header) {
	for k, vs := range h {
		if k == MoneySpansHeader {
			t.spansList = append(t.spansList, vs...)
			return
		}
	}

	return
}

/*
// storeMoneyMaps updates the spansMaps field of a tracker
func (t *HTTPTracker) storeMoneyMaps() {
	for k, vs := range h {
		if k == MoneySpansHeader {
			t.Span.String()
			return
		}
	}

	return
}
*/

// Returns a HTTPTracker object.
func (t *HTTPTracker) HTTPTracker() *HTTPTracker {
	return t
}

func (t *HTTPTracker) Switch() {
	t.m.RLock()
	defer t.m.RUnlock()

	if t.done {
		t.done = false
	} else if !t.done {
		t.done = true
	}
}

// UpdateMaps updates the spans maps.  Used when handling device spans.
func (t *HTTPTracker) UpdateMaps(maps []map[string]string) {
	t.m.RLock()
	defer t.m.RUnlock()

	t.spansMaps[len(t.spansMaps)-1] = maps[len(maps)-1]
}

// UpdateMaps gets the spans maps.  Used when handling talaria and device span bridge.
func (t *HTTPTracker) GetMaps() []map[string]string {
	return t.spansMaps
}

// UpdateSpan updaets the span of a device.  Its specifically used to add spans to a device's tracker.
func (t *HTTPTracker) UpdateSpan(m map[string]string) {
	t.m.RLock()
	defer t.m.RUnlock()

	t.span, _ = BuildSpanFromMap(m)
}

// TrackerFromContext extracts a tracker contained in a given context.
func TrackerFromContext(ctx context.Context) (*HTTPTracker, bool) {
	t, ok := ctx.Value(contextKeyTracker).(*HTTPTracker)
	return t.HTTPTracker(), ok
}
