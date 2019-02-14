package money

import (
	"net/http"
)

type client interface {
	Do(*http.Request) (*http.Response, error)
}

type Options struct{
	Finish bool
}

type moneyClient struct {
	client
	finish bool
}

func NewMoneyTransactor(o *Options) moneyClient {
	return moneyClient{finish: o.Options}
}

func (mc moneyClient) Do(request *http.Request) (*http.Response, error) {
	tracker, err := ExtractTrackerFromRequest(resp)
	if err == nil {
		request = InjectTrackerIntoRequest(SetRequestMoneyHeader(tracker, request))
	}

	resp, err := mc.Client.Do(request)
	if err != nil {
		return nil, err
	}

	tracker, err := ExtractTrackerFromResponse(resp)
	if err != nil {
		return nil, err
	}

	tracker.Finish()
	SetResponseMoneyHeader(tracker, response.Header())

	if mc.Finisher == true {
		WriteMapsToStringResult()
	}

	return response, nil
}

func (mc moneyClient) WriteMoneyResponse() {
	mc.finisher := true
}

func (mc MoneyClient) Monetize(next client) client {
	return moneyClient{client: next}
}








