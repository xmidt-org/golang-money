package money

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
)

// Header keys
//
// MoneyHeaders holds a trace context of a span and are how nodes recognize
// if a Money trace needs to continue.
// MoneySpansHeaders hold the result of a finished span.
// TODO: Add Parent Span Header.
const (
	MoneyHeader      = "X-Money-Trace"
	MoneySpansHeader = "X-Money-Spans"

	// money-trace context keys
	tIDKey = "trace-id"
	pIDKey = "parent-id"
	sIDKey = "span-id"
)

// simpleResponseWriter is the core decorated http.ResponseWriter
type simpleResponseWriter struct {
	http.ResponseWriter
	code int
}

// WriteMoneySpansHeader writes a finished span's results to a simpleReponseWriter
func (rw simpleResponseWriter) WriteMoneySpansHeader(r Result) {
	var o = new(bytes.Buffer)

	h := rw.Header()

	success := rw.code < 400

	o.WriteString(r.String())
	o.WriteString(";response-code=" + strconv.Itoa(rw.code))
	o.WriteString(fmt.Sprintf(";success=" + strconv.FormatBool(success)))
	h.Set(MoneySpansHeader, o.String())
}

// checkHeaderForMoneySpan checks if a http header contains a MoneyHeader
func checkHeaderForMoneyTrace(h http.Header) bool {
	_, ok := h[MoneyHeader]
	return ok
}

// checkHeaderForMoneySpan checks if a http header contains a MoneySpansHeader
func checkHeaderForMoneySpan(h http.Header) bool {
	_, ok := h[MoneySpansHeader]
	return ok
}

// ExtractTracker extracts a tracker cotained in a given request.
func ExtractTracker(request *http.Request) (*HTTPTracker, error) {
	val := request.Context().Value(contextKeyTracker)
	t, ok := val.(*HTTPTracker)
	if !ok {
		return nil, errRequestDoesNotContainTracker
	}

	return t.HTTPTracker(), nil
}

// InjectTracker injects a tracker into a request.
func InjectTracker(request *http.Request, ht *HTTPTracker) *http.Request {
	ctx := context.WithValue(request.Context(), contextKeyTracker, ht)
	return request.WithContext(ctx)
}
