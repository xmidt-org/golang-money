package money

import (
	"context"
)

// Spanner is an interface which represents different ways spanning occurs.
type Spanner interface {
	Start(context.Context, Span) Tracker

	SubTrace(context.Context, Span) (*HTTPTracker, error)
}
