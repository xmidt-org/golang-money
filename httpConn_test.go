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
	"reflect"
	"testing"
	"time"
	//"net/http"
	//"strings"
)

func setupMoneyArray() (*Money, *Money) {
	m1 := new(Money)
	m1.spanId = 23456
	m1.traceId = "money for nothing"
	m1.parentId = 6748
	m1.spanName = "spanner"
	m1.startTime, _ = time.Parse(time.RFC3339Nano, "0")
	m1.spanDuration = 5
	m1.errorCode = 202
	m1.spanSuccess = true

	m2 := new(Money)
	m2.spanId = 98765
	m2.traceId = "money for nothing"
	m2.parentId = 23456
	m2.spanName = "namest"
	m2.startTime, _ = time.Parse(time.RFC3339Nano, "10")
	m2.spanDuration = 7
	m2.errorCode = 257
	m2.spanSuccess = true

	return m1, m2
}

func TestDelAllMNYHeaders(t *testing.T) {
	/*
		for k, _ := range rw.Header() {
			if strings.ToLower(k) == strings.ToLower(HEADER) {
				rw.Header().Del(k)
			}
		}

		return rw
	*/
	/*
		m1, m2 := setupMoneyArray()
		var mnys []*Money
		mnys = append(mnys, m1)
		mnys = append(mnys, m2)

		rw := *new(http.ResponseWriter)

		for _, mny := range mnys {
			rw.Header().Add("Money", mny.ToString())
		}

		c := 0
		for k, _ := range rw.Header() {
			if strings.ToLower(k) == strings.ToLower(HEADER) {
				c++
			}
		}

		rw = DelAllMNYHeaders(rw)
		d1 := len(rw.Header()) - len(mnys)
		d2 := len(rw.Header()) - c
		if d1 != d2 {
			t.Errorf("Money header removal failed")
		}
	*/
}

func TestAddAllMNYHeaders(t *testing.T) {
	/*
		for x := range allMNY {
			rw.Header().Add(HEADER, allMNY[x].ToString())
		}

		return rw
	*/
}

func TestMyGrandParentIs(t *testing.T) {
	m1, m2 := setupMoneyArray()
	var mnys []Money
	mnys = append(mnys, *m1)
	mnys = append(mnys, *m2)
	var id int64 = 98765
	var expectedId int64 = 23456

	mnys, mnyId, found := myGrandParentIs(mnys, id)

	if !found {
		t.Error("Descendent not found")
	}
	if mnyId != expectedId {
		t.Error("Wrong descendent found")
	}
	if len(mnys) > 1 {
		t.Error("Descendent not removed")
	}

	for _, m := range mnys {
		if m.spanId == id {
			t.Error("Current id was not removed correctly")
		}
	}

}

func TestCopyOfMNY(t *testing.T) {
	m1, m2 := setupMoneyArray()
	var mnys []*Money
	mnys = append(mnys, m1)
	mnys = append(mnys, m2)
	newMnys := copyOfMNY(mnys)

	if reflect.TypeOf(newMnys) == reflect.TypeOf(mnys) {
		t.Errorf("copy []Money is not a copy")
	}

	for _, m := range mnys {
		found := false
		for _, n := range newMnys {
			if m.spanId == n.spanId {
				found = true
				break
			}
		}

		if !found {
			t.Error("Money values are not the same")
			break
		}
	}
}

func TestGetCurrentHeader(t *testing.T) {
	m1, m2 := setupMoneyArray()
	var mnys []*Money
	mnys = append(mnys, m1)
	mnys = append(mnys, m2)
	var expectedId int64 = 98765

	mny := GetCurrentHeader(mnys)

	if mny.spanId != expectedId {
		t.Errorf("Incorrect header found.  expected: %v, got: %v\n", expectedId, mny.spanId)
	}
}
