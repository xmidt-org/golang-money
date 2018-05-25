package money

import (
	"net/http"
	"testing"
)

func TestSpanStart(t *testing.T) {
	h := make(http.Header)
	End := Start(h)

	sr := &SpanReport{
		Name:    "testSpan",
		AppName: "testService",
		TC: &TraceContext{
			PID: 123,
			SID: 123,
			TID: "trace-xyz",
		},
		Success: true,
		Code:    200,
	}

	End(sr)

	if len(h[MoneySpansHeader]) != 1 {
		t.Error("expected a money span header to be present")
	}
}

func TestSpanStrong(t *testing.T) {
	//TODO:
}
