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
	moneySpansHeader = "X-Money-Spans"

	//money-trace context keys
	tIDKey = "trace-id"
	pIDKey = "parent-id"
	sIDKey = "span-id"
)

//MainSpan is a handy way to decorate your handler with money capabilities
//TODO: need to explain about how to use the money trace context values in in the request context
func MainSpan(appName string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				rw := &RWInterceptor{ResponseWriter: w, Code: http.StatusOK, Body: new(bytes.Buffer)}

				if r.Header.Get(MoneyHeader) != "" {
					var (
						tc  *traceContext
						err error
					)

					if tc, err = decodeTraceContext(r.Header.Get(MoneyHeader)); err == nil {
						r = r.WithContext(traceCxt(tc))

						var done = Start(rw.Header())

						defer func() {
							//finish this main handler span
							done(&SpanReport{
								Name:    "HTTPHandler",
								AppName: appName,
								Code:    rw.Code,
								Success: rw.Code < 400,
							})

							//flush temporary responseWriter buffer
							//into the original one
							rw.Flush()
						}()
					}
				}

				h.ServeHTTP(rw, r)
			})
	}
}

//RWInterceptor allows money trace spans to be int
type RWInterceptor struct {
	http.ResponseWriter
	Code int
	Body *bytes.Buffer
}

//Write simply saves the last array of bytes written. Such data
//is then written to the original responseWriter once Flush() is called
func (rw *RWInterceptor) Write(b []byte) (int, error) {
	return rw.Body.Write(b)
}

//WriteHeader saves the last code written to the it. Such code
//is then written to the original responseWriter once Flush() is called
func (rw *RWInterceptor) WriteHeader(code int) {
	rw.Code = code
}

//Flush transfers the temporary buffer data into the original responseWriter
func (rw *RWInterceptor) Flush() (int, error) {
	rw.ResponseWriter.WriteHeader(rw.Code)
	return rw.ResponseWriter.Write(rw.Body.Bytes())
}
