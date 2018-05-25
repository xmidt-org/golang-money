// Copyright 2016 Comcast Cable Communications Management, LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// End Copyright

package money

import (
	"bytes"
	"net/http"
)

//Header keys
const (
	MoneyHeader      = "X-Money-Trace"
	MoneySpansHeader = "X-Money-Spans"

	//money-trace context keys
	tIDKey = "trace-id"
	pIDKey = "parent-id"
	sIDKey = "span-id"
)

//MainSpan is the way to decorate your handler with money capabilities
//It also places money trace context information in the incoming request's contexts
//This is done so you can pass such trace context information to subsequent systems using money as well.
//Important: this decorator changes the behavior of the incoming responseWriter into the handler
//by allowing multiple body and code writes and keeping only the last writes. This allows the decorator
//to send back money span headers
func MainSpan(appName string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {

				rw := &rwInterceptor{ResponseWriter: w, Code: http.StatusOK, Body: new(bytes.Buffer)}
				defer rw.Flush()

				if r.Header.Get(MoneyHeader) != "" {
					var (
						tc  *TraceContext
						err error
					)

					if tc, err = decodeTraceContext(r.Header.Get(MoneyHeader)); err == nil {
						r = r.WithContext(traceCxt(tc))

						var done = Start(rw.Header())

						//finish this main handler span
						defer done(&SpanReport{
							Name:    "HTTPHandler",
							AppName: appName,
							Code:    rw.Code,
							Success: rw.Code < 400,
							TC:      tc,
						})
					}
				}

				h.ServeHTTP(rw, r)
			})
	}
}

//ForwardMoneySpanHeaders copies all money span headers from an input header to a target
//This is useful to pass money span headers received from remote systems to your system's
//response writer headers
func ForwardMoneySpanHeaders(from http.Header, to http.Header) {
	for k, vs := range from {
		if k == MoneySpansHeader {
			for _, v := range vs {
				to.Add(k, v)
			}
		}
	}
}

//rwInterceptor allows temporary buffering of
//body and code for an original responseWriter
type rwInterceptor struct {
	http.ResponseWriter
	Code int
	Body *bytes.Buffer
}

//Write simply saves the last array of bytes written. Such data
//is then written to the original responseWriter once Flush() is called
func (rw *rwInterceptor) Write(b []byte) (int, error) {
	rw.Body.Reset() //starting fresh
	return rw.Body.Write(b)
}

//WriteHeader saves the last code written to the it. Such code
//is then written to the original responseWriter once Flush() is called
func (rw *rwInterceptor) WriteHeader(code int) {
	rw.Code = code
}

//Flush transfers the temporary buffer data into the original responseWriter
func (rw *rwInterceptor) Flush() (int, error) {
	rw.ResponseWriter.WriteHeader(rw.Code)
	return rw.ResponseWriter.Write(rw.Body.Bytes())
}
