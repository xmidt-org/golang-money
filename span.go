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
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Span models all the data related to a Span
// It is a superset to what is specified in the
// spec: https://github.com/Comcast/money/wiki#what-is-captured
type Span struct {
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
	TC string

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
func NewSpan(spanName string, tc *TraceContext) *Span {
	return &Span{
		Name: spanName,
		TC:   tc,
	}
}

// Changes a maps values to type string.
func mapFieldToString(m map[string]interface{}) map[string]string {
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
				n[k] = typeInferenceTC(m[k])
			} else if k == "Err" {
				n[k] = "Error"
			}
		}
	}

	return n
}

// Map returns a string map representation of the span
func (s *Span) Map() map[string]string {
	var m map[string]interface{}

	// Receive a map of string to objects
	r, _ := json.Marshal(s)

	json.Unmarshal(r, &m)

	return mapFieldToString(m)
}

// String returns the string representation of the span
func (s *Span) String() string {
	var o = new(bytes.Buffer)

	o.WriteString("span-name=" + s.Name)
	o.WriteString(";app-name=" + s.AppName)
	o.WriteString(";span-duration=" + strconv.FormatInt(s.Duration.Nanoseconds()/1e3, 10)) //span duration in microseconds
	//	o.WriteString(";span-duration=" + strconv.FormatInt(s.Duration.Nanoseconds()/1e3, 10)) //span duration in microseconds
	o.WriteString(";span-success=" + strconv.FormatBool(s.Success))

	o.WriteString(";span-id=" + string(s.TC.SID))
	o.WriteString(";trace-id=" + s.TC.TID)
	o.WriteString(";parent-id=" + string(s.TC.PID))

	o.WriteString(fmt.Sprintf(";start-time=%v", +s.StartTime.UTC().UnixNano()))

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

// String returns the string representation of the result.
func (r *Result) String() string {
	var o = new(bytes.Buffer)

	o.WriteString("span-name=" + r.Name)
	o.WriteString(";app-name=" + r.AppName)
	o.WriteString(";host=" + r.Host)
	o.WriteString(";trace-context=" + r.TC)
	o.WriteString(";start-time=" + r.StartTime.String())
	o.WriteString(";span-duration=" + r.Duration.String())

	return o.String()
}

// BuildSpanFromMap builds a http span from a tracker
func BuildSpanFromMap(t map[string]string) (*Span, error) {
	span := new(Span)
	//span.TC = t["TC"]

	// TODO: all possible error codes or alternate route.
	if t["Code"] == "400" {
		span.Code = 400
	} else {
		span.Code = 401
	}

	if t["Success"] == "false" {
		span.Success = false
	} else {
		span.Success = true
	}

	start, err := time.Parse(t["StartTime"], "2011-01-19")
	if err != nil {
		return nil, err
	}

	span.StartTime = start
	duration, err := parseTime(t["Duration"])
	if err != nil {
		return nil, err
	}

	span.Duration = duration

	span.Name = t["Name"]
	span.AppName = t["AppName"]
	span.Err = errors.New(t["Err"])
	span.Host = t["Host"]

	return span, nil
}

func parseTime(t string) (time.Duration, error) {
	var mins, hours int
	var err error

	parts := strings.SplitN(t, ":", 2)

	switch len(parts) {
	case 1:
		mins, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
	case 2:
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}

		mins, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
	default:
		return 0, fmt.Errorf("invalid time: %s", t)
	}

	if mins > 59 || mins < 0 || hours > 23 || hours < 0 {
		return 0, fmt.Errorf("invalid time: %s", t)
	}

	return time.Duration(hours)*time.Hour + time.Duration(mins)*time.Minute, nil
}
