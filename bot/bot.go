package bot

import (
	"log"
	"sync"

	"github.com/RichardKnop/arbitrage/types"
)

// Bot ...
type Bot struct {
	Exchanges []types.Exchange
	Tickers   map[string]map[string]*types.Ticker
	quit      chan int
	wg        *sync.WaitGroup
}

// New returns new Bot instance
func New(exchanges ...types.Exchange) *Bot {
	return &Bot{
		Exchanges: exchanges,
		Tickers:   make(map[string]map[string]*types.Ticker),
		quit:      make(chan int),
		wg:        new(sync.WaitGroup),
	}
}

// Run ...
func (b *Bot) Run() error {
	tickers := make(chan *types.Ticker)
	errChan := make(chan error)

	for _, e := range b.Exchanges {
		go func() {
			b.wg.Add(1)

			if err := e.Run(tickers); err != nil {
				log.Print(err)
			}

			b.wg.Done()
		}()
	}

	go func() {
		for {
			select {
			case ticker := <-tickers:
				log.Print(ticker)
			case <-b.quit:
				errChan <- nil
			default:
			}
		}
	}()

	return <-errChan
}

// Quit ...
func (b *Bot) Quit() {
	// Trigger graceful shutdown of all exchange processes
	for _, e := range b.Exchanges {
		if err := e.Quit(); err != nil {
			log.Print(err)
		}
	}

	// Wait for quit process of exchanges to complete
	log.Print("Waiting for all exchanges to quit gracefully ")
	b.wg.Wait()

	b.quit <- 1
}
