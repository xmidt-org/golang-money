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
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type span struct {
	//the user gives us these values
	Name    string
	AppName string
	TC      *TraceContext
	Success bool
	Code    int

	//we deduce these values
	StartTime time.Time
	Duration  time.Duration
	Host      string
}

//SpanReport contains money trace parameters used for reporting back on spans
type SpanReport struct {
	//Name of the Span (i.e HTTPHandler)
	Name string

	//Name of the application/service running the Span (i.e. Scytale in XMiDT)
	AppName string

	//money trace context for this span
	TC *TraceContext

	//whether or not this span is defined as "successful"
	Success bool

	//Optional: status code for the operation (i.e 200 for an HTTP)
	Code int
}

//Start begins the life of a span
//returns a func whose calling represents the end of the span
func Start(h http.Header) func(*SpanReport) {
	var start = time.Now()

	return func(sr *SpanReport) {
		var span = span{
			Name:      sr.Name,
			AppName:   sr.AppName,
			TC:        sr.TC,
			Success:   sr.Success,
			Code:      sr.Code,
			StartTime: start,
			Duration:  time.Since(start),
		}

		h.Add(MoneySpansHeader, span.String())
	}
}

//String() returns the string representation of the current span
//TODO:
//-revise to see if better approach could be taken: maybe through field tags
//-does this align with the documentation?
func (s *span) String() string {
	var o = new(bytes.Buffer)

	o.WriteString("span-name=" + s.Name)
	o.WriteString(";span-duration=" + strconv.FormatInt(int64(s.Duration), 10))
	o.WriteString(";span-success=" + strconv.FormatBool(s.Success))

	o.WriteString(";span-id=" + strconv.FormatInt(int64(s.TC.SID), 10))
	o.WriteString(";trace-id=" + s.TC.TID)
	o.WriteString(";parent-id=" + strconv.FormatInt(int64(s.TC.PID), 10))
	o.WriteString(";start-time=" + s.StartTime.Format(time.RFC3339Nano))

	if s.Code != 0 {
		o.WriteString(fmt.Sprintf(";http-response-code=%v", s.Code))
	}

	if hostname, err := os.Hostname(); err == nil {
		o.WriteString(";host=" + hostname)
	}

	return o.String()
}
