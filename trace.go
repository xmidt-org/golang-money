package money

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//Trace Context decoding errors
var (
	errPairsCount = errors.New("expecting three pairs in trace context")
	errBadPair    = errors.New("expected trace context header to have pairs")
	errBadTrace   = errors.New("malformatted trace context header")
)

//TraceContext encapsutes all the core information of any given span
//In a single trace, the TID is the same across all spans and
//the SID and the PID is what links all spans together
type TraceContext struct {
	TID string //Trace ID
	SID int64  //Span ID
	PID int64  //Parent ID
}

//decodeTraceContext returns a TraceContext from the given value "raw"
//raw is typically taken directly from http.Request headers
//for now, it is overly strict with the expected format
//TODO: could we use regex here instead for simplicity?
func decodeTraceContext(raw string) (tc *TraceContext, err error) {
	tc = new(TraceContext)

	pairs := strings.Split(raw, ";")

	if len(pairs) != 3 {
		return nil, errPairsCount
	}

	seen := make(map[string]bool)

	for _, pair := range pairs {
		kv := strings.Split(pair, "=")

		if len(kv) != 2 {
			return nil, errBadPair
		}

		var k, v = kv[0], kv[1]

		switch {
		case k == tIDKey && !seen[k]:
			tc.TID, seen[k] = v, true

		case k == sIDKey && !seen[k]:
			var pv int64
			if pv, err = strconv.ParseInt(v, 10, 64); err != nil {
				return nil, err
			}

			tc.SID, seen[k] = pv, true
		case k == pIDKey && !seen[k]:
			var pv int64
			if pv, err = strconv.ParseInt(v, 10, 64); err != nil {
				return nil, err
			}
			tc.PID, seen[k] = pv, true

		default:
			return nil, errBadTrace
		}
	}

	return
}

//EncodeTraceContext encodes the TraceContext into a string.
//This is useful if you want to pass your trace context over an outgoing request
func EncodeTraceContext(tc *TraceContext) string {
	return fmt.Sprintf("%s=%v;%s=%v;%s=%v", pIDKey, tc.PID, sIDKey, tc.SID, tIDKey, tc.TID)
}

//SubTrace creates a child trace context for current
func SubTrace(current *TraceContext) *TraceContext {
	rand.Seed(time.Now().Unix())
	return &TraceContext{
		PID: current.SID,
		SID: rand.Int63(),
		TID: current.TID,
	}
}
