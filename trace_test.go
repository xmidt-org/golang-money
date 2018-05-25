package money

import (
	"context"
	"reflect"
	"testing"
)

func TestDecodeTraceContext(t *testing.T) {
	tests := []struct {
		name string
		i    string
		o    *TraceContext
		e    error
	}{
		{
			i: "",
			o: nil,
			e: errPairsCount,
		},
		{
			name: "ideal",
			i:    "trace-id=de305d54-75b4-431b-adb2-eb6b9e546013;parent-id=3285573610483682037;span-id=3285573610483682037",
			o: &TraceContext{
				PID: 3285573610483682037,
				SID: 3285573610483682037,
				TID: "de305d54-75b4-431b-adb2-eb6b9e546013",
			},
			e: nil,
		},

		{
			name: "duplicateEntries",
			i:    "parent-id=1;parent-id=3285573610483682037;span-id=3285573610483682037",
			o:    nil,
			e:    errBadTrace,
		},
		{
			name: "badPair",
			i:    "parent-id=de305d54-75b=4-431b-adb2-eb6b9e546013;parent-id=3285573610483682037;span-id=3285573610483682037",
			o:    nil,
			e:    errBadPair,
		},

		{
			name: "empty",
			i:    "",
			o:    nil,
			e:    errBadTrace,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualO, actualE := decodeTraceContext(test.i)
			if actualE != test.e || !reflect.DeepEqual(actualO, test.o) {
				t.Errorf("I was expecting '%v' '%v'", test.e, test.o)
				t.Errorf("but got '%v' '%v'", actualE, actualO)
			}
		})
	}
}

func TestDecodeTraceContextOtherCases(t *testing.T) {
	t.Run("NonIntSID", func(t *testing.T) {
		i := "trace-id=de305d54-75b4-431b-adb2-eb6b9e546013;parent-id=3285573610483682037;span-id=NaN"
		tc, e := decodeTraceContext(i)

		if tc != nil || e == nil {
			t.Errorf("expected tc to be nil and error to be non-nil but got '%v' and '%v'", tc, e)
		}
	})

	t.Run("NonIntPID", func(t *testing.T) {
		i := "trace-id=de305d54-75b4-431b-adb2-eb6b9e546013;parent-id=NaN;span-id=123"
		tc, e := decodeTraceContext(i)

		if tc != nil || e == nil {
			t.Errorf("expected tc to be nil and error to be non-nil but got '%v' and '%v'", tc, e)
		}
	})
}

func TestEncodeTraceContext(t *testing.T) {
	in := &TraceContext{
		PID: 1,
		TID: "one",
		SID: 1,
	}
	actual, expected := EncodeTraceContext(in), "parent-id=1;span-id=1;trace-id=one"

	if actual != expected {
		t.Errorf("Expected %v but got %v", expected, actual)
	}
}

func TestSubtrace(t *testing.T) {
	current := &TraceContext{
		TID: "123",
		SID: 1,
	}
	st := SubTrace(current)

	if st.PID != current.SID {
		t.Errorf("Expected pid to be %v but got %v", current.PID, st.PID)
	}

	if st.SID == 0 {
		t.Error("Expected sid to be defined")
	}

	if st.TID != current.TID {
		t.Errorf("Expected tid to be %v but got %v", current.TID, st.TID)
	}
}

func TestPassThroughContext(t *testing.T) {

	tests := []struct {
		name string
		in   context.Context
		eVal string
		eOk  bool
	}{
		{
			name: "noValue",
			in:   context.TODO(),
			eVal: "",
			eOk:  false,
		},

		{
			name: "ideal",
			in:   context.WithValue(context.Background(), contextKeyMoneyTraceHeader, "testVal"),
			eVal: "testVal",
			eOk:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aVal, aOk := PassThroughTraceContext(test.in)

			if aVal != test.eVal || aOk != test.eOk {
				t.Errorf("expected '%s' and '%v' but got '%s' and '%v'", test.eVal, test.eOk, aVal, aOk)
			}
		})
	}
}

func TestMainSpanChildContext(t *testing.T) {

	testTraceCtx := &TraceContext{
		TID: "test-trace",
		SID: 123,
		PID: 123,
	}

	tests := []struct {
		name string
		in   context.Context
		eVal *TraceContext
		eOk  bool
	}{
		{
			name: "noValue",
			in:   context.TODO(),
			eVal: nil,
			eOk:  false,
		},

		{
			name: "ideal",
			in:   context.WithValue(context.Background(), contextKeyChildMoneyTrace, testTraceCtx),
			eVal: testTraceCtx,
			eOk:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			aVal, aOk := MainSpanChildContext(test.in)

			if aOk != test.eOk || !reflect.DeepEqual(aVal, test.eVal) {
				t.Errorf("expected '%v' and '%v' but got '%v' and '%v'", test.eVal, test.eOk, aVal, aOk)
			}
		})
	}
}
