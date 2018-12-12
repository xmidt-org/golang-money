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
func WriteMoneySpansHeader(r Result, rw http.ResponseWriter, code int) {
	var o = new(bytes.Buffer)

	h := rw.Header()

	success := code < 400

	o.WriteString(r.String())
	o.WriteString(";response-code=" + strconv.Itoa(code))
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
	err, ok := h[MoneySpansHeader]
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

/*
// getMoneyTraceHeader grabs a money header from an http span.
func getMoneyTraceHeader(h http.Header) string {
	value := h.Get(MoneyHeader)
	if len(value) == 0 {
		return
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		panic(err)
	}

	return &b
}
*\

/*
// Writes all spans made under this tracker to response header
func CompleteList(ht *HTTPTracker, w http.ResponseWriter) http.ResponseWriter {
	list, err := tracker.SpansList()
	w.Header().Set("X-Money-Spans")
	for i, list := range list {
		w.Header().Add("X-Money-Spans", list[i])
	}

	return w
}
*
