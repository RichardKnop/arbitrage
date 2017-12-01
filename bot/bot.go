package bot

import (
	"log"

	"github.com/RichardKnop/arbitrage/types"
)

// Bot ...
type Bot struct {
	Exchanges []types.Exchange
	Tickers   map[string]map[string]*types.Ticker
	quit      chan int
}

// New returns new Bot instance
func New(exchanges ...types.Exchange) *Bot {
	return &Bot{
		Exchanges: exchanges,
		Tickers:   make(map[string]map[string]*types.Ticker),
		quit:      make(chan int),
	}
}

// Run ...
func (b *Bot) Run() error {
	tickers := make(chan *types.Ticker)

	for _, e := range b.Exchanges {
		go e.Run(tickers)
	}

	for {
		select {
		case <-b.quit:
			return nil
		case ticker := <-tickers:
			log.Print(ticker)
		}
	}
}

// Quit ...
func (b *Bot) Quit() error {
	b.quit <- 1

	for _, e := range b.Exchanges {
		if err := e.Quit(); err != nil {
			return err
		}
	}

	return nil
}
