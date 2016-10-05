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
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"
)

func getTestMoney() (*Money, *Money) {
	m1 := new(Money)
	m1.spanId = 23456
	m1.traceId = "Money For Nothing"
	m1.parentId = 6748
	m1.spanName = "Money-Span-Name-A"
	m1.startTime, _ = time.Parse(time.RFC3339Nano, "0")
	m1.spanDuration = 5
	m1.errorCode = 202
	m1.spanSuccess = true

	m2 := new(Money)
	m2.spanId = 98765
	m2.traceId = "Money For Nothing"
	m2.parentId = 23456
	m2.spanName = "Money-Span-Name-B"
	m2.startTime, _ = time.Parse(time.RFC3339Nano, "10")
	m2.spanDuration = 7
	m2.errorCode = 257
	m2.spanSuccess = true

	return m1, m2
}

func getTestMoneyArray(moneys ...*Money) []*Money {
	var mnys []*Money
	for _, money := range moneys {
		mnys = append(mnys, money)
	}
	
	return mnys
}

func makeTestServer(f Handler) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(
		f,
	)

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}
	
	return server, httpClient
}

func makeTestMoneyRequest(mnys []*Money, url string, client *http.Client) *http.Response {
	newReq, e := http.NewRequest("GET", url, nil)
	if e != nil {
		log.Error("request error: ", e)
	}
	
	for _, m := range mnys {
		newReq.Header.Add(HEADER, m.ToString())
	}
	
	resp, e := client.Do(newReq)
	if e != nil {
		log.Error("response error: ", e)
	}
	
	return resp
}

func verifyTestResponseHeaderCount(mnys []*Money, resp *http.Response) int {
	count := 0
	for _, v := range resp.Header[ http.CanonicalHeaderKey(HEADER) ] {
		for _, m := range mnys {
			if m.ToString() == v {
				count++
				break
			}
		}
	}
	
	return count
}

func verifyTestResponseHeaderCountAndCurrentSpan(mnys []*Money, resp *http.Response, spanName string) (int, string) {
	count := 0
	currentSpan := ""
	for _, v := range resp.Header[ http.CanonicalHeaderKey(HEADER) ] {
		for _, m := range mnys {
			if m.ToString() == v {
				count++
				break
			}
		}
		
		if strings.Contains(v, spanName) {
			currentSpan = v
			count++
		}
	}
	
	return count, currentSpan
}

func TestDelAllMNYHeaders(t *testing.T) {
	mnys := getTestMoneyArray(getTestMoney())

	handler := http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			headers := DelAllMNYHeaders(req.Header)
			for k, v := range headers {
				for i:=0; i < len(v); i++ {
					if i == 0 {
						rw.Header().Set(k, v[i])
					} else {
						rw.Header().Add(k, v[i])
					}
				}
			}
			rw.Write([]byte("test delete money headers"))
		},
	)
	server, client := makeTestServer(handler)
	defer server.Close()
	
	resp := makeTestMoneyRequest(mnys, server.URL, client)
	defer resp.Body.Close()

	count := verifyTestResponseHeaderCount(mnys, resp)

	if count != 0 {
		t.Errorf("Expected Money header count not correct, Got: %v, Expected: 0", count)
	}
}

func TestAddAllMNYHeaders(t *testing.T) {
	mnys := getTestMoneyArray(getTestMoney())
	
	handler := http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			AddAllMNYHeaders(rw.Header(), mnys)
			rw.Write([]byte("test adding money headers"))
		},
	)
	server, client := makeTestServer(handler)
	defer server.Close()
	
	var noMNY []*Money
	resp := makeTestMoneyRequest(noMNY, server.URL, client)
	defer resp.Body.Close()

	count := verifyTestResponseHeaderCount(mnys, resp)

	if count != len(mnys) {
		t.Errorf("Expected Money header count not correct, Got: %v, Expected: %v", count, len(mnys))
	}
}

