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
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Span models all the data related to a Span
// It is a superset to what is specified in the
// spec: https://github.com/Comcast/money/wiki#what-is-captured
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

// Result models the result fields of a span.
type Result struct {
	// Name of the Span (i.e HTTPHandler)
	Name string

	// Start Time

	// Name of the application/service running the Span (i.e. Scytale in XMiDT)
	AppName string

	// StartTime

	// Code is an abstract value which is up to the span code to supply.
	// It is not necessary to enforce that this is an HTTP status code.
	// The translation into an HTTP status code should take place elsewhere.
	Code int

	// Whether or not this span is defined as "successful"
	Success bool

	Err error

	StartTime time.Time

	Duration time.Duration

	Host string
}

// NewSpan returns a new span instance.
func NewSpan(spanName string, tc *TraceContext) Span {
	return Span{
		Name: spanName,
		TC:   tc,
	}
}

type SpanMap map[string]string

// Changes a maps values to type string.
func mapFieldToString(m map[string]interface{}) SpanMap {
	n := make(map[string]string)

	for k, v := range m {
		switch v.(type) {
		case float64:
			if k == "Duration" {
				var i = int64(m[k].(float64))
				var d = time.Duration(i).Nanoseconds()
				n[k] = fmt.Sprintf("%v"+"ns", d)
			} else if k == "Code" {
				n[k] = fmt.Sprintf("%.f", m[k].(float64))
			}
		case bool:
			n[k] = strconv.FormatBool(m[k].(bool))
		case string:
			n[k] = m[k].(string)
		case map[string]interface{}:
			if k == "TC" {
				n[k] = encodeTC(m[k])
			} else if k == "Err" {
				n[k] = "Error"
			}
		}
	}

	return n
}

// Map returns a string map representation of the span
func (s *Span) Map() (SpanMap, error) {
	var m map[string]interface{}

	// Receive a map of string to objects
	r, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}

	json.Unmarshal(r, &m)

	return mapFieldToString(m), nil
}

// String returns the string representation of the span
func (s *Span) String() string {
	var o = new(bytes.Buffer)

	o.WriteString("span-name=" + s.Name)
	o.WriteString(";app-name=" + s.AppName)
	o.WriteString(";span-duration=" + fmt.Sprintf("%v"+"ns", s.Duration.Nanoseconds()))
	o.WriteString(";span-success=" + strconv.FormatBool(s.Success))
	o.WriteString(";" + encodeTraceContext(s.TC))
	o.WriteString(";start-time=" + s.StartTime.Format("2006-01-02T15:04:05.999999999Z07:00"))

	if s.Host != "" {
		o.WriteString(";host=" + s.Host)
	}

	if s.Code != 0 {
		o.WriteString(fmt.Sprintf(";response-code=%v", s.Code))
	}

	if s.Err != nil {
		o.WriteString(fmt.Sprintf(";err=%v", s.Err))
	}

	return o.String()
}
