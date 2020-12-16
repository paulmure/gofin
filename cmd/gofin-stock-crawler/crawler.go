package main

import (
	"errors"
	"sync"
	"time"

	"github.com/paulmure/gofin/mysql"
)

func crawlSymbol(symbol, dsn string, emptyQuote mysql.Quote, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		err := insertQuote(symbol, dsn, emptyQuote)
		if err != nil {
			if errors.As(err, &marketClosedError{}) {
				eventLogger.Println(err.Error())
			} else {
				errorLogger.Println(err.Error())
			}
			return
		}
		time.Sleep(time.Minute)
	}
}

func crawlSymbols(symbolsMap map[string]mysql.Quote, dsn string) {
	nSymbols := len(symbolsMap)
	var wg sync.WaitGroup
	wg.Add(nSymbols)

	for symbol, emptyQuote := range symbolsMap {
		go crawlSymbol(symbol, dsn, emptyQuote, &wg)
	}

	wg.Wait()
}
