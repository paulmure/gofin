package main

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/paulmure/gofin/mysql"
	"github.com/piquette/finance-go/quote"
)

func initEmptyQuotes(symbols []string, dsn string) (map[string]mysql.Quote, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	res := make(map[string]mysql.Quote, 0)
	for _, symbol := range symbols {
		emptyQuote, err := mysql.EmptyQuote(symbol, db)
		if err != nil {
			return nil, err
		}
		res[symbol] = emptyQuote
	}

	return res, nil
}

func realTimeQuote(symbol string, q mysql.Quote) (mysql.Quote, error) {
	// get a realtime quote from Yahoo Finance
	market, err := quote.Get(symbol)
	if err != nil {
		return q, err
	}

	// make sure the market is open at the time of the quote
	if market.MarketState != "REGULAR" {
		return q, marketClosedError{symbol}
	}

	// convert time
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return q, err
	}
	q.Timestamp = time.Now().In(loc)

	q.Bid = market.Bid
	q.Ask = market.Ask
	q.BidSize = market.BidSize
	q.AskSize = market.AskSize
	return q, nil
}

func insertQuote(symbol, dsn string, emptyQuote mysql.Quote) error {
	// get a real time market quote
	market, err := realTimeQuote(symbol, emptyQuote)
	if err != nil {
		return err
	}

	// set up a database connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// insert market quote into database
	err = market.Insert(db)
	if err != nil {
		return err
	}

	return nil
}
