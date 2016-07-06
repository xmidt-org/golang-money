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
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	HEADER = "X-MoneyTrace"
)

type Money struct {
	spanId       int64
	traceId      string
	parentId     int64
	spanName     string
	startTime    time.Time
	spanDuration int64
	errorCode    int
	spanSuccess  bool
}

func (mny *Money) showme() {
	fmt.Printf("spanId:       %d\n", mny.spanId)
	fmt.Printf("traceId:      %s\n", mny.traceId)
	fmt.Printf("parentId:     %d\n", mny.parentId)
	fmt.Printf("spanName:     %s\n", mny.spanName)
	fmt.Printf("startTime:    %v\n", mny.startTime)
	fmt.Printf("spanDuration: %v\n", mny.spanDuration)
	fmt.Printf("errorCode:    %d\n", mny.errorCode)
	fmt.Printf("spanSuccess:  %v\n", mny.spanSuccess)
}

func (m *Money) setSpanId(val int64)        {m.spanId = val}
func (m *Money) setTraceId(val string)      {m.traceId = val}
func (m *Money) setParentId(val int64)      {m.parentId = val}
func (m *Money) setSpanName(val string)     {m.spanName = val}
func (m *Money) setStartTime(val time.Time) {m.startTime = val}
func (m *Money) setSpanDuration(val int64)  {m.spanDuration = val}
func (m *Money) setErrorCode(val int)       {m.errorCode = val}
func (m *Money) setSpanSuccess(val bool)    {m.spanSuccess = val}

func StringToObject(headerValue string) *Money {
	var mny Money

	pairs := strings.Split(strings.TrimSuffix(headerValue, ";"), ";")
	for p := 0; p < len(pairs); p++ {
		kv := strings.SplitN(pairs[p], "=", 2)

		if len(kv) == 2 {
			key := strings.TrimSpace(kv[0])
			val := strings.TrimSpace(kv[1])

			switch key {
			case "span-id":
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					log.Error("Unable to convert Money span-id string value to int64: %s", val)
				}
				mny.spanId = int64(i)

			case "trace-id":
				mny.traceId = val

			case "parent-id":
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					log.Error("Unable to convert Money parent-id string value to int64: %s", val)
				}
				mny.parentId = int64(i)

			case "span-name":
				mny.spanName = val

			case "start-time":
				t, err := time.Parse(time.RFC3339Nano, val)
				if err != nil {
					log.Error("Unable to convert Money start-time string value to time: %s", val)
				}
				mny.startTime = t

			case "span-duration":
				i, err := strconv.ParseInt(val, 10, 64)
				if err != nil {
					log.Error("Unable to convert Money span-duration string value to int64: %s", val)
				}
				mny.spanDuration = i

			case "error-code":
				i, err := strconv.ParseInt(val, 10, 0)
				if err != nil {
					log.Error("Unable to convert Money error-code string value to int: %s", val)
				}
				mny.errorCode = int(i)

			case "span-success":
				b, err := strconv.ParseBool(val)
				if err != nil {
					log.Error("Unable to convert Money span-success string value to bool: %s", val)
				}
				mny.spanSuccess = b

			case "http-response":
				log.Debug("Money key ignored: %s", key)

			case "response-duration":
				log.Debug("Money key ignored: %s", key)

			default:
				log.Debug("Money key unknown: %s", key)
			}
		} else {
			log.Error("Money header, bad key/value pair: %v.  Header: %v", kv, headerValue)
		}
	}

	return &mny
}

func (mny *Money) ToString() string {
	var result string

	result = "span-id=" + strconv.FormatInt(int64(mny.spanId), 10)
	result += ";trace-id=" + mny.traceId
	result += ";parent-id=" + strconv.FormatInt(int64(mny.parentId), 10)
	result += ";span-name=" + mny.spanName
	result += ";start-time=" + mny.startTime.Format(time.RFC3339Nano)
	result += ";span-duration=" + strconv.FormatInt(int64(mny.spanDuration), 10)
	result += ";error-code=" + strconv.FormatInt(int64(mny.errorCode), 10)
	result += ";span-success=" + strconv.FormatBool(mny.spanSuccess)

	return result
}

func newSpanId(parentid int64) int64 {
	return parentid + 1
}

func (mny *Money) AddResults(errorCode int, spanSuccess bool) *Money {
	mny.errorCode = errorCode
	mny.spanSuccess = spanSuccess
	mny.spanDuration = int64(time.Since(mny.startTime) / time.Microsecond)

	return mny
}

func NewChild(parentHeader, spanName string) *Money {
	pMNY := StringToObject(parentHeader)
	cMNY := new(Money)

	cMNY.spanId = newSpanId(pMNY.spanId)
	cMNY.traceId = pMNY.traceId
	cMNY.parentId = pMNY.spanId
	cMNY.spanName = spanName
	cMNY.startTime = time.Now().UTC()

	return cMNY
}
