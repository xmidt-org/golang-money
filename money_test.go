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
	"testing"
	"time"
)

type testLogger struct {
	Logger
}

func (l *testLogger) Debug(params ...interface{}) { fmt.Println(params) }
func (l *testLogger) Error(params ...interface{}) { fmt.Println(params) }

func init() {
	tl := new(testLogger)
	SetLogger(tl)
}


func TestShowme(t *testing.T) {
	mny := new(Money)
	mny.showme()
}

func TestStringToObject(t *testing.T) {
	// string to object test

	// correct value type test
	headerval := "trace-id  =  test trace id;parent-id=  12345;%^&;span-id  =12346;span-name= TestSpanName;start-time =2015-10-09T20:30:46.782538292Z;span-duration = 3000083865;error-code=400;http-response=999;response-duration=0;foo=bar;span-success=false"
	traceId := "test trace id"
	parentId := int64(12345)
	spanId := int64(12346)
	spanName := "TestSpanName"
	startTime := "2015-10-09T20:30:46.782538292Z"
	spanDuration := int64(3000083865)
	errorCode := 400
	//httpResponse := 999
	//responseDuration := int64(0)
	//foo := "bar"
	spanSuccess := false
	
	mny := StringToObject(headerval)

	if mny.spanId != spanId {
		t.Errorf("spanId expected %v, got %v", spanId, mny.spanId)
	}
	if mny.traceId != traceId {
		t.Errorf("traceId expected v%, got %v", traceId, mny.traceId)
	}
	if mny.parentId != parentId {
		t.Errorf("parentId expected %v, got %v", parentId, mny.parentId)
	}
	if mny.spanName != spanName {
		t.Errorf("expected spanName %v, got %v", spanName, mny.spanName)
	}
	st, _ := time.Parse(time.RFC3339Nano, startTime )
	if mny.startTime != st {
		t.Errorf("expected startTime %v, got %v", st, mny.startTime.Format(time.RFC3339Nano))
	}
	if mny.spanDuration != spanDuration {
		t.Errorf("expected spanDuration %v, got %v", spanDuration, mny.spanDuration)
	}
	if mny.errorCode != errorCode {
		t.Errorf("expected errorCode %v, got %v", errorCode, mny.errorCode)
	}
	if mny.spanSuccess != spanSuccess {
		t.Errorf("expected spanSuccess %v, got %v", spanSuccess, mny.spanSuccess)
	}
	
	// incorrect value type test
	bad_headerval := "trace-id=97531;parent-id=BadParentID;%^&;span-id=BadSpanID;span-name=false;start-time=BadStartTime;span-duration=BadSpanDuration;error-code=BadErrorCode;http-response=BadHttpResponse;response-duration=BadResponseDuration;foo=87654;span-success=BadSpanSuccess"
	StringToObject(bad_headerval)
}

func TestToString(t *testing.T) {
	expect := "span-id=12346;trace-id=test trace id;parent-id=12345;span-name=WebPA-Service;start-time=2015-10-09T20:30:46.782538292Z;span-duration=3000083865;error-code=400;span-success=false"

	mny := new(Money)
	mny.spanId = int64(12346)
	mny.traceId = "test trace id"
	mny.parentId = int64(12345)
	mny.spanName = "WebPA-Service"
	st, _ := time.Parse(time.RFC3339Nano, "2015-10-09T20:30:46.782538292Z")
	mny.startTime = st
	mny.spanDuration = int64(3000083865)
	mny.errorCode = 400
	mny.spanSuccess = false

	result := mny.ToString()
	if result != expect {
		t.Errorf("Object to String failed\n expected: %v\n   result: %v", expect, result)
	}
}

func TestNewChild(t *testing.T) {
	headerval := "trace-id  =  test trace id;parent-id=  12345;span-id  =12346;span-name= WebPA-Service;start-time =2015-10-09T20:30:46.782538292Z"
	mny := NewChild(headerval, "WebPA-Service")

	if mny.spanId != int64(12347) {
		t.Errorf("expected spanId 12347, got %v", mny.spanId)
	}
	if mny.parentId != int64(12346) {
		t.Errorf("expected parentId 12346, got %v", mny.parentId)
	}
}

func TestAddResults(t *testing.T) {
	headerval := "trace-id  =  test trace id;parent-id=  12345;span-id  =12346;span-name= WebPA-Service;start-time =2015-10-09T20:30:46.782538292Z"
	mny := NewChild(headerval, "WebPA-Service")

	mny.AddResults(200, true)

	if mny.errorCode != 200 {
		t.Errorf("expected errorCode to be 200, got %d", mny.errorCode)
	}
	if mny.spanSuccess != true {
		t.Errorf("expected spanSuccess to be true, got %v", mny.spanSuccess)
	}
	if mny.spanDuration <= 0 {
		t.Errorf("expected spanDuration to be greater than 0, got %v", mny.spanDuration)
	}
}
