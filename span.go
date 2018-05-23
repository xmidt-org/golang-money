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
	"os"
	"strconv"
	"time"
)

//Span represents a single money trace span
type Span struct {
	ID        int64
	Name      string
	Code      int //HTTP Response Code
	TraceID   string
	ParentID  int64
	StartTime time.Time
	Duration  int64
	Success   bool
	Host      string
}

//Start begins the life of a span
//returns a func whose calling represents the end of the span
//w is expected to be the response writer to the calling system
//r is expected to be the incoming request from the calling system
func Start(o *HandlerOptions, w, r *http.Request) func() {
	var start = time.Now()

	return func() {
		//TODO: here we would calculate the duration of the span
		// and place the span into the responsewriter headers

		//we want to flush the
	}
}

//String() returns the string representation of the current span
//TODO:
//-revise to see if better approach could be taken: maybe through field tags
//-does this align with the documentation?
func (s *Span) String() string {
	var o = new(bytes.Buffer)

	o.WriteString("span-name=" + s.Name)
	o.WriteString(";span-duration=" + strconv.FormatInt(int64(s.Duration), 10))
	o.WriteString(";span-success=" + strconv.FormatBool(s.Success))

	o.WriteString(";span-id=" + strconv.FormatInt(int64(s.ID), 10))
	o.WriteString(";trace-id=" + s.TraceID)
	o.WriteString(";parent-id=" + strconv.FormatInt(int64(s.ParentID), 10))
	o.WriteString(";start-time=" + s.StartTime.Format(time.RFC3339Nano))

	hostname, _ := os.Hostname()
	o.WriteString(";host=" + hostname)

	return o.String()
}
