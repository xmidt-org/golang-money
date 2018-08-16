package money

import (
	"net/http"
	"strconv"
)

func WriteCompleteHeaders(r *Result, h http.Header) {

	// Add adds headers as pictured as a stack - MoneySpansHeader
	// will appear first.
	//
	// This is not coupled with the WriteHeader method within spanner.go
	// because it allows users to easily write to headers in their own application using
	// the finished result details.
	h.Add("span-duration", r.Duration.String())
	h.Add("start-time", r.StartTime.String())
	h.Add("span-success", strconv.FormatBool(r.Success))
	h.Add("app-name", r.AppName)
	h.Add("span-name", r.Name)
	h.Add(sIDKey, string(r.sIDKey))
	h.Add(pIDKey, string(r.pIDKey))
	h.Add(tIDKey, r.tIDKey)
	h.Add(MoneySpansHeader, "X-MoneySpans")
}
