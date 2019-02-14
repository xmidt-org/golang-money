package money

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Time stub to simulate a clock
type Clock interface {
	Start() time.Time
	End(t time.Time) time.Duration
}

type stubClock struct{}

func (stubClock) Start() time.Time { return time.Now() }

func (stubClock) End(t time.Time) time.Duration { return time.Since(t) }

func createMockTC() *TraceContext {
	return &TraceContext{
		TID: "test-trace",
		SID: 1,
		PID: 1,
	}
}

func createMockSpan() *Span {
	var tc = createMockTC()
	var s time.Time
	var d time.Duration
	var err = fmt.Errorf("err")

	return &Span{
		Name:      "test-span",
		AppName:   "test-app",
		TC:        tc,
		Success:   true,
		Code:      1,
		StartTime: s,
		Duration:  d,
		Host:      "localhost",
		Err:       err,
	}
}

/*
func TestNewSpan(t *testing.T) {
	var tc = createMockTC()
	var expected = Span{
		Name: "test-span",
		TC:   tc,
	}

	assert.Equal(t, expected, NewSpan("test-span", tc))
}
*/

func TestMapFieldToString(t *testing.T) {
	var s = createMockSpan()
	var m map[string]interface{}

	var ourClock stubClock
	var startTime = ourClock.Start()
	time.Sleep(200000000)
	var duration = ourClock.End(startTime)

	s.StartTime, s.Duration = startTime, duration

	r, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(r, &m)

	var actual = mapFieldToString(m)

	var expected = map[string]string{
		"Name":      "test-span",
		"AppName":   "test-app",
		"TC":        "parent-id=1;span-id=1;trace-id=test-trace",
		"Success":   "true",
		"Code":      "1",
		"Duration":  fmt.Sprintf("%v"+"ns", duration.Nanoseconds()),
		"StartTime": startTime.Format("2006-01-02T15:04:05.999999999Z07:00"),
		"Err":       "Error",
		"Host":      "localhost",
	}

	// assert.Equal(t, expected, n)
	ok := reflect.DeepEqual(expected, actual)
	if !ok {
		t.Errorf("Expected:\n%v\n\nActual:\n%v", expected, actual)
	}
}

func TestMap(t *testing.T) {
	s := createMockSpan()

	var ourClock stubClock
	var startTime = ourClock.Start()
	time.Sleep(200000000)
	var duration = ourClock.End(startTime)

	s.StartTime, s.Duration = startTime, duration

	sm := s.Map()

	var expected = map[string]string{
		"Name":      "test-span",
		"AppName":   "test-app",
		"TC":        "parent-id=1;span-id=1;trace-id=test-trace",
		"Success":   "true",
		"Code":      "1",
		"Duration":  fmt.Sprintf("%v"+"ns", duration.Nanoseconds()),
		"StartTime": startTime.Format("2006-01-02T15:04:05.999999999Z07:00"),
		"Err":       "Error",
		"Host":      "localhost",
	}

	assert.Equal(t, expected, sm)
}

func FuncBuildSpanFromMap(t *testing.T) {
	var (
		duration = Duration(12)
		//times = []string{    ,    ,      }
		m = map[string]string{
			"Name":      "test-span",
			"AppName":   "test-app",
			"TC":        "parent-id=1;span-id=1;trace-id=test-trace",
			"Success":   "true",
			"Code":      "1",
			"Duration":  fmt.Sprintf("%v"+"ns", duration.Nanoseconds()),
			"StartTime": startTime.Format("2006-01-02T15:04:05.999999999Z07:00"),
			"Err":       "Error",
			"Host":      "localhost",
		}
	)

	span, err := buildSpanFromMap(m)
	if err != nil {
		t.Errorf(err)
	}
}

/*
func TestString(t *testing.T) {

	var ourClock stubClock
	var startTime = ourClock.Start()
	time.Sleep(200000000)
	var duration = ourClock.End(startTime)

	var startTimeString = startTime.Format("2006-01-02T15:04:05.999999999Z07:00")
	var durationString = fmt.Sprintf("%v"+"ns", duration.Nanoseconds())

	var s = &Span{
		Name:      "test-span",
		AppName:   "test-app",
		TC:        createMockTC(),
		Success:   true,
		Code:      1,
		StartTime: startTime,
		Duration:  duration,
		Host:      "localhost",
	}

	var expected = "span-name=test-span" +
		";app-name=test-app" +
		";span-duration=" + durationString +
		";span-success=true" +
		";parent-id=1" +
		";span-id=1" +
		";trace-id=test-trace" +
		";start-time=" + startTimeString +
		";host=localhost" +
		";response-code=1"

	assert.Equal(t, s.String(), expected)
}
*/
