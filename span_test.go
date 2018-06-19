package money

import (
	"testing"
	"time"
)

func TestSpanString(t *testing.T) {
	startTime, e := time.Parse(time.RFC3339, "1970-01-01T00:00:01+00:00") //1 second into epoch time = 1,000,000 microseconds
	if e != nil {
		panic(e)
	}

	i := &Span{
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
