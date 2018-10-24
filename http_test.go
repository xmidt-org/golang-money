package money

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func NewMockResult() Result {
	var (
		startTime = time.Now()
		duration  = time.Since(startTime)
	)

	return Result{
		Name:      "test",
		TC:        "test",
		AppName:   "test",
		Code:      0,
		Success:   false,
		Err:       errors.New("test"),
		StartTime: startTime,
		Duration:  duration,
		Host:      "test",
	}
}

func TestCheckHeaderForMoneyTrace(t *testing.T) {
	var r = httptest.NewRequest("GET", "localhost:9090/test", nil)
	r.Header.Set(MoneyHeader, "test")

	var ok = checkHeaderForMoneyTrace(r.Header)
	if !ok {
		t.Fatalf("should contain money header")
	}

}

func TestCheckHeaderForMoneySpan(t *testing.T) {
	var r = httptest.NewRequest("GET", "localhost:9090/test", nil)
	r.Header.Set(MoneySpansHeader, "test")

	var ok = checkHeaderForMoneySpan(r.Header)
	if !ok {
		t.Fatalf("should contain money span header")
	}
}

func TestWriteMoneySpanHeader(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		var mockResponseWriter = simpleResponseWriter{
			code:           http.StatusOK,
			ResponseWriter: w,
		}

		var Result = NewMockResult()
		mockResponseWriter.WriteMoneySpansHeader(Result)
	}

	var req = httptest.NewRequest("GET", "http://example.com/foo", nil)
	var w = httptest.NewRecorder()
	handler(w, req)

	var resp = w.Result()
	var ok = checkHeaderForMoneySpan(resp.Header)
	spew.Dump(resp.Header)
	if !ok {
		spew.Dump(resp.Header)
		t.Fatalf("should contain Money span header")
	}
}

func TestExtractTracker(t *testing.T) {
	var (
		mockTC = &TraceContext{
			PID: 1,
			SID: 1,
			TID: "1",
		}

		mockSpan = &Span{
			Name: "spantest",
			TC:   mockTC,
		}

		mockHT = &HTTPTracker{
			span: mockSpan,
		}
	)

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, err := ExtractTracker(r)
			if err != nil {
				t.Error("Expected tracker to be present")
			}

		})

	r := httptest.NewRequest("GET", "localhost:9090/test", nil)

	var ctx = context.WithValue(r.Context(), contextKeyTracker, mockHT)

	r = r.WithContext(ctx)

	handler.ServeHTTP(nil, r)
}

func TestInjectTracker(t *testing.T) {
	var (
		mockTC = &TraceContext{
			PID: 1,
			SID: 1,
			TID: "1",
		}

		mockSpan = &Span{
			Name: "spantest",
			TC:   mockTC,
		}

		mockHT = &HTTPTracker{
			span: mockSpan,
		}
	)

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, err := ExtractTracker(r)
			if err != nil {
				t.Error("Expected tracker to be present")
			}

		})

	r := httptest.NewRequest("GET", "localhost:9090/test", nil)

	handler.ServeHTTP(nil, InjectTracker(r, mockHT))
}