func TestCopyOfMNY(t *testing.T) {
	mnys := getTestMoneyArray(getTestMoney())
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
	mnys := getTestMoneyArray(getTestMoney())
	var expectedId int64 = 98765

	mny := GetCurrentHeader(mnys)

	if mny.spanId != expectedId {
		t.Errorf("Incorrect header found.  expected: %v, got: %v\n", expectedId, mny.spanId)
	}
}


func TestStart(t *testing.T) {
	mnys := getTestMoneyArray(getTestMoney())

	spanName := "TEST_SPAN_NAME"
	handler := http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			rec := httptest.NewRecorder()
			rw, req, rec = Start(rw, req, rec, spanName)
			
			rw.Write([]byte("test money start function"))
		},
	)
	server, client := makeTestServer(handler)
	defer server.Close()
	
	resp := makeTestMoneyRequest(mnys, server.URL, client)
	defer resp.Body.Close()

	count, foundNewSpan := verifyTestResponseHeaderCountAndCurrentSpan(mnys, resp, spanName)

	if count != (len(mnys) + 1) {
		t.Errorf("Expected Money header count not correct, Got: %v, Expected: %v", count, (len(mnys) + 1))
	}

	if foundNewSpan == "" {
		t.Errorf("New Money span was not found.  Expected to find Money span with span name: %v", spanName)
	}
}


func TestFinish(t *testing.T) {
	m1, m2 := getTestMoney()
	mnys := getTestMoneyArray(m1, m2)

	spanName := "TEST_SPAN_NAME"

	c := new(Money)
	c.spanId = newSpanId(m1.spanId)
	c.traceId = m1.traceId
	c.parentId = m1.spanId
	c.spanName = spanName
	c.startTime = time.Now().UTC()
	mnys = append(mnys, c)

	recBody := "test money finish function"
	handler := http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			rec := httptest.NewRecorder()
			for _, m := range mnys {
				rec.Header().Add(HEADER, m.ToString())
			}
			rec.WriteHeader(222)
			rec.Write([]byte(recBody))
			Finish(rw, req, rec)
		},
	)
	server, client := makeTestServer(handler)
	defer server.Close()
	
	var noMNY []*Money
	resp := makeTestMoneyRequest(noMNY, server.URL, client)
	defer resp.Body.Close()

	count, currentSpan := verifyTestResponseHeaderCountAndCurrentSpan(mnys, resp, spanName)

	if count != len(mnys) {
		t.Errorf("Expected Money header count not correct, Got: %v, Expected: %v", count, len(mnys))
	}

	if currentSpan == "" {
		t.Errorf("New Money span was not found.  Expected to find Money span with span name: %v", spanName)
	}

	mny := StringToObject(currentSpan)
	if mny.spanDuration <= 0 || mny.errorCode != 222 || mny.spanSuccess != true {
		t.Errorf("Money header was not updated when finished. %v", currentSpan)
	}
	
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Error readying response body: ", err)
	}
	
	if string(respBody) != recBody {
		t.Error("Response body does not match ResponseWriter body before Finish function. Expected: %s, Got: %s", recBody, string(respBody))
	}
}

func TestBegin(t *testing.T) {	
	mnys := getTestMoneyArray(getTestMoney())
	
	spanName := "TEST_SPAN_NAME"
	req, err := http.NewRequest("GET", "http://www.google.com/", nil)
	if err != nil {
		log.Error("request error: ", err)
	}
	
	for _, m := range mnys {
		req.Header.Add(HEADER, m.ToString())
	}
	
	allMNY, cMNY := Begin(req, spanName)
	
	if len(allMNY) != len(mnys)+1 {
		t.Error("Request money array is not the correct size.  expected: %d, got: %d", len(mnys)+1, len(allMNY))
	}
	
	if cMNY.spanName != spanName {
		t.Error("Newly created child money span name not set correctly. expected: %s, got: %s", spanName, cMNY.spanName)
	}
	
	var found bool
	for _, am := range allMNY {
		if am == cMNY {
			found = true
			break
		}
	}
	
	if !found {
		t.Error("New child money not found in Money spans.\n%+v\n%+v", allMNY, cMNY)
	}
}

