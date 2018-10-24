package money

import "net/http"

//SpanForwardingOptions allows gathering data from an HTTP response
//into string-encoded golang money spans
//application code is responsible to only inspect the response and if otherwise, put back data
//(i.e if body is read)
//An use case for this is extracting WRP spans into golang money spans
type SpanForwardingOptions func(*http.Response) []string
