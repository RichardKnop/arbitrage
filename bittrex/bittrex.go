// Package bittrex wraps the exchange API, see: https://bittrex.com/home/api
package bittrex

import (
	"errors"
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
	// BaseHost will be used by default
	BaseHost = "https://bittrex.com/api/v1.1"
	// GetMarketsEndpoint is public endpoint to get open markets
	GetMarketsEndpoint = "/public/getmarkets"
	// GetCurrenciesEndpoint is a public endpoint to get traded currencies
	GetCurrenciesEndpoint = "/public/getcurrencies"
	// GetTickerEndpoint is a public endpoint to get tickers
	GetTickerEndpoint = "/public/getticker"
)

var (
	// ErrEmptyResult is returned on edge case when response's success flag is true but result is null for some reason
	ErrEmptyResult = errors.New("Empty result")
	// GetTickerInterval spaces out ticker requests a bit so we don't DDOS the exchange
	GetTickerInterval = 10 * time.Millisecond
	// RefreshInterval specifies how often we want to get updated tickers from the exchange
	RefreshInterval = 10 * time.Second
)

// Exchange wraps methods that interact with exchange
type Exchange struct {
	host   string
	client *http.Client
	quit   chan int
	wg     *sync.WaitGroup
}

// New returns new instance of Exchange
func New(host string) *Exchange {
	if host == "" {
		host = BaseHost
	}
	return &Exchange{
		host:   host,
		client: new(http.Client),
	}
}

// GetName returns a unique identifier for this exchange
func (e *Exchange) GetName() string {
	return Name
}

// Run ...
func (e *Exchange) Run(tickers chan *types.Ticker) error {
	e.quit = make(chan int)
	e.wg = new(sync.WaitGroup)

	for {
		markets, err := e.GetMarkets()
		if err != nil {
			log.Printf("[%s] Get markets error: %v", e.GetName(), err)
			continue
		}

		select {
		case <-e.quit:
			return nil
		default:
			for _, m := range markets {
				e.wg.Add(1)
				go e.getTicker(m.MarketName, tickers)

				<-time.After(GetTickerInterval)
			}

			<-time.After(RefreshInterval)
		}
	}
}

// Quit ...
func (e *Exchange) Quit() error {
	log.Printf("[%s] Quitting the running goroutine", e.GetName())
	e.quit <- 0

	log.Printf("[%s] Wait for ticker goroutines to finish", e.GetName())
	e.wg.Wait()

	return nil
}

func (e *Exchange) getTicker(market string, tickers chan *types.Ticker) error {
	defer e.wg.Done()

	ticker, err := e.GetTicker(market)
	if err != nil {
		log.Printf("[%s] Get ticker error: %v", e.GetName(), err)
		return err
	}

	tickers <- &types.Ticker{
		Bid:  decimal.NewFromFloat(ticker.Bid),
		Ask:  decimal.NewFromFloat(ticker.Ask),
		Last: decimal.NewFromFloat(ticker.Last),
		Time: time.Now(),
	}

	return nil
}
