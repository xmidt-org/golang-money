package money

import (
	"net/http"
	"testing"
	"time"
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

func TestSpanString(t *testing.T) {
	startTime, e := time.Parse(time.RFC3339, "1970-01-01T00:00:01+00:00") //1 second into epoch time = 1,000,000 microseconds
	if e != nil {
		panic(e)
	}

	i := &span{
		Name:    "test-span",
		AppName: "test-app",
		TC: &TraceContext{
			TID: "test-trace",
			SID: 1,
			PID: 1,
		},
		Success:   true,
		Code:      200,
		StartTime: startTime,
		Duration:  time.Second,
		Host:      "localhost",
	}

	var expected = "span-name=test-span;app-name=test-app;span-duration=1000000;span-success=true;span-id=1;trace-id=test-trace;parent-id=1;start-time=1000000;host=localhost;http-response-code=200"

	if i.String() != expected {
		t.Errorf("expected '%s' but got '%s", expected, i.String())
	}
}
