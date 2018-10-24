package money

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
func TestNewHTTPTracker(t *testing.T) {
	var (
		c = context.Background()

		mockTC = &TraceContext{
			PID: 1,
			SID: 1,
			TID: "1",
		}

		mockSpan = &Span{
			Name: "spantest",
			TC:   mockTC,
		}

		mockHS = &HTTPSpanner{}

		ht = NewHTTPTracker(c, mockSpan, mockHS)

		expected = &HTTPTracker{
			span: &Span{
				Name: "spantest",
				TC: &TraceContext{
					PID: 1,
					SID: 1,
					TID: "1",
				},
			},
			HTTPSpanner: &HTTPSpanner{},
		}
	)

	assert.Equal(t, expected, ht)

}
*/

func TestSubTraceDone(t *testing.T) {
	var (
		c = context.Background()

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
			done: true,
		}

		mockHS = &HTTPSpanner{}
	)

	_, err := mockHT.SubTrace(c, mockHS)
	if err != errTrackerNotFinished {
		t.Fatalf("Tracker should be finished")
		return
	}
}

func TestSubTraceNotDone(t *testing.T) {
	var (
		c = context.Background()

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
			done: false,
		}

		mockHS = &HTTPSpanner{}
	)

	_, err := mockHT.SubTrace(c, mockHS)
	if err != nil {
		t.Fatalf("Tracker should not be finished")
		return
	}
}

func TestFinishDone(t *testing.T) {
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
			done: true,
		}
	)

	_, err := mockHT.Finish()
	if err != errTrackerAlreadyFinished {
		t.Fatalf("Tracker should not be done")
		return
	}
}

func TestFinishNotDone(t *testing.T) {
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
			done: false,
		}
	)

	_, err := mockHT.Finish()
	if err == errTrackerAlreadyFinished {
		t.Fatalf("Tracker should be done")
		return
	}
}

func TestStringDone(t *testing.T) {
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
			done: true,
		}
	)

	_, err := mockHT.String()
	if err != nil {
		t.Fatalf("Tracker should not be done")
		return
	}
}

func TestStringNotDone(t *testing.T) {
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
			done: false,
		}
	)

	_, err := mockHT.String()
	if err != errTrackerNotFinished {
		t.Fatalf("Tracker should be done")
		return
	}
}

func TestMapDone(t *testing.T) {
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
			done: true,
		}
	)

	_, err := mockHT.Map()
	if err != nil {
		t.Fatalf("Tracker should not be done")
		return
	}
}

func TestMapNotDone(t *testing.T) {
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
			done: false,
		}
	)

	_, err := mockHT.Map()
	if err != errTrackerNotFinished {
		t.Fatalf("Tracker should be done")
		return
	}
}

/*
// TestDecorateTransactorIfTrackerExistInRequest tests a transactor to see if a a tracker is injected in a request.
func TestDecorateTransactor(t *testing.T) {
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
			done: false,
		}
	)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := ExtractTracker(r)
		if err != nil {
			t.Fatalf("No Tracker In Context")
		}
	}))
	defer ts.Close()

	var (
		client     = ts.Client()
		transactor = client.Do
		url        = ts.URL
		r, _       = http.NewRequest("POST", url, nil)
	)

	r.RequestURI = ""

	transactor = mockHT.DecorateTransactor(transactor)

	var _, err = transactor(r)
	if err != nil {
		fmt.Print(err)
	}
}
*/

func TestStoreMoneySpans(t *testing.T) {
	var (
		client     = http.Client{}
		transactor = client.Do

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
			done: false,
		}
	)

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		o := new(bytes.Buffer)
		h := rw.Header()

		o.WriteString("test")
		h.Set(MoneySpansHeader, o.String())
	}))
	defer ts.Close()

	var url = ts.URL
	var r, _ = http.NewRequest("POST", url, nil)

	var res, _ = transactor(r)
	mockHT.storeMoneySpans(res.Header)
	if len(mockHT.spansList) < 1 {
		t.Fatalf("Trackers list did not append")
	}
}

func TestTrackerFromContext(t *testing.T) {
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
			done: false,
		}
	)

	var ctx = context.WithValue(context.Background(), contextKeyTracker, mockHT)
	var _, ok = TrackerFromContext(ctx)
	if !ok {
		t.Fatalf("No tracker in context")
	}
}

//func TestSpansList(t *testing.T) {}

//func TestSpansMaps(t *testing.T) {}

//func TestDecorateTransactor(t *testing.T) {}