func TestAddResponseDiffToMoney(t *testing.T) {
	m1, _ := getTestMoney()
	mnys := getTestMoneyArray(m1)

	fsn1 := "FINISH-SPAN-NAME1"
	fsn2 := "FINISH-SPAN-NAME2"

	h2code := 234
	h2success := true
	handler2 := http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			h2aMNY, h2cMNY := Begin(req, fsn2)
			rw = End(rw, h2aMNY, h2cMNY, h2code, h2success)
			rw.WriteHeader(234)
			rw.Write([]byte(fsn2))
		},
	)
	
	server2, client2 := makeTestServer(handler2)
	defer server2.Close()
	
	h1code := 222
	h1success := true
	handler1 := http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			h1aMNY, h1cMNY := Begin(req, fsn1)
			rp := makeTestMoneyRequest(h1aMNY, server2.URL, client2)
			
			h1aMNY = AddResponseDiffToMoney(rp, h1aMNY)
			
			rw = End(rw, h1aMNY, h1cMNY, h1code, h1success)
			rw.WriteHeader(h1code)
			rw.Write([]byte(fsn1))
		},
	)
	
	server1, client1 := makeTestServer(handler1)
	defer server1.Close()
	
	resp := makeTestMoneyRequest(mnys, server1.URL, client1)
	defer resp.Body.Close()

	count := 0
	fsn1Money := ""
	fsn2Money := ""
	for _, v := range resp.Header[ http.CanonicalHeaderKey(HEADER) ] {
		if strings.Contains(v, fsn1) {
			fsn1Money = v
		} else if strings.Contains(v, fsn2) {
			fsn2Money = v
		}
		
		count++
	}

	if count != len(mnys) + 2 {
		t.Errorf("Expected Money header count not correct, Got: %v, Expected: %v", count, len(mnys)+2)
	}
	
	if fsn2Money == "" {
		t.Errorf("Money header was not found when finished. %v", fsn2Money)
	}
	
	if fsn1Money == "" {
		t.Errorf("Money header was not found when finished. %v", fsn1Money)
	}
}

func TestEnd(t *testing.T) {
	m1, m2 := getTestMoney()
	mnys := getTestMoneyArray(m1, m2)
	
	spanName := "TEST-END-SPANNAME"
	cMNY := NewChild(m2.ToString(), spanName)
	mnys = append(mnys, cMNY)
	
	code := 234
	success := true
	handler := http.HandlerFunc(
		func(rw http.ResponseWriter, req *http.Request) {
			End(rw, mnys, cMNY, code, success)
			
			rw.WriteHeader(code)
			rw.Write([]byte(spanName))
		},
	)
	
	server, client := makeTestServer(handler)
	defer server.Close()
	
	var noMNY []*Money
	resp := makeTestMoneyRequest(noMNY, server.URL, client)
	defer resp.Body.Close()
	
	count := 0
	var childMNY *Money
	for _, v := range resp.Header[ http.CanonicalHeaderKey(HEADER) ] {
		if strings.Contains(v, spanName) {
			childMNY = StringToObject(v)
		}
		
		count++
	}
	
	if count != len(mnys) {
		t.Error("ResponseWriter money headers count is not the correct size. expected: %d, got: %d", len(mnys), count)
	}
	
	if childMNY.errorCode != code {
		t.Error("Current span's money header error code was no updated.  expected: %d, got: %d", code, childMNY.errorCode)
	}
	
	if childMNY.spanSuccess != success {
		t.Error("Current span's money header success value was no expected.  expected: %v, got: %v", success, childMNY.spanSuccess)
	}
}

