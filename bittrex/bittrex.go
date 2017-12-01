// Package bittrex wraps the exchange API, see: https://bittrex.com/home/api
package bittrex

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/RichardKnop/arbitrage/types"
	"github.com/shopspring/decimal"
)

const (
	// Name is a unique exchange name
	Name = "bittrex"
)

var (
	// ErrEmptyResult is returned on edge case when response's success flag is true but result is null for some reason
	ErrEmptyResult = errors.New("Empty result")
)

// Exchange wraps methods that interact with exchange
type Exchange struct {
	cnf        *Config
	client     *http.Client
	quit       chan int
	wg         *sync.WaitGroup
	batch      []string
	batchCount int
}

// New returns new instance of Exchange
func New(cnf *Config) *Exchange {
	return &Exchange{
		cnf: cnf,
		client: &http.Client{
			Timeout: time.Duration(2 * time.Second), // set timeout to reasonably low period
		},
		quit:  make(chan int),
		wg:    new(sync.WaitGroup),
		batch: make([]string, cnf.BatchSize),
	}
}

// GetName returns a unique identifier for this exchange
func (e *Exchange) GetName() string {
	return Name
}

// Run ...
func (e *Exchange) Run(tickers chan *types.Ticker) error {
	errChan := make(chan error)

	go func() {
		errChan <- e.getTickersInBatches(tickers)
	}()

	return <-errChan
}

// Quit ...
func (e *Exchange) Quit() error {
	log.Printf("[%s] Quitting the running goroutine", e.GetName())
	e.quit <- 1

	log.Printf("[%s] Wait for ticker goroutines to finish", e.GetName())
	e.wg.Wait()

	return nil
}

func (e *Exchange) getTickersInBatches(tickers chan *types.Ticker) error {
	for {
		// Get all available markets
		markets, err := e.GetMarkets()
		if err != nil {
			return fmt.Errorf("[%s] Get markets error: %v", e.GetName(), err)
		}

		for i, m := range markets {
			// Capture quit channel here so we can exit the loop
			select {
			case <-e.quit:
				return nil
			default:
			}

			// Add the market to batch slice
			e.batch[e.batchCount] = m.MarketName
			e.batchCount++

			// If we have filled the batch slice or this is the last iteration in the loop
			if e.batchCount == e.cnf.BatchSize-1 || i == len(markets)-1 {
				// Execute batch of ticker requests
				for _, marketName := range e.batch {
					e.wg.Add(1)
					go e.getTicker(marketName, tickers)
				}

				// Reset the batch
				e.batchCount = 0
				e.batch = make([]string, e.cnf.BatchSize)

				// Space out batch requests
				<-time.After(e.cnf.BatchInterval)
			}
		}
	}

	return nil
}

func (e *Exchange) getTicker(marketName string, tickers chan *types.Ticker) error {
	defer e.wg.Done()

	// If the market name is empty string, ignore
	if marketName == "" {
		return nil
	}

	// Get the ticker for this market name
	ticker, err := e.GetTicker(marketName)
	if err != nil {
		return fmt.Errorf("[%s] Get ticker for '%s' error: %v", e.GetName(), marketName, err)
	}

	// Push the ticker to the upstream channel
	tickers <- &types.Ticker{
		Bid:  decimal.NewFromFloat(ticker.Bid),
		Ask:  decimal.NewFromFloat(ticker.Ask),
		Last: decimal.NewFromFloat(ticker.Last),
		Time: time.Now(),
	}

	return nil
}
