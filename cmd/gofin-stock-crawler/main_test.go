package main

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestInsertQueryFormatter(t *testing.T) {
	dbDSN := os.Getenv(dsn)
	if dbDSN == "" {
		t.Error("Invalid DSN")
	}

	db, err := sql.Open("mysql", dbDSN)

	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	q := IntradayQuoteEntry{TickerID: 1, BidSize: 5, AskSize: 5, Bid: 10.11, Ask: 11.12}

	insert, err := db.Query(q.formatInsertQuery())

	if err != nil {
		t.Errorf("Got an error when inserting query: %v", err)
	}
	defer insert.Close()
}

func TestTickerQuery(t *testing.T) {
	id := getTickerID("APPL")
	if id != 1 {
		t.Errorf("Expected APPL ticker id to be: %d, got: %d", 1, id)
	}

	id = getTickerID("TEST1")
}

func TestGetQuote(t *testing.T) {
	q, err := getRealTimeQuote("TSLA", 42)
	if err != nil {
		t.Error(err)
	}
	if q.TickerID != 42 {
		t.Errorf("Expected ticker_id to be: %d, got: %d", 42, q.TickerID)
	}
	if q.Ask != 0 {
		t.Errorf("Expected ask to be: %d, got: %.2f", 0, q.Ask)
	}
	if q.Bid != 0 {
		t.Errorf("Expected bid to be: %d, got: %.2f", 0, q.Bid)
	}
	if q.BidSize != 10 {
		t.Errorf("Expected bid size to be: %d, got: %d", 10, q.BidSize)
	}
	if q.AskSize != 11 {
		t.Errorf("Expected ask size to be: %d, got: %d", 11, q.AskSize)
	}
}
