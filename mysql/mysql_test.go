// These tests are interactive in nature,
// make sure the values mathces your local database before proceeding.
package mysql_test

import (
	"database/sql"
	"log"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/paulmure/gofin/mysql"
)

func TestEmptyTicker(t *testing.T) {
	db, err := sql.Open("mysql", os.Getenv("US_STOCKS_DSN"))
	if err != nil {
		log.Fatal(err)
	}

	q, err := mysql.EmptyQuote("ISNS", db)
	if err != nil {
		log.Fatal(err)
	}
	if q.TickerID != 1 {
		t.Errorf("ticker_id for ISNS should be 1, got %d", q.TickerID)
	}

	q, err = mysql.EmptyQuote("TEST", db)
	if err != nil {
		log.Fatal(err)
	}
	if q.TickerID != 4 {
		t.Errorf("ticker_id for TEST should be 3, got %d", q.TickerID)
	}
}

func TestInsertQuote(t *testing.T) {
	db, err := sql.Open("mysql", os.Getenv("US_STOCKS_DSN"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	q, err := mysql.EmptyQuote("TEST", db)
	if err != nil {
		log.Fatal(err)
	}

	q.Bid = 333.12
	q.Ask = 661.33
	q.BidSize = 100
	q.AskSize = 400
	q.Timestamp = time.Now()
	err = q.Insert(db)
	if err != nil {
		log.Fatal(err)
	}
}
