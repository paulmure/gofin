package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/piquette/finance-go/quote"
)

var (
	dsn         string
	logPath     string
	errorLogger *log.Logger
)

// IntradayQuoteEntry models a row of data in the intraday_by_minute table
type IntradayQuoteEntry struct {
	TickerID int     `json:"ticker_id"`
	BidSize  int     `json:"bid_size"`
	AskSize  int     `json:"ask_size"`
	Bid      float64 `json:"bid"`
	Ask      float64 `json:"ask"`
}

// Ticker is used to receive a ticker id from ticker table
type Ticker struct {
	ID int `json:"id"`
}

func (q IntradayQuoteEntry) formatInsertQuery() string {
	header := "INSERT INTO intraday_by_minute " +
		"(ticker_id, quote_timestamp, bid, ask, bid_size, ask_size) "
	values := fmt.Sprintf("VALUES (%d, NOW(), %.2f, %.2f, %d, %d);",
		q.TickerID, q.Bid, q.Ask, q.BidSize, q.AskSize)
	return header + values
}

func insertQuote(q IntradayQuoteEntry) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		errorLogger.Println("Failed to connect to SQL database")
	}
	defer db.Close()

	insert, err := db.Query(q.formatInsertQuery())

	if err != nil {
		errorLogger.Println("Failed to insert quote to SQL databse")
	}
	defer insert.Close()
}

func getTickerID(symbol string) (id int) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Try to query for a ticker id
	id, prs := queryTickerID(symbol, db)

	if !prs {
		// ticker symbol does not exist in database
		// insert the symbol and try again
		insertTicker(symbol, db)
		id, prs = queryTickerID(symbol, db)
		if !prs {
			// still can't insert symbol
			log.Fatal("Failed to insert new ticker into database")
		}
	}

	return
}

func queryTickerID(symbol string, db *sql.DB) (int, bool) {
	results, err := db.Query(fmt.Sprintf("SELECT id FROM ticker WHERE symbol = '%s';", symbol))

	if err != nil {
		log.Fatal(err)
	}
	defer results.Close()

	var t Ticker
	for results.Next() {
		err = results.Scan(&t.ID)

		if err != nil {
			log.Fatal(err)
		}
		return t.ID, true
	}

	return 0, false
}

func insertTicker(symbol string, db *sql.DB) {
	insert, err := db.Query(fmt.Sprintf("INSERT INTO ticker (symbol) VALUES ('%s');", symbol))
	if err != nil {
		log.Fatal(err)
	}
	insert.Close()
}

func getRealTimeQuote(symbol string, tickerID int) (quoteEntry IntradayQuoteEntry, retErr error) {
	q, err := quote.Get(symbol)
	if err != nil {
		errorLogger.Println(fmt.Sprintf("Failed to get quote for %s", symbol))
		retErr = errors.New("Failed to get quote")
		return
	}
	if q.MarketState != "REGULAR" {
		errorLogger.Fatal("Maket closed")
	}

	quoteEntry.TickerID = tickerID
	quoteEntry.BidSize = q.BidSize
	quoteEntry.AskSize = q.AskSize
	quoteEntry.Bid = q.Bid
	quoteEntry.Ask = q.Ask
	return
}

func init() {
	// Get DSN from environment for connecting to MySQL database
	dsn = os.Getenv("US_STOCKS_DSN")
	if dsn == "" {
		log.Fatal(fmt.Errorf("Cannot find DSN from environment"))
	}

	// Create log file
	file, err := os.OpenFile(os.Args[2], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	errorLogger = log.New(file, "ERROR: ", log.LstdFlags|log.Lshortfile)
}

func waitTill30() {
	for {
		if time.Now().Minute() >= 30 {
			break
		}
		time.Sleep(time.Second)
	}
}

func main() {
	symbol := os.Args[1]
	tickerID := getTickerID(symbol)

	log.Println("Waiting till market opens")
	waitTill30()
	log.Println("Starting data collection")

	for {
		q, err := getRealTimeQuote(symbol, tickerID)
		if err == nil {
			insertQuote(q)
		}
		time.Sleep(time.Minute)
	}
}
