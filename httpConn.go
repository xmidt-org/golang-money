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
	"net/http"
	"net/http/httptest"
	"strings"
)

// Money's version of an http.Handler to decorate
type Handler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// HandlerFunc is a function type that implements the Handler interface
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	f(rw, req)
}

// A Decorator wraps a Handler with extra behaviour
type Decorator func(Handler) Handler

// Used when decorating an http.HandlerFunc()
// spanName: will set the span-name
func AddToHandler(spanName string) Decorator {
	return func(f Handler) Handler {
		return HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get(HEADER) == "" {
				f.ServeHTTP(rw, req)
			} else {
				rec := httptest.NewRecorder()
				defer Finish(rw, req, rec)

				rw, req, rec = Start(rw, req, rec, spanName)

				f.ServeHTTP(rec, req)
			}
		})
	}
}

// Decorate decorates a Handler h with all the given Decorators, in order.
func Decorate(h Handler, ds ...(func(Handler) Handler)) Handler {
	decorated := h
	for _, decorate := range ds {
		decorated = decorate(decorated)
	}
	return decorated
}

// Removes the Money header from the http.Header
func DelAllMNYHeaders(header http.Header) http.Header {
	for h := range header {
		if strings.EqualFold(h, HEADER) {
			header.Del(h)
			return header
		}
	}

	return header
}

// Adds Money values to the http.Header
func AddAllMNYHeaders(header http.Header, allMNY []*Money) http.Header {
	for x := range allMNY {
		header.Add(HEADER, allMNY[x].ToString())
	}

	return header
}

// Updates the Money Header values for a *http.Request
func UpdateMNYHeaderReq(req *http.Request, allMNY []*Money) *http.Request {
	header := DelAllMNYHeaders(req.Header)
	header = AddAllMNYHeaders(header, allMNY)
	req.Header = header

	return req
}

// Updates the Money Header values for a http.ResponseWriter
func UpdateMNYHeaderRW(rw http.ResponseWriter, allMNY []*Money) http.ResponseWriter {
	header := DelAllMNYHeaders(rw.Header())
	header = AddAllMNYHeaders(header, allMNY)
	if mnyheader, ok := header[http.CanonicalHeaderKey(HEADER)]; ok {
		for i := 0; i < len(mnyheader); i++ {
			if i == 0 {
				rw.Header().Set(HEADER, mnyheader[i])
			} else {
				rw.Header().Add(HEADER, mnyheader[i])
			}
		}
	}

	return rw
}

// UpdateRWMoney adds any different values between the passed Money array and 
// the http.ResponseWriter to the http.ResponseWriter
func AddMoneyDiffToRW(mnys []*Money, rw http.ResponseWriter) http.ResponseWriter {
	var rwMNY []*Money
	for _, m := range rw.Header()[ http.CanonicalHeaderKey(HEADER) ] {
		rwMNY = append(rwMNY, StringToObject(m))
	}

	for _, m := range mnys {
		found := false
		
		for _, rm := range rwMNY {
			if m.spanId == rm.spanId {
				found = true
				break
			}
		}
		
		if !found {
			rw.Header().Add(HEADER, m.ToString())
		}
	}
	
	return rw
}

// Updates the Money Header values for a httptest.ResponseRecorder
func UpdateMNYHeaderRec(rec *httptest.ResponseRecorder, allMNY []*Money) *httptest.ResponseRecorder {
	header := DelAllMNYHeaders(rec.HeaderMap)
	header = AddAllMNYHeaders(header, allMNY)
	if mnyheader, ok := header[http.CanonicalHeaderKey(HEADER)]; ok {
		for i := 0; i < len(mnyheader); i++ {
			if i == 0 {
				rec.HeaderMap.Set(HEADER, mnyheader[i])
			} else {
				rec.HeaderMap.Add(HEADER, mnyheader[i])
			}
		}
	}

	return rec
}

// AddResponseMoney takes the additional money objects from the response then
// adds them into the money object array
func AddResponseDiffToMoney(resp *http.Response, allMNY []*Money) []*Money {
	var respMNY []*Money
	for _, m := range resp.Header[ http.CanonicalHeaderKey(HEADER) ] {
		respMNY = append(respMNY, StringToObject(m))
	}
	
	for _, rm := range respMNY {
		found := false
		
		for _, am := range allMNY {
			if rm.spanId == am.spanId {
				found = true
			}
		}
		
		if !found {
			allMNY = append(allMNY, rm)
		}
	}
	
	return allMNY
}

