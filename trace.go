package money

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
	"net/http"
	"strconv"
	"strings"
)

type contextKey int

const (
	//ContextKeyMoneyTrace is used to put/get the MoneyTrace header value to/from contexts
	ContextKeyMoneyTrace contextKey = iota
)

type traceContext struct {
	TID string //Trace ID
	SID int64  //Span ID
	PID int64  //Parent ID
}

//traceContext returns a traceContext from the request headers
//for now, it is overly strict with the expected format
//we assume this is called knowing that there is a moneyHeader
func decodeTraceContext(raw string) (tc *traceContext, err error) {
	tc = new(traceContext)

	pairs := strings.Split(raw, ";")

	if len(pairs) != 3 {
		err = errors.New("expecting only three pairs in trace context header")
		return
	}

	seen := make(map[string]bool)

	for _, pair := range pairs {
		kv := strings.Split(pair, "=")

		if len(kv) != 2 {
			return nil, errors.New("expected trace context header to have pairs")
		}

		var k, v = kv[0], kv[1]

		switch {
		case k == tIDKey && !seen[tIDKey]:
			tc.TID, seen[k] = v, true

		case k == sIDKey && !seen[sIDKey]:
			var pv int64
			if pv, err = strconv.ParseInt(v, 10, 64); err != nil {
				return
			}

			tc.SID, seen[k] = pv, true
		case k == pIDKey && !seen[pIDKey]:
			var pv int64
			if pv, err = strconv.ParseInt(v, 10, 64); err != nil {
				return
			}
			tc.PID, seen[k] = pv, true

		default:
			return nil, errors.New("malformatted trace context header")
		}
	}
	return
}

func encodeTraceContext(tc *traceContext) string {
	return fmt.Sprintf("%s=%v;%s=%v;s=%v", pIDKey, tc.PID, sIDKey, tc.SID, tIDKey, tc.TID)
}

func next(current *traceContext) *traceContext {
	return &traceContext{
		PID: current.SID,
		SID: newSID(),
		TID: current.TID,
	}
}

//injectTraceContext places the given
func injectContext(moneyTraceHeaderVal string, r *http.Request) {
	var ctx = context.WithValue(r.Context(), ContextKeyMoneyTrace, moneyTraceHeaderVal)
	r.WithContext(ctx)
}

//newSID returns a random int64 between 0 and maxint64
//-1 means that an error occurred
func newSID() (r int64) {
	r = -1
	if i, e := rand.Int(rand.Reader, big.NewInt(math.MaxInt64)); e == nil {
		r = i.Int64()
	}
	return
}
