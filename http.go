package money

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"

	money "github.com/Comcast/golang-money"
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

// simpleResponseWriter is the core decorated http.ResponseWriter
type simpleResponseWriter struct {
	http.ResponseWriter
	code int
}

/*
//TODO:
func RunMoney(ctx context.Context, statusCode int) error {
	tracker, ok := money.TrackerFromContext(ctx)
	if ok {
		result, err := tracker.Finish()
		if err != nil {
			return err
		}

		money.WriteMoneySpansHeader(result, w, deviceResponseModel.StatusCode)
	}

	return nil
}
*/

// WriteMoneySpansHeader writes a finished span's results to a responseWriter's header.
func WriteMoneySpansHeader(r Result, rw http.ResponseWriter, code interface{}) {
	var o = new(bytes.Buffer)

	h := rw.Header()

	success := code < 400
	o.WriteString(r.String())
	switch v := code.(type) {
	case int:
		m, _ := i.(int)
		o.WriteString(";response-code=" + strconv.Itoa(code))
	case int64:
		m, _ := i.(*int64)
		o.WriteString(";response-code=" + strconv.ParseInt(&code))
	}

	o.WriteString(fmt.Sprintf(";success=" + strconv.FormatBool(success)))
	h.Add(MoneySpansHeader, o.String())
}

func CheckHeaderForMoneyTrace(h http.Header) bool {
	return checkHeaderForMoneyTrace(h)
}

// checkHeaderForMoneySpan checks if a http header contains a MoneyHeader
func checkHeaderForMoneyTrace(h http.Header) bool {
	_, ok := h[MoneyHeader]
	return ok
}

func CheckHeaderForMoneySpan(h http.Header) bool {
	return checkHeaderForMoneySpan(h)
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
	case int64:
		m, _ := i.(int64)
		return CheckDeviceResponseForMoney(m)
	}

	o.WriteString(";response-code=" + strconv.ParseInt(&code))
	o.WriteString(fmt.Sprintf(";success=" + strconv.FormatBool(success)))
	h.Add(MoneySpansHeader, o.String())
}

func CheckHeaderForMoneyTrace(h http.Header) bool {
	return checkHeaderForMoneyTrace(h)
}

// checkHeaderForMoneySpan checks if a http header contains a MoneyHeader
func checkHeaderForMoneyTrace(h http.Header) bool {
	_, ok := h[MoneyHeader]
	return ok
}

func CheckHeaderForMoneySpan(h http.Header) bool {
	return checkHeaderForMoneySpan(h)
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
	case int64:
		m, _ := i.(int64)
		return CheckDeviceResponseForMoney(m)
	}

	o.WriteString(";response-code=" + strconv.ParseInt(&code))
	o.WriteString(fmt.Sprintf(";success=" + strconv.FormatBool(success)))
	h.Add(MoneySpansHeader, o.String())
}

func CheckHeaderForMoneyTrace(h http.Header) bool {
	return checkHeaderForMoneyTrace(h)
}

// checkHeaderForMoneySpan checks if a http header contains a MoneyHeader
func checkHeaderForMoneyTrace(h http.Header) bool {
	_, ok := h[MoneyHeader]
	return ok
}

func CheckHeaderForMoneySpan(h http.Header) bool {
	return checkHeaderForMoneySpan(h)
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
