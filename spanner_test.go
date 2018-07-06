package money

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHTTPSpanner(t *testing.T) {
	t.Run("DI", testNewHTTPSpannerNil)
	t.Run("Start", testStart)
	t.Run("DecorationNoMoneyContext", testDecorate)
	t.Run("DecorationMoneyContext", testDecorateWithMoney)
}

func testNewHTTPSpannerNil(t *testing.T) {
	var spanner *HTTPSpanner
	if spanner.Decorate("test", nil) != nil {
		t.Error("Decoration should leave handler unchanged")
	}
}

func testStart(t *testing.T) {
	var spanner = NewHTTPSpanner()
	if spanner.Start(context.Background(), Span{}) == nil {
		t.Error("was expecting a non-nil response")
	}
}

func testDecorate(t *testing.T) {
	var spanner = NewHTTPSpanner()

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, ok := TrackerFromContext(r.Context())
			if ok {
				t.Error("Tracker should not have been injected")
			}
		})
	decorated := spanner.Decorate("test", handler)
	decorated.ServeHTTP(nil, httptest.NewRequest("GET", "localhost:9090/test", nil))
}

func testDecorateWithMoney(t *testing.T) {
	var spanner = NewHTTPSpanner()

	handler := http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			_, ok := TrackerFromContext(r.Context())
			if !ok {
				t.Error("Expected tracker to be present")
			}
		})
	decorated := spanner.Decorate("test", handler)
	inputRequest := httptest.NewRequest("GET", "localhost:9090/test", nil)
	inputRequest.Header.Add(MoneyHeader, "trace-id=abc;parent-id=1;span-id=1")
	decorated.ServeHTTP(nil, inputRequest)
}

//create a test that simply finishes the tracker that was started
