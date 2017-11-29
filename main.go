package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/RichardKnop/arbitrage/bittrex"
	"github.com/RichardKnop/arbitrage/bot"
)

func main() {
	quit := make(chan int)

	// Run the bot
	b := bot.New(bittrex.New(""))

	// Signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	// Goroutine Handle SIGINT and SIGTERM signals
	go func() {
		for {
			select {
			case s := <-sig:
				log.Printf("Signal received: %v", s)
				b.Quit()
				quit <- 0
				return
			}
		}
	}()

	b.Run()

	<-quit
}
