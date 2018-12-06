package money

import "net/http"

// An use case for this is extracting WRP spans into golang money spans
type SpanForwardingOptions func(*http.Response) []string

type TransactorTypes {
	WRPConverter SpanForwardingOptions
}

type WRPConverter func(*http.Request) (*HTTPTracker, error)

// WRP messages are messagepack? 
// Converts WRP messages into money compatible ones
// TODO: this is subjected to change X-Xmidt-Span may not be the appropraite span. 
func wrpConverter(resp http.Response) string {
	//deserialize response 


	return moneySpan
}

func WRPConverterON() SpanForwardingOptions {
	return func(tt TransactorTypes) {
		tt.WRPConverter = wrpConverter
	}
}

"parent": "parent string" – the parent id to tie the entry to
"name": "my name string" – the name of this entry (could later be used as a parent)
"start": 1234565, – the starting time of the window
"duration": 3434, – the duration of the window
"status": 200 – the http/ccsp response code for the trace entry





// 
and array of maps with key value pairs 
