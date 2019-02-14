package money

import (
	"net/http"
)

type Client interface {
	Transact(*http.Request) (*http.Response, error)
}

func NewTransactors(options TransactorOptions) *Transactors {
	t := new(Transactors)

	options(t)

	return t
}

// Transactors embeds Transactor Interface to allow any type.
type Transactors struct {
	Tr1d1umDec DecoratedClient
	ScytaleDec DecoratedClient
}

type Transactor func(*http.Request) (*http.Response, error)

// DecorateTransactor decorates basic transactors.
//
// In the xmidt ecosystem scytale is the only node that uses a normal transactor.
func DecorateTransactor(t Transactor) Transactor {
	return func(r *http.Request) (*http.Response, error) {
		tracker, err := ExtractTrackerFromRequest(r)
		if err == nil {
			if resp, err := t(r); err == nil {
				return resp, nil
			}
		}

		r = SetRequestMoneyHeader(tracker, r)
		if resp, err := t(r); err == nil {
			tracker, err := ExtractTrackerFromResponse(resp)
			err = tracker.Finish()
			if err != nil {
				return nil, err
			}

			return resp, nil
		}

		return nil, err
	}
}

/*
func (t *Transactors) DecorateTransactor(c DecoratedClient) DecoratedClient {
	switch {
	case t.Tr1d1umDec != nil:
		return t.Tr1d1umDec(c)
	case t.ScytaleDec != nil:
		return t.ScytaleDec(c)
	}

	return nil
}
*/
