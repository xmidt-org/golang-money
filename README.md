Golang-money
===

Distributed Tracing using Go.  

This is the golang implementation of [Money](https://github.com/Comcast/money)

[![Build Status](https://travis-ci.org/Comcast/golang-money.svg?branch=master)](https://travis-ci.org/Comcast/golang-money) 
[![codecov.io](http://codecov.io/github/Comcast/golang-money/coverage.svg?branch=master)](http://codecov.io/github/Comcast/golang-money?branch=master) 
[![Go Report Card](https://goreportcard.com/badge/github.com/Comcast/golang-money)](https://goreportcard.com/report/github.com/Comcast/golang-money) 


## Requirements 

[Go](http://golang.org) 1.9 or newer.

## Getting Started

Send a curl with the following header:

```
"curl http://localhost:12345/ -i -H "X-Money-Trace: trace-id=de305d54-75b4-431b-adb2-eb6b9e546013;parent-id=3285573610483682037;span-id=3285573610483682037"
```

Or post a request from your application:

In Go:

```
req, err := http.NewRequest("GET", "http://localhost:12345/", nil)
if err != nil {
	// handle err
}
req.Header.Set("X-Moneytrace", "trace-id=de305d54-75b4-431b-adb2-eb6b9e546013;parent-id=3285573610483682037;span-id=3285573610483682037")

resp, err := http.DefaultClient.Do(req)
if err != nil {
	// handle err
}
defer resp.Body.Close()
```

A header that contains a trace context {span-id, trace-id, parent-id} and the header, X-MoneyTrace, is a MoneyHeader. MoneyHeaders are required for the Money Decorator
to continue a trace.

Whats returned from a host that processes a X-MoneyTrace is a response that contains a X-MoneySpans header.  This header contains details from the host that just received your http request:


|Span Data   |Description                                      |
|------------|-------------------------------------------------|
|spanId      |current span's identifier                        |
|traceId     |name for the trace                               |
|parentId    |current span's parent identifier                 |
|AppName     |name of applciation/service who created the span |
|spanName    |current span's name                              |
|startTime   |current span's start time                        |
|spanDuration|current span's duration time                     |
|errorCode   |current span's error code                        |
|spanSuccess |Was the current span successful                  |
|host        |host who created this span                       | 

All X-MoneySpans responses are recorded in the HTTPTracker's object list as a concatented string of spand data.  

```
X-Money-Spans: span-id=-460900382554701468;trace-id=de305d54-75b4-431b-adb2-eb6b9e546013;parent-id=-460900382554701468;
span-name=spanName;app-name=appName;start-time=1412550594494;
;http-response-code=500;span-success=false;span-duration=120004.0;host=myHost; 
```

#### Using Money into your application:  

Create a HTTPSpanner using NewHTTPSpanner function.  This function is the plane where specific HTTSpanner options are chosen from.

Currently there are four options to choose from: 

(1) Starter: When a node/machine is responsible for starting a trace and injecting the first HTTPTracker. 
(2) SubTracer: When a node/machine is responsible for subtracing and the HTTPTracker object needs to be forwarded to a new http request. 
(3) End:  When a node/machine is responsible for subtracing but does not need to forward the HTTPTracker to a transactor. 
(4) Off: Turns off the decorator. 

OPTIONS:
 StarterON()
 SubtracerON()
 EnderON()
 Off

```
myHTTPSpanner := NewHTTPSpanner(OPTION)
```

### 2

#### Examples: 

Decorate your httpSpanner using a Alice style decorator as well as your transactor. 
```
chain := alice.New(s.Decorate).Then(yourSpecialHandler)
httptracker.DecorateTransactor(yourSpecialTransactor)
```

Example01: Using StarterON()

```
func yourSpecialHandler1(w http.ResponseWriter, r *http.Request) {
	msg := "Request Received 1\n"
	w.Write([]byte(msg))
}

func main() {
	h1 := http.HandlerFunc(yourSpecialHandler)

	// Create a HTTPSpanner object.
	s := money.NewHTTPSpanner(StarterON())

	// Decorate HTTPSpanner. Note that Decorate turns a handler.
	chain := alice.New(s.Decorate).Then(yourSpecialHandler)

	if err := http.ListenAndServe(":12345", chain); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
```

Example02: Using SubtracerON() 

````
package main

import (
	"log"
	"net/http"

	"github.com/Comcast/golang-money"
	"github.com/justinas/alice"
)

func responseHandler1(w http.ResponseWriter, r *http.Request) {
	msg := "Request Received 1\n"
	w.Write([]byte(msg))
}

// Your transactor should utilize the same methods functions in the same
// fashion as below.
func transactor(r *http.Request) (*http.Response, error) {
	client := &http.Client{}

	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func main() {
	// Create a HTTPSpanner object.
	htChannel := make(chan *money.HTTPTracker)
	h1 := http.HandlerFunc(responseHandler1)

	// Create a HTTPSpanner object with SpanDecoderWithChannelON
	// This spanner takes in a channel that is later used to receive
	// httpTrackers from httpTrackers wrangled by handlers.
	s := spanner.NewHTTPSpanner(spanner.SpanDecoderWithChannelON(htChannel))

	chain := alice.New(s.Decorate).Then(h1)

	if err := http.ListenAndServe(":12345", chain); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	// Receive httpTracker from handler.
	httpTracker := <-htChannel

	result := httpTracker.Finish()

	// Forward httpTracker to destination described above, localhost:123456.
	transactor = httpTracker.DecorateTransactor(transactor)
    
	req, err := http.NewRequest("GET", ":12346", nil)
	if err != nil {
		return nil, err
	}

    response, err := transactor(req)
    if err != nil {
        return nil, err
    }

    // response will need to be utilized  
 }
````
#### Creating new HTTPSpanner options

(1) Make a new field in the HTTPSpanner struct and a httpspanner container for options that require specific parts.  
(2) Create a option in httpspanneroptions.go
(3) Create a new process in httpspannerprocess.go
(4) Implement the newly created option and process into spanner's Decorate. 

#### How is Money data interpreted?

Currently for our use case money spans are read from Splunk logs.  Future implementations could have Money spans forwarded to a microservice that help understand system anamolies by using HTTPTracker's lists and map of list fields. 

