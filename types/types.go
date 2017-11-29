package types

import (
	"time"

	"github.com/shopspring/decimal"
)

// Ticker ...
type Ticker struct {
	Bid  decimal.Decimal
	Ask  decimal.Decimal
	Last decimal.Decimal
	Time time.Time
}

// Exchange ...
type Exchange interface {
	Run(tickers chan *Ticker) error
	Quit() error
}
