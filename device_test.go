package money

import (
	"net/http/httptest"
	"testing"

	"github.com/Comcast/webpa-common/device"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

/*
func TestDecode(t *testing.T) {

}

*/

func TestEncode(t *testing.T) {
	var (
		tracker = 
		bridge
		input
	)
	
	// 1. test regulary 
	bridge.Encode


	// 2. test if body sustains content 
	


	// d


}


//
func TestCheckRequest(t *testing.T) {
	var (
		tracker = &HTTPTracker{}
		bridge  = NewBridge(tracker)
		input   = httptest.NewRequest("PUT", nil, nil)
	)

	if ok = bridge.check(input); ok {
		assert.True(t, ok)
	}
}

// checks
func TestCheckDevice(t *testing.T) {
	var (
		tracker = &HTTPTracker{}
		bridge  = NewBridge(tracker)
		input   = &device.Response{}
	)

	if ok = bridge.check(input); ok {
		assert.True(t, ok)
	}
}

func TestBuild(t *testing.T) {
	var (
		tracker = &HTTPTracker{}
		bridge  = NewBridge(tracker)
		res     = device.Response{}
	)

	bridge.Build(res)
	spew.Sdump(bridge.deviceTracker.spanMaps, bridge.talariaTracker.spanMaps)
	if len(bridge.deviceTracker) < len(bridge.talariaTracker) {
		t.Fatalf("talariaTracker's map list should be larger then the devices")
	}
}

func TestWrite(t *testing.T) {}

func TestbuildDeviceSpanFromTrackerMap(t *testing.T) {
	var (
		tracker = &HTTPTracker{}
		bridge  = NewBridge(tracker)
		res     = device.Response{}
	)

	buildDeviceSpanFromTrackerMap(

	builDeviceSpanFromTrackerMap()

}

func TestResult(t *testing.T) {}

func TestWriteDeviceSpanToHeaders(t *testing.T) {}

func TestCheckDeviceResponseForMoney(t *testing.T) {}
