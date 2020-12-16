package main

import (
	"fmt"
	"log"
	"os"
)

var (
	errorLogger *log.Logger
	eventLogger *log.Logger
)

type marketClosedError struct {
	symbol string
}

func (err marketClosedError) Error() string {
	return fmt.Sprintf("market closed when querying for %s quote", err.symbol)
}

func init() {
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
	if len(os.Args) < 3 {
		log.Fatal(fmt.Errorf("invalid arguments, please enter log path followed by ticker symbols to crawl"))
	}

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
	eventLogger.Println("Starting data collection")

	crawlSymbols(symbolsMap, dsn)
}
