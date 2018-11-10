package money

import (
	"fmt"
	"net/http"
)

// Option types for a HTTPSpanner used to quickly change a HTTPSpanners state.
type HTTPSpannerOptions func(*HTTPSpanner)

// A list of a HTTPSpanners different states.
type SubTracer func(*http.Request) (*HTTPTracker, error)
type Starter func(*http.Request) (*Span, error)
type Ender func(*http.Request) (*HTTPTracker, error)
type Off func(*http.Request)

// subTracer extracts a tracker from a request's context. Its used
// in the subtracer option when there already exists a tracker to subtrace from.
func subTracer(r *http.Request) (*HTTPTracker, error) {
	return ExtractTracker(r)
}

// starter decodes money headers off a request. Its used
// in the starter option when an http tracker has yet to be created.
func starter(r *http.Request) (*Span, error) {
	tc, err := decodeTraceContext(r.Header.Get(MoneyHeader))
	if err != nil {
		fmt.Print(err)
	}

	s := NewSpan("HTTPSpan", tc)

	return s, err
}

// SubTracerON is an option to use the decorator as a subtracer.
func SubTracerON() HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.st = StarterContainer{}
		hs.sb.function = subTracer
		hs.ed = EnderContainer{}
		hs.s = false
	}
}

// StarterON is an option to use the decorator as a starter.
func StarterON() HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.st.function = starter
		hs.sb = SubTracerContainer{}
		hs.ed = EnderContainer{}
		hs.s = false
	}
}

// End is an option to use the decorator as a Ender
func EnderON() HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.st = StarterContainer{}
		hs.sb = SubTracerContainer{}
		hs.st = StarterContainer{}
		hs.ed.function = subTracer
		hs.s = false
	}
}

// RedirectorON is an option to use the decorator as a redirector.
func RedirectorON() HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.st = StarterContainer{}
		hs.sb = SubTracerContainer{}
		hs.ed = EnderContainer{}
		hs.s = false
	}
}

// SpannerOFF turns off all of HTTPSpanner's possible states.
// TODO: this could removed by changing the logic in the httpspanner struct
func SpannerOFF() HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.sb = SubTracerContainer{}
		hs.st = StarterContainer{}
		hs.ed = EnderContainer{}
		hs.s = true
	}
}
