package money

import "net/http"

//SpanForwardingOptions allows gathering data from an HTTP response
//An use case for this is extracting WRP spans into golang money spans
type SpanForwardingOptions func(*http.Response) []string
