package bittrex

import (
	"time"
)

const (
	// DefaultBatchSize ...
	DefaultBatchSize = 5
	// DefaultBatchInterval ...
	DefaultBatchInterval = 250 * time.Millisecond
)

// Config stores Bittrex configuration options
type Config struct {
	Host          string
	BatchSize     int           // specifies how many ticker requests we send at once before waiting for next batch
	BatchInterval time.Duration // to space out ticker requests a bit so we don't DDOS the exchange
}
