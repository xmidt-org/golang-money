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

func AddToHandler(spanName string) Decorator {
	return func(f Handler) Handler {
		return HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Header.Get(HEADER) == "" {
				f.ServeHTTP(rw, req)
			} else {
				rec := httptest.NewRecorder()
				defer Finish(rw, req, rec)

				rw, req = Start(rw, req, spanName)

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

func DelAllMNYHeaders(rw http.ResponseWriter) http.ResponseWriter {
	for k := range rw.Header() {
		if strings.ToLower(k) == strings.ToLower(HEADER) {
			rw.Header().Del(k)
		}
	}

	return rw
}

func AddAllMNYHeaders(rw http.ResponseWriter, allMNY []*Money) http.ResponseWriter {
	for x := range allMNY {
		rw.Header().Add(HEADER, allMNY[x].ToString())
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

func Start(rw http.ResponseWriter, req *http.Request, spanName string) (http.ResponseWriter, *http.Request) {
	var allMNY []*Money
	for k, v := range req.Header {
		if strings.ToLower(k) == strings.ToLower(HEADER) {
			for i := range v {
				allMNY = append(allMNY, StringToObject(v[i]))
			}
		}
	}

	pMNY := GetCurrentHeader(allMNY)
	cMNY := NewChild(pMNY.ToString(), spanName)
	allMNY = append(allMNY, cMNY)

	rw = DelAllMNYHeaders(rw)
	rw = AddAllMNYHeaders(rw, allMNY)

	return rw, req
}

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
	for k, v := range rw.Header() {
		if strings.ToLower(k) == strings.ToLower(HEADER) {
			for i := range v {
				allMNY = append(allMNY, StringToObject(v[i]))
			}
		}
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

	rw = DelAllMNYHeaders(rw)
	rw = AddAllMNYHeaders(rw, allMNY)

	rw.WriteHeader(mny.errorCode)
	rw.Write([]byte(rec.Body.String()))

	return rw, req
}
