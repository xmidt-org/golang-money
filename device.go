package money

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Comcast/webpa-common/device"
	"github.com/Comcast/webpa-common/wrp"
)

type DecodeRequestFunc func(io.Reader, wrp.Format) (*device.Request, error)

// DecorateDecodeRequestWithMoney extracts a money span from requests httpheader and meshes it with a httpbody so spans can be passed through a device's websocket.
func DecorateDecodeRequestWithMoney(dr DecodeRequestFunc, req *http.Request, f wrp.Format) (*device.Request, error) {
	if ok := CheckHeaderForMoneyTrace(req.Header()); ok {
		ht, err := ExtractTracker(req)
		if err != nil {
			return nil, err
		}

		ht, err := ht.SubTrace(req.Context(), nil)
		if err != nil {
			return nil, err
		}

		b, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		spansMaps := []byte(ht.SpansMaps())

	} else {
		return DecodeRequestFunc(dcr, req, f), nil
	}
}

// HTTPTrackerFromDeviceResponse builds a httptracker from a device.Response
func HTTPTrackerFromDeviceResponse(device.Response) *HTTPTracker {

	return httpTracker
}
