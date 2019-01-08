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
const (
	MoneyHeader      = "X-Money-Trace"
	MoneySpansHeader = "X-Money-Spans"

	// money-trace context keys
	tIDKey = "trace-id"
	pIDKey = "parent-id"
	sIDKey = "span-id"
)

// WriteMoneySpansHeader writes a finished span's results to a responseWriter's header.
func WriteMoneySpansHeader(r Result, rw http.ResponseWriter, code interface{}) {
	var o = new(bytes.Buffer)

	h := rw.Header()

	o.WriteString(r.String())
	switch code.(type) {
	case int:
		c := code.(int)
		o.WriteString(";response-code=" + strconv.Itoa(c))
		success := c < 400
		o.WriteString(fmt.Sprintf(";success=" + strconv.FormatBool(success)))
	case int64:
		i := code.(*int64)
		c := int(*i)
		o.WriteString(";response-code=" + strconv.Itoa(c))
		success := c < 400
		o.WriteString(fmt.Sprintf(";success=" + strconv.FormatBool(success)))
	}

	h.Add(MoneySpansHeader, o.String())
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

// CheckHeaderForMoneyTrace checks if a http header contains a MoneyTrace
func CheckHeaderForMoneyTrace(h http.Header) bool {
	return checkHeaderForMoneyTrace(h)
}

func checkHeaderForMoneyTrace(h http.Header) bool {
	_, ok := h[MoneyHeader]
	return ok
}

// CheckHeaderForMoneySpan checks if a http header contains a MoneySpansHeader
func CheckHeaderForMoneySpan(h http.Header) bool {
	return checkHeaderForMoneySpan(h)
}

func checkHeaderForMoneySpan(h http.Header) bool {
	fmt.Print(2)
	_, ok := h[MoneySpansHeader]
	return ok
}
