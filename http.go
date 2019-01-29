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

func CheckHeaderForMoneyTrace(h http.Header) bool {
	_, ok := h[MoneyHeader]
	return ok
}

func CheckHeaderForMoneySpan(h http.Header) bool {
	_, ok := h[MoneySpansHeader]
	return ok
}

// checkForTrackerInContext checks if a context contains a tracker
func CheckForTrackerInContext(ctx context.Context) bool {
	_, ok := ctx.Value(contextKeyTracker).(*HTTPTracker)
	return ok
}

// MapsToStringResult returns the a string of all the traces created under this tracker
func MapsToStringResult(m []map[string]string) string {
	var o = new(bytes.Buffer)
	for _, v := range m {
		for k, x := range v {
			o.WriteString(k + "=" + x + ";")
		}
	}

	return o.String()
}

/*
// NewMoneyResponseHeader clears the reponse header and injects money
func NewMoneyResponse(r *http.Response,          ) *http.Response {
	moneyResponse := response
	for k := range moneyResponse.Header {
		delete(m, k)
	}

	ctx := context.WithValue(request.Context(), contextKeyTracker, ht)
	return request.WithContext(ctx)
}
*/
