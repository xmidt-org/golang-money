package money

import (
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
			i: "trace-id=de305d54-75b4-431b-adb2-eb6b9e546013;parent-id=3285573610483682037;span-id=3285573610483682037",
			o: &TraceContext{
				PID: 3285573610483682037,
				SID: 3285573610483682037,
				TID: "de305d54-75b4-431b-adb2-eb6b9e546013",
			},
			e: nil,
		},

		{
			i: "parent-id=1;parent-id=3285573610483682037;span-id=3285573610483682037",
			o: nil,
			e: errBadTrace,
		},
		{
			i: "parent-id=de305d54-75b=4-431b-adb2-eb6b9e546013;parent-id=3285573610483682037;span-id=3285573610483682037",
			o: nil,
			e: errBadPair,
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
