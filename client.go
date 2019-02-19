package money

import (
	"log"
	"net/http"
)

type client interface {
	Do(*http.Request) (*http.Response, error)
}

type Options struct {
	Finish bool
}

type moneyClient struct {
	finish bool
	client
	moneyLog       log.Logger
	responseWriter http.ResponseWriter
}

func NewMoneyTransactor(o *Options) moneyClient {
	return moneyClient{finish: o.Finish}
}

func (mc moneyClient) Do(request *http.Request) (*http.Response, error) {
	tracker, err := ExtractTrackerFromRequest(request)
	if err != nil {
		// run client as normally
		resp, err := mc.client.Do(request)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	// run client with a span
	//
	// set the request's MoneyHeader with a new trace
	request = SetRequestMoneyHeader(tracker, request)
	resp, err := mc.client.Do(request)
	if err != nil {
		return nil, err
	}

	tracker, err = ExtractTrackerFromResponse(resp)
	if err != nil {
		return nil, err
	}

	// if span does not participate in round trips
	if tracker.CheckOneWay() {
		maps, err := tracker.SpansMap()
		if err != nil {
			return nil, err
		}

		// mc.moneyLog.Log(logging.MessageKey(), mapsToStringResult(maps))
	}

	return resp, nil
}

func (mc moneyClient) Monetize(next client) client {
	return moneyClient{client: next}
}
