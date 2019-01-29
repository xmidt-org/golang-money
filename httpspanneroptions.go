package money

import (
	"fmt"
	"net/http"
)

// HTTPSpannerOptions types for a HTTPSpanner used to declare a HTTPSpanners state.
type HTTPSpannerOptions func(*HTTPSpanner)

// A list of a HTTPSpanners different states.
type SubTracer func(*http.Request) (*HTTPTracker, error)
type Starter func(*http.Request) (*Span, error)
type Ender func(*http.Request) (*HTTPTracker, error)

// subTracer extracts a tracker from a request's context. Its used
// in the subtracer option when there already exists a tracker to subtrace from.
func subTracer(r *http.Request) (*HTTPTracker, error) {
	return ExtractTracker(r)
}

// starter decodes money headers off a request. Its used
// in the starter option when an http tracker has yet to be created.
func starter(r *http.Request) (*Span, error) {
	tc, err := DecodeTraceContext(r.Header.Get(MoneyHeader))
	if err != nil {
		fmt.Print(err)
	}

	s := NewSpan("HTTPSpan", tc)

	return s, err
}

// SubTracerON is an option to use the decorator as a subtracer.
func ScytaleON() HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.Tr1d1um = nil
		hs.Scytale = subTracer
		hs.Talaria = nil
	}
}

// StarterON is an option to use the decorator as a starter.
func Tr1d1umON() HTTPSpannerOptions {
	return func(hs *HTTPSpanner) {
		hs.Tr1d1um = starter
		hs.Scytale = nil
		hs.Talaria = nil
	}
}
