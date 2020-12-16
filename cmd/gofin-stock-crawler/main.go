package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	marketLoc   *time.Location
	errorLogger *log.Logger
	eventLogger *log.Logger
)

type marketClosedError struct {
	symbol string
}

func (err marketClosedError) Error() string {
	return fmt.Sprintf("market closed when querying for %s quote", err.symbol)
}

func waitTillMarketOpens() {
	for {
		t := time.Now().In(marketLoc)
		switch {
		case t.Hour() < 9:
			time.Sleep(time.Hour)
		case t.Minute() < 29:
			time.Sleep(time.Minute)
		case t.Minute() == 29:
			time.Sleep(time.Second)
		default:
			return
		}
	}
}

func init() {
	if len(os.Args) < 3 {
		log.Fatal(fmt.Errorf("invalid arguments, please enter log path followed by ticker symbols to crawl"))
	}

	err := errors.New("")
	marketLoc, err = time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}

	// Create log file
	file, err := os.OpenFile(os.Args[1], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	// register loggers
	errorLogger = log.New(file, "ERROR: ", log.LstdFlags|log.Lshortfile)
	eventLogger = log.New(file, "EVENT: ", log.LstdFlags|log.Lshortfile)
}

func main() {
	// Get DSN from environment for connecting to MySQL database
	dsn := os.Getenv("US_STOCKS_DSN")
	if dsn == "" {
		log.Fatal(fmt.Errorf("cannot find DSN from environment"))
	}

	symbolsMap, err := initEmptyQuotes(os.Args[2:], dsn)
	if err != nil {
		log.Fatal(err)
	}

	eventLogger.Println("Waiting till market opens")
	waitTillMarketOpens()
	eventLogger.Println("Starting data collection")

	crawlSymbols(symbolsMap, dsn)
}
