package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	timeFormat string = "2006-01-02 15:04:05"
)

// Quote modles a row of data in the us_stocks.minute table.
type Quote struct {
	TickerID  int       `json:"ticker_id"`
	Timestamp time.Time `json:"timestamp"`
	Bid       float64   `json:"bid"`
	Ask       float64   `json:"ask"`
	BidSize   int       `json:"bid_size"`
	AskSize   int       `json:"ask_size"`
}

// EmptyQuote initializes a Quote struct to have the right
// TickerID for the given symbol.
// If the symbol does not exist in the database,
// insert it first.
func EmptyQuote(symbol string, db *sql.DB) (Quote, error) {
	res := Quote{}
	id, err := queryTickerID(symbol, db)
	if err != nil {
		return res, err
	}

	// ticker_id is the primary key with auto increment,
	// so the lowest possible value for it is 1
	// if queryTickerID returns 0,
	// it means that the given symbol does not exist
	// in the current database
	if id == 0 {
		// Insert the given symbol into the database
		err := insertTicker(symbol, db)
		if err != nil {
			return res, err
		}
		id, err = queryTickerID(symbol, db)
		if err != nil {
			return res, err
		}

		// If id still comes back as 0,
		// it means there an error in the insertion process
		if id == 0 {
			return res, errors.New("failed to insert new ticker into databse")
		}
	}

	res.TickerID = id
	return res, nil
}

// Insert will take a Quote struct and insert it into the given database connection.
func (q Quote) Insert(db *sql.DB) error {
	// Format query
	header := "INSERT INTO us_stocks.minute (ticker_id, timestamp, bid, ask, bid_size, ask_size) "
	timestamp := q.Timestamp.Format(timeFormat)
	values := fmt.Sprintf("VALUES (%d, '%s', %.2f, %.2f, %d, %d);", q.TickerID, timestamp, q.Bid, q.Ask, q.BidSize, q.AskSize)

	insert, err := db.Query(header + values)
	if err != nil {
		return err
	}

	defer insert.Close()
	return nil
}
