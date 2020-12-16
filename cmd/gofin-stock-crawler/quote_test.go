package main

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestInitQuotes(t *testing.T) {
	symbols := []string{"ISNS", "NGD", "TKC", "TEST"}
	dsn := os.Getenv("US_STOCKS_DSN")

	symbolsMap, err := initEmptyQuotes(symbols, dsn)
	if err != nil {
		log.Fatal(err)
	}

	if symbolsMap["ISNS"].TickerID != 1 {
		t.Errorf("expected ISNS id to be 1, got %d", symbolsMap["ISNS"].TickerID)
	}

	if symbolsMap["NGD"].TickerID != 2 {
		t.Errorf("expected NGD id to be 2, got %d", symbolsMap["NGD"].TickerID)
	}

	if symbolsMap["TKC"].TickerID != 3 {
		t.Errorf("expected TKC id to be 3, got %d", symbolsMap["TKC"].TickerID)
	}

	fmt.Printf("TEST id = %d\n", symbolsMap["TEST"].TickerID)
}

// Note: This test is interactive, and checking for market state
// is temporarily commented out when testing this function.
func TestRealTimeQuotes(t *testing.T) {
	symbols := []string{"TEST"}
	dsn := os.Getenv("US_STOCKS_DSN")

	symbolsMap, err := initEmptyQuotes(symbols, dsn)
	if err != nil {
		log.Fatal(err)
	}

	market, err := realTimeQuote("MSFT", symbolsMap["TEST"])
	if err != nil {
		log.Fatal(err)
	}

	if market.Bid != 214.04 {
		t.Errorf("expected bid to be 214.04, got %.2f", market.Bid)
	}

	if market.Ask != 214.15 {
		t.Errorf("expected bid to be 214.15, got %.2f", market.Ask)
	}

	if market.BidSize != 12 {
		t.Errorf("expected bid size to be 12, got %d", market.BidSize)
	}

	if market.AskSize != 12 {
		t.Errorf("expected ask size to be 12, got %d", market.AskSize)
	}
}

func TestInsertQuotes(t *testing.T) {
	symbols := []string{"TEST"}
	dsn := os.Getenv("US_STOCKS_DSN")

	symbolsMap, err := initEmptyQuotes(symbols, dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = insertQuote("MSFT", dsn, symbolsMap["TEST"])
	if err != nil {
		log.Fatal(err)
	}
}