// AddResponseMoneyToRW take the different money values between a http.Response and 
// http.ResponseWriter and applies those different values from the http.Response to the
// http.ResponseWriter.
func AddResponseDiffToRW(rw http.ResponseWriter, resp *http.Response) http.ResponseWriter {
	var rwMNY []*Money
	for _, m := range rw.Header()[ http.CanonicalHeaderKey(HEADER) ] {
		rwMNY = append(rwMNY, StringToObject(m))
	}
	
	var respMNY []*Money
	for _, m := range resp.Header[ http.CanonicalHeaderKey(HEADER) ] {
		respMNY = append(respMNY, StringToObject(m))
	}
	
	for _, rp := range respMNY {
		found := false
		
		for _, rw := range rwMNY {
			if rw.spanId == rp.spanId {
				found = true
				break
			}
		}
		
		if !found {
			rw.Header().Add(HEADER, rp.ToString())
		}
	}
	
	return rw
}


func copyOfMNY(MNYs []*Money) []Money {
	var mnys []Money
	for i := 0; i < len(MNYs); i++ {
		mnys = append(mnys, *MNYs[i])
	}

	return mnys
}

// Finds the Money value for this instance
func GetCurrentHeader(MNYs []*Money) Money {
	mnys := copyOfMNY(MNYs)

	pID := mnys[0].parentId
	for len(mnys) > 1 {
		found := false
		for i := 0; i < len(mnys); i++ {
			if mnys[i].spanId == pID {
				pID = mnys[i].parentId
				found = true
				mnys = append(mnys[:i], mnys[i+1:]...)

				break
			}
		}

		if !found {
			mnys = append(mnys[:0], mnys[1:]...)
			pID = mnys[0].parentId
		}
	}

	return mnys[0]
}

// Decorator start of the money trace process
func Start(rw http.ResponseWriter, req *http.Request, rec *httptest.ResponseRecorder, spanName string) (http.ResponseWriter, *http.Request, *httptest.ResponseRecorder) {
	var allMNY []*Money
	for _, m := range req.Header[ http.CanonicalHeaderKey(HEADER) ] {
		allMNY = append(allMNY, StringToObject(m))
	}

	pMNY := GetCurrentHeader(allMNY)
	cMNY := NewChild(pMNY.ToString(), spanName)
	allMNY = append(allMNY, cMNY)

	req = UpdateMNYHeaderReq(req, allMNY)
	rw = UpdateMNYHeaderRW(rw, allMNY)
	rec = UpdateMNYHeaderRec(rec, allMNY)

	return rw, req, rec
}

// Decorator finish of the money trace process
func Finish(rw http.ResponseWriter, req *http.Request, rec *httptest.ResponseRecorder) (http.ResponseWriter, *http.Request) {
	for k, v := range rec.HeaderMap {
		for i := range v {
			if i == 0 {
				rw.Header().Set(k, v[i])
			} else {
				rw.Header().Add(k, v[i])
			}
		}
	}

	var allMNY []*Money
	for _, v := range rw.Header()[ http.CanonicalHeaderKey(HEADER) ] {
		allMNY = append(allMNY, StringToObject(v))
	}

	mny := GetCurrentHeader(allMNY)

	if rec.Code >= 400 {
		mny.AddResults(rec.Code, false)
	} else {
		mny.AddResults(rec.Code, true)
	}

	// update current mny object
	for x := range allMNY {
		if allMNY[x].spanId == mny.spanId {
			allMNY[x] = &mny
			break
		}
	}

	req = UpdateMNYHeaderReq(req, allMNY)
	rw = UpdateMNYHeaderRW(rw, allMNY)

	rw.WriteHeader(mny.errorCode)
	rw.Write([]byte(rec.Body.String()))

	return rw, req
}

// Begin will obtain the money headers from the request,
// Build a Money object array from those.
// Create a new child money object and add it to the array.
// return the Money array, and the new child Money object
func Begin(req *http.Request, spanName string) ([]*Money, *Money) {
	var allMNY []*Money
	for _, m := range req.Header[ http.CanonicalHeaderKey(HEADER) ] {
		allMNY = append(allMNY, StringToObject(m))
	}

	pMNY := GetCurrentHeader(allMNY)
	cMNY := NewChild(pMNY.ToString(), spanName)
	allMNY = append(allMNY, cMNY)

	return allMNY, cMNY
}

// End adds the results to the current/child span
// returns a http.ResponseWriter with the updated Money array
func End(rw http.ResponseWriter, allMNY []*Money, cMNY *Money, statusCode int, spanSuccess bool) http.ResponseWriter {
	cMNY.AddResults(statusCode, spanSuccess)
	
	// update current mny object
	for x := range allMNY {
		if allMNY[x].spanId == cMNY.spanId {
			allMNY[x] = cMNY
			break
		}
	}
	
	rw = UpdateMNYHeaderRW(rw, allMNY)

	return rw
}
