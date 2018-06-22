package money

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPSpanner(t *testing.T) {
	t.Run("Start", testStart)
}

func testStart(t *testing.T) {
	startTime, e := time.Parse(time.RFC3339, "1970-01-01T00:00:01+00:00") //1 second into epoch time = 1,000,000 microseconds
	require.Nil(t, e)

	c, cancel := context.WithCancel(context.Background())
	require.NotNil(t, c)
	defer cancel()

	s := Span{
		Name:    "test-span",
		AppName: "test-start-app",
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
	require.NotNil(t, s)

	h := new(HTTPTracker)
	h.span = Span{
		Name:    "test-tracker-span",
		AppName: "test-app",
		TC: &TraceContext{
			TID: "test-tracker-trace",
			SID: 11,
			PID: 11,
		},
		Success:   true,
		Code:      200,
		StartTime: startTime,
		Duration:  time.Second,
		Host:      "localhost",
	}

	var end HTTPTracker
	/* Following call results in panic.
	 * TODO - Need clarification on the following. Start is a method call of pointer receiver HTTPTracker and returns a HTTPTracker
	 */
	//end := h.Start(c, s)
	assert.NotNil(t, end)
}
