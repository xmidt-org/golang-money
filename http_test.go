package money

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainSpan(t *testing.T) {
	var handler = http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://localhost", nil)
	r.Header.Add(MoneyHeader, "parent-id=123;span-id=123;trace-id=78")

	MainSpan("testing")(handler).ServeHTTP(w, r)

	if l := len(w.Header()[MoneySpansHeader]); l != 1 {
		t.Errorf("expected headers to have only 1 span but it has %v instead", l)
	}
}

func TestRWInterceptor(t *testing.T) {
	var w = httptest.NewRecorder()

	var rw = &rwInterceptor{
		ResponseWriter: w,
		Code:           http.StatusOK,
		Body:           new(bytes.Buffer),
	}

	rw.Write([]byte("body1"))
	rw.Write([]byte("body2"))
	rw.WriteHeader(404)
	rw.WriteHeader(500)
	rw.Header().Add("header-test", "test")

	rw.Flush() //need to flush for buffers to be copied

	if v := w.Header().Get("header-test"); v != "test" {
		t.Errorf("expected header value for key 'header-test' to be 'test' but got '%v'", v)
	}

	if b := w.Body.String(); b != "body2" {
		t.Errorf("expected body to be 'body' but got '%v'", b)
	}

	if w.Code != 500 {
		t.Errorf("expected '500' but got '%v'", w.Code)
	}
}
