package money

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Comcast/webpa-common/device"
	"github.com/Comcast/webpa-common/wrp"
)

// Device is an interface that provides the approproiate functions to work with money spans between talaria and the devices.
type Device interface {
	Decode(DecodeRequestFunc, *http.Request, wrp.Format) (*device.Request, error)
	Encode(EncodeResponseFunc, http.ResponseWriter, *device.Response, wrp.Format) (error)
}

// Bridge holds the data necessary for money http x wrp transfers
type Bridge struct {
	request        *device.Request
	response       *device.Response
	talariaTracker *HTTPTracker
	deviceTracker  *HTTPTracker
}

// NewBridge creates a new bridge.
func NewBridge(httpTracker *HTTPTracker) {
	return Bridge{
		talariaTracker: httpTracker,
	}
}

type DecodeRequestFunc func(io.Reader, wrp.Format) (*device.Request, error)

// Decode injects money into the requests response body.
func (b Bridge) Decode(dr DecodeRequestFunc, req *http.Request, f wrp.Format) (*device.Request, error) {
	if ok = b.check(req); ok {
		contents, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		message := new(wrp.Message)
		message.Traces := b.talariaTracker.SpansMaps

		decoder := wrp.NewDecoderBytes(contents, f).Decode(message)

		// add http tracker to the message and get a new encoder.
		encoder := wrp.NewEncoderBytes(contents, wrp.FormatFromContentType(httpRequest.Header.Get("Content-Type"), wrp.Msgpack))

		// obtain a new io reader with our encoder 
		// TODO: TEST i need to test if this will wipe the customers body. 
		msg, err = wrp.MustEncode(message, wrp.FormatFromContentType(httRequest.Header.Get("Content-Type"), wrp.Msgpack))
		if err != nil {
			return nil, err
		}

		req.Body := bytes.NewReader(msg)

		return DecodeRequestFunc(req.Body, f), nil
	} else {
		return DecodeRequestFunc(req, f), nil
		}
}


type EncodeResponseFunc func(http.ResponseWriter, *device.Response, wrp.Format) (error)

// Encode decorates a EncodeResponseFunc and writes money to http responses.
func(b Bridge) Encode(EncodeResponseFunc, rw http.ResponseWriter, res *device.Response, f wrp.Format) (err error) {
	if ok = b.check(res); ok {
		b.build(res)
		b.write(rw)

		return EncodeResponseFunc(rw, res, f)
	}

	return EncodeResponseFunc(rw, res, f)
}

// Checks if the input contains a money span. 
func (b Bridge) check(i interface{}) bool {
	switch v := i.(Type) {
	case http.Request:
		m, _ := i.(http.Request)
		return CheckHeaderForMoneyTrace(m.Header())
	case Response:
		m, _ := i.(*device.Response)
		return CheckDeviceResponseForMoney(m)
	case default:
		return false
	}
}

// Build builds a httptrackers from a device.Response so spans can be written to http responses.
func (b Bridge) build(res device.Response) {
	trackers := res.Response.Message.Traces

	b.buildSpan[len(trackers)-1])
	b.deviceTracker.spanMaps[len(b.deviceTracker.spansMaps)-1] = trackers[len(trackers)-1]

	b.buildSpan(trackers[len(trackers)-2])
	b.talariaTracker.spanMaps[len(b.talariaTracker.spansMaps)-2] = trackers[len(trackers)-2]
}

// Write writes talaria and caduceus spans to the http response
func (b Bridge) write(rw http.ResponseWriter) {
	deviceResult := b.result()
	b.writeDeviceSpanToHeaders(b.deviceTracker, rw)

	talariaResult, _ := b.talariaTracker.Finish()
	WriteMoneySpansHeader(talariaResult)
}

// BuildSpan builds a span from map object, specifically a device tracker.
func (b Bridge) buildSpan(t map[string]string) {
	t.span.Name = t[Name]
	t.span.TC = t[TC]
	t.span.AppName = t[AppName]
	t.span.Code = t[Code]
	t.span.Success = t[Success]
	t.span.Err = t[Err]
	t.span.StartTime = t[StartTime]
	t.span.Duration = t[Duration]
	t.span.Host = t[Host]

	return
}

// Result turns spins up a result object of a device tracker
func (b Bridge) result() money.Result {
		return Result{
			Name:      b.deviceTracker.span.Name,
			TC:        b.deviceTracker.span.TC,
			AppName:   b.deviceTracker.span.AppName,
			Code:      b.deviceTracker.span.Code,
			Success:   t.deviceTraker.span.Success,
			Err:       t.deviceTracker.span.Err,
			StartTime: t.deviceTracker.span.StartTime,
			Duration:  t.deviceTracker.span.Duration,
			Host:      t.deviceTracker.span.Host,
		}
}

// write deviceSpanToHeaders writes a 
func (b Bridge) writeDeviceSpanToHeaders(ht *money.HTTPTracker) bool {
	tracker := b.deviceTracker
}

// CheckDeviceResponseForMoney checks if a device response contains money 
func CheckDeviceResponseForMoney(res *device.Response) bool {
		if res.Message.Traces != nil {
			return true
		} else {
			return false
		}
}