func TestAddMoneyDiffToRW(t *testing.T) {
	m1, m2 := getTestMoney()
	mnys := getTestMoneyArray(m1, m2)
	
	fsn1 := "FINISH-SPAN-NAME1"
	fsn2 := "FINISH-SPAN-NAME2"
	
	h2code := 234
	handler2 := Decorate(
		http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(h2code)
				rw.Write([]byte(fsn2))
			},
		),
	AddToHandler(fsn2))
	
	server2, client2 := makeTestServer(handler2)
	defer server2.Close()
	
	var rwMNY []*Money
	handler1 := Decorate(
		http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				var h1aMNY []*Money
				for _, m := range req.Header[ http.CanonicalHeaderKey(HEADER) ] {
					h1aMNY = append(h1aMNY, StringToObject(m))
				}
				
				rp := makeTestMoneyRequest(h1aMNY, server2.URL, client2)
				
				var rpMNY []*Money
				for _, m := range rp.Header[ http.CanonicalHeaderKey(HEADER) ] {
					rpMNY = append(rpMNY, StringToObject(m))
				}
				rw = AddMoneyDiffToRW(rpMNY, rw)
				for _, m := range rw.Header()[ http.CanonicalHeaderKey(HEADER) ] {
					rwMNY = append(rwMNY, StringToObject(m))
				}
			},
		),
	AddToHandler(fsn1))
	
	server1, client1 := makeTestServer(handler1)
	defer server1.Close()
	
	resp := makeTestMoneyRequest(mnys, server1.URL, client1)
	defer resp.Body.Close()
	
	if len(rwMNY) != len(mnys)+2 {
		t.Errorf("ResponseWriter money header was not the correct length. Expected: %d, Got: %d", len(mnys)+2, len(rwMNY))
	}
	
	var fsn1Found bool
	var fsn2Found bool
	for _, m := range rwMNY {
		if m.spanName == fsn1 {
			fsn1Found = true
		} else if m.spanName == fsn2 {
			fsn2Found = true
		}
	}
	
	if !fsn1Found {
		t.Errorf("Extra Response Money span (%s) was not found in Money object.", fsn1)
	}
	if !fsn2Found {
		t.Errorf("Extra Response Money span (%s) was not found in Money object.", fsn2)
	}
}

func TestRecursiveDepth(t *testing.T) {
	m1, _ := getTestMoney()
	mnys := getTestMoneyArray(m1)

	fsn1 := "FINISH-SPAN-NAME1"
	fsn2 := "FINISH-SPAN-NAME2"

	h2code := 234
	handler2 := Decorate(
		http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(h2code)
				rw.Write([]byte(fsn2))
			},
		),
	AddToHandler(fsn2))
	
	server2, client2 := makeTestServer(handler2)
	defer server2.Close()
	
	h1code := 222
	handler1 := Decorate(
		http.HandlerFunc(
			func(rw http.ResponseWriter, req *http.Request) {
				var h1aMNY []*Money
				for _, m := range req.Header[ http.CanonicalHeaderKey(HEADER) ] {
					h1aMNY = append(h1aMNY, StringToObject(m))
				}
				
				rp := makeTestMoneyRequest(h1aMNY, server2.URL, client2)
				rw = AddResponseDiffToRW(rw, rp)
				
				rw.WriteHeader(h1code)
				rw.Write([]byte(fsn1))
			},
		),
	AddToHandler(fsn1))
	
	server1, client1 := makeTestServer(handler1)
	defer server1.Close()
	
	resp := makeTestMoneyRequest(mnys, server1.URL, client1)
	defer resp.Body.Close()

	count := 0
	fsn1Money := ""
	fsn2Money := ""
	for _, v := range resp.Header[ http.CanonicalHeaderKey(HEADER) ] {
		if strings.Contains(v, fsn1) {
			fsn1Money = v
		} else if strings.Contains(v, fsn2) {
			fsn2Money = v
		}
		
		count++
	}

	if count != len(mnys) + 2 {
		t.Errorf("Expected Money header count not correct, Got: %v, Expected: %v", count, len(mnys)+2)
	}
	
	if fsn2Money == "" {
		t.Errorf("Money header was not found when finished. %v", fsn2)
	}
	
	if fsn1Money == "" {
		t.Errorf("Money header was not found when finished. %v", fsn1)
	}
}
