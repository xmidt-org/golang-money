We have adopted OTEL instead.

# Money

## Distributed Tracing using Go
This is the Go implementation of [Money](https://github.com/Comcast/money)

[![Build Status](https://github.com/xmidt-org/golang-money/actions/workflows/ci.yml/badge.svg)](https://github.com/xmidt-org/golang-money/actions/workflows/ci.yml)
[![codecov.io](http://codecov.io/github/xmidt-org/golang-money/coverage.svg?branch=main)](http://codecov.io/github/xmidt-org/golang-money?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/xmidt-org/golang-money)](https://goreportcard.com/report/github.com/xmidt-org/golang-money)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=xmidt-org_golang-money&metric=alert_status)](https://sonarcloud.io/dashboard?id=xmidt-org_golang-money)
[![Apache V2 License](http://img.shields.io/badge/license-Apache%20V2-blue.svg)](https://github.com/xmidt-org/golang-money/blob/main/LICENSE)
[![GitHub Release](https://img.shields.io/github/release/xmidt-org/golang-money.svg)](CHANGELOG.md)
[![GoDoc](https://pkg.go.dev/badge/github.com/xmidt-org/golang-money)](https://pkg.go.dev/github.com/xmidt-org/golang-money)


### A Money header looks like the following
```
Money: trace-id=YourTraceId;parent-id=12345;span-id=12346;span-name=YourSpanName;start-time=2016-02-15T20:30:46.782538292Z;span-duration=3000083865;error-code=200;span-success=true
```

|Span Data   |Description                     |
|------------|--------------------------------|
|spanId      |current span's identifier       |
|traceId     |name for the trace              |
|parentId    |current span's parent identifier|
|spanName    |current span's name             |
|startTime   |current span's start time       |
|spanDuration|current span's duration time    |
|errorCode   |current span's error code       |
|spanSuccess |Was the current span successful |

### Functionality to handle the Money header can be added in two ways
#### 1. Decorate you're handlers with the Money handler
```
Money.Decorate( [http.Handler], Money.AddToHandler( [spanName] ))
```

#### 2. Use the Money Begin and End functions by adding them to your http.Handler

### Start server and make a request that includes a Money header

The basics to start a Money trace are a trace id name and starting span id number.
```
Money:trace-id=YourTraceId;span-id=12345;
```
