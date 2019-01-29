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

// Decorator is a function that decorates Clients
type Decorator func(Client) Client

// Transactors embeds Transactor Interface to allow any type.
type Transactors struct {
	tr1d1umDec Decorator
	scytaleDec Decorator
}

func (t *Transactors) DecorateTransactor(c Client) Client {
	switch {
	case t.tr1d1umDec != nil:
		return t.tr1d1umDec(c)
	case t.scytaleDec != nil:
		return t.scytaleDec(c)
	}

	return nil
}

// TransactorOptions used to declare DecorateTransactors state by adjusting a
// structs field
type TransactorOptions func(*Transactors)
