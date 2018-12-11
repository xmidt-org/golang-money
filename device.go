package money

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Comcast/webpa-common/device"
	"github.com/Comcast/webpa-common/wrp"
)

// Device is an interface that provides the approproiate functions to work with money spans between talaria and the devices.
type Device interface {
	Decode(DecodeRequestFunc, *http.Request, wrp.Format) (*device.Request, error)
	Encode(EncodeResponseFunc, http.ResponseWriter, *device.Response, wrp.Format) error
}

// Bridge holds the data necessary for money http x wrp transfers
type Bridge struct {
	request        *device.Request
	response       *device.Response
	talariaTracker *HTTPTracker
	deviceTracker  *HTTPTracker
}

// NewBridge creates a new bridge.
func NewBridge(httpTracker *HTTPTracker) *Bridge {
	return &Bridge{
		talariaTracker: httpTracker,
	}
}

type DecodeRequestFunc func(io.Reader, wrp.Format) (*device.Request, error)

// Decode injects money into the requests response body.
func (b Bridge) Decode(d DecodeRequestFunc, req *http.Request, f wrp.Format) (*device.Request, error) {
	if ok := b.check(req); ok {
		var contents *[]byte
		contents, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		message := new(wrp.Message)
		message.Traces = b.talariaTracker.spansMaps

		decoder := wrp.NewDecoderBytes(contents, f).Decode(message)

		// add http tracker to the message and get a new encoder.
		encoder := wrp.NewEncoderBytes(contents, f)

		format, err := wrp.FormatFromContentType(req.Header.Get("Content-Type"), f)
		if err != nil {
			return nil, err
		}

		// obtain a new io reader with our encoder
		// TODO: TEST i need to test if this will wipe the customers body.
		msg := wrp.MustEncode(message, f)

		req.Body = bytes.NewReader(msg)

		req, err := d(req.Body, f)
		if err != nil {
			return nil, err
		}

		return req, nil
	} else {
		req, err := d(req.Body, f)
		if err != nil {
			return nil, err
		}

		return req, nil
	}
}

type EncodeResponseFunc func(http.ResponseWriter, *device.Response, wrp.Format) error

// Encode decorates a EncodeResponseFunc and writes money to http responses.
func (b Bridge) Encode(encoder EncodeResponseFunc, rw http.ResponseWriter, res *device.Response, f wrp.Format) (err error) {
	if ok := b.check(res); ok {
		b.build(res)
		b.write(rw)

		return encoder(rw, res, f)
	}

	return encoder(rw, res, f)
}

// Checks if the input contains a money span.
func (b Bridge) check(i interface{}) bool {
	switch v := i.(type) {
	case http.Request:
		m, _ := i.(http.Request)
		return CheckHeaderForMoneyTrace(m.Header)
	case device.Response:
		m, _ := i.(*device.Response)
		return CheckDeviceResponseForMoney(m)
	default:
		return false
	}
}

// Build builds a httptrackers from a device.Response so spans can be written to http responses.
func (b Bridge) build(res *device.Response) {
	trackers := res.Message.Traces

	b.buildSpan[len(trackers)-1]
	b.deviceTracker.spansMaps[len(b.deviceTracker.spansMaps)-1] = trackers[len(trackers)-1]

	b.buildSpan(trackers[len(trackers)-2])
	b.talariaTracker.spansMaps[len(b.talariaTracker.spansMaps)-2] = trackers[len(trackers)-2]
}

// Write writes talaria and caduceus spans to the http response
func (b Bridge) write(rw http.ResponseWriter) error {
	deviceSpan, err := b.deviceTracker.String()
	if err != nil {
		return err
	}

	h := rw.Header()
	h.Add(MoneySpansHeader, deviceSpan)

	talariaResult, _ := b.talariaTracker.Finish()
	if err != nil {
		return err
	}

	WriteMoneySpansHeader(talariaResult, rw)

	return nil
}

// BuildSpan builds a devies span from map object, specifically a device tracker.
func (b Bridge) buildSpan(t map[string]string) error {
	tc, err := DecodeTraceContext(t["TC"])
	if err != nil {
		return err
	}
	b.deviceTracker.span.TC = tc

	// TODO: all possible error codes or alternate route.
	if t["Code"] == "400" {
		b.deviceTracker.span.Code = 400
	} else {
		b.deviceTracker.span.Code = 401
	}

	if t["Success"] == "false" {
		b.deviceTracker.span.Success = false
	} else {
		b.deviceTracker.span.Success = true
	}

	b.deviceTracker.span.StartTime, err = time.Parse(t["StartTime"], "2011-01-19")
	if err != nil {
		return err
	}

	// TODO: ensure you can parse type durations like this
	b.deviceTracker.span.Duration, err = parseTime(t["Duration"])
	if err != nil {
		return err
	}

	b.deviceTracker.span.Name = t["Name"]
	b.deviceTracker.span.AppName = t["AppName"]
	b.deviceTracker.span.Err = errors.New(t["Err"])
	b.deviceTracker.span.Host = t["Host"]

	return nil
}

// CheckDeviceResponseForMoney checks if a device response contains money
func CheckDeviceResponseForMoney(res *device.Response) bool {
	if res.Message.Traces != nil {
		return true
	} else {
		return false
	}
}

// parseTime converts a string type to a duration type
func parseTime(t string) (time.Duration, error) {
	var mins, hours int
	var err error

	parts := strings.SplitN(t, ":", 2)

	switch len(parts) {
	case 1:
		mins, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}
	case 2:
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			return 0, err
		}

		mins, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, err
		}
	default:
		return 0, fmt.Errorf("invalid time: %s", t)
	}

	if mins > 59 || mins < 0 || hours > 23 || hours < 0 {
		return 0, fmt.Errorf("invalid time: %s", t)
	}

	return time.Duration(hours)*time.Hour + time.Duration(mins)*time.Minute, nil
}
