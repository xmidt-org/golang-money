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
	"strconv"
	"time"
)

//Span models all the data related to a Span
//It is a superset to what is specified in the
//spec: https://github.com/Comcast/money/wiki#what-is-captured
type Span struct {
	//the user gives us these values
	Name    string
	AppName string
	TC      *TraceContext
	Success bool
	Code    int
	Err     error

	//we deduce these values
	StartTime time.Time
	Duration  time.Duration
	Host      string
}

// Result models the result fields of a span.  The zero value of this struct
// indicates a successful span execution.
type Result struct {
	//Name of the Span (i.e HTTPHandler)
	Name string

	//Name of the application/service running the Span (i.e. Scytale in XMiDT)
	AppName string

	//whether or not this span is defined as "successful"
	Success bool

	// Code is an abstract value which is up to the span code to supply.
	// It is not necessary to enforce that this is an HTTP status code.
	// The translation into an HTTP status code should take place elsewhere.
	Code int

	// Err is just the error reported by the span
	Err error
}

//String() returns the string representation of the span
//TODO: update with new variables added to Span like err
func (s *Span) String() string {
	var o = new(bytes.Buffer)

	o.WriteString("span-name=" + s.Name)
	o.WriteString(";app-name=" + s.AppName)
	o.WriteString(";span-duration=" + strconv.FormatInt(s.Duration.Nanoseconds()/1e3, 10)) //span duration in microseconds
	o.WriteString(";span-success=" + strconv.FormatBool(s.Success))

	o.WriteString(";span-id=" + strconv.FormatInt(int64(s.TC.SID), 10))
	o.WriteString(";trace-id=" + s.TC.TID)
	o.WriteString(";parent-id=" + strconv.FormatInt(int64(s.TC.PID), 10))

	o.WriteString(fmt.Sprintf(";start-time=%v", s.StartTime.UTC().UnixNano()/1e3)) //UTC time since epoch in microseconds

	if s.Host != "" {
		o.WriteString(";host=" + s.Host)
	}

	if s.Code != 0 {
		o.WriteString(fmt.Sprintf(";http-response-code=%v", s.Code))
	}

	if s.Err != nil {
		o.WriteString(fmt.Sprintf(";err=%v", s.Code))
	}

	return o.String()
}
