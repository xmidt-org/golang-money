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

// trackExtractor extracts a tracker from a request's context. Is used
// in the subtracer option when there already exist a tracker to subtrace from.
func subTracer(r *http.Request) (*HTTPTracker, error) {
	return ExtractTracker(r)
}

// spanDecoder decodes money headers off a request. Is used
// in the starter option when an http tracker has yet to be created.
func spanDecoder(r *http.Request) (*Span, error) {
	tc, err := decodeTraceContext(r.Header.Get(MoneyHeader))
	if err != nil {
		fmt.Print(err)
	}

	s := NewSpan("HTTPSpan", tc)

	return s, err
}

// SubTracerON is an option to use the decorator as a subtracer.
func SubTracerON(ch chan<- *HTTPTracker) HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.sb.function, hs.sb.htChannel = subTracer, ch
		hs.st = StarterContainer{}
		hs.ed = EnderContainer{}
		hs.s = false
	}
}

// StarterON is an option to use the decorator as a starter.
func StarterON(ch chan<- *HTTPTracker) HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.st.function, hs.st.htChannel = spanDecoder, ch
		hs.sb = SubTracerContainer{}
		hs.ed = EnderContainer{}
		hs.s = false
	}
}

// End is an option to use the decorator as a Ender
func EnderON() HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.ed.function = subTracer
		hs.sb = SubTracerContainer{}
		hs.st = StarterContainer{}
		hs.s = false
	}
}

// SpannerOff turns off all of HTTPSpanner's possible states.
func SpannerOff() HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.sb = SubTracerContainer{}
		hs.st = StarterContainer{}
		hs.ed = EnderContainer{}
		hs.s = true
	}
}
