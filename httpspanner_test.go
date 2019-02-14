package money

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockResponse struct {
	http.ResponseWriter
	headers    http.Header
	statusCode int
}

func NewMockResponse() *MockResponse {
	return &MockResponse{
		headers: make(http.Header),
	}
}

func TestNewHTTPSpanner(t *testing.T) {
	t.Run("Start", testStart)
	//	t.Run("DecorationNoOptions", testDecorateNoOptions)
	//	t.Run("DecorationSpanDecoderON", testDecorateSpanDecoderON)
	//	t.Run("DecorationSubTracerON", testDecorateSubTracerON)
}

func testNewHTTPSpannerNil(t *testing.T) {
	var spanner *HTTPSpanner
	if spanner.Decorate(nil) != nil {
		t.Error("Decoration should leave handler unchanged")
	}
}

func testStart(t *testing.T) {
	var spanner = NewHTTPSpanner(nil)
	if spanner.start(context.Background(), &Span{}) == nil {
		t.Error("was expecting a non-nil response")
	}
}

// Test non-edge spans
func TestDecorateSubTracerON(t *testing.T) {
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
		spanner = NewHTTPSpanner(ScytaleON())
	)

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, ok := TrackerFromContext(r.Context())
			if !ok {
				t.Error("Expected tracker to be present")
			}

		})

	decorated := spanner.Decorate(handler)
	inputRequest := httptest.NewRequest("GET", "localhost:9090/test", nil)
	inputRequest.Header.Add(MoneyHeader, "trace-id=abc;parent-id=1;span-id=1")

	var r = httptest.NewRecorder()
	decorated.ServeHTTP(r, InjectTrackerIntoRequest(inputRequest, mockHT))

}

// Tests edge server spans
func TestDecorateStarterON(t *testing.T) {
	var (
		spanner = NewHTTPSpanner(Tr1d1umON())
	)

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, ok := TrackerFromContext(r.Context())
			if !ok {
				t.Error("Expected tracker to be present")
			}

		})

	decorated := spanner.Decorate(handler)
	inputRequest := httptest.NewRequest("GET", "localhost:9090/test", nil)
	inputRequest.Header.Add(MoneyHeader, "trace-id=abc;parent-id=1;span-id=1")
	var r = httptest.NewRecorder()
	decorated.ServeHTTP(r, inputRequest)
}
