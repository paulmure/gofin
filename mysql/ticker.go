package mysql

import (
	"database/sql"
	"fmt"
)

type ticker struct {
	ID     int    `json:"id"`
	Symbol string `json:"symbol"`
}

func queryTickerID(symbol string, db *sql.DB) (int, error) {
	query := fmt.Sprintf("SELECT id FROM us_stocks.ticker WHERE symbol='%s';", symbol)
	result, err := db.Query(query)
	if err != nil {
		return 0, err
	}
	defer result.Close()

	id := 0
	for result.Next() {
		var t ticker
		result.Scan(&t.ID)
		id = t.ID
	}

	return id, nil
}

func insertTicker(symbol string, db *sql.DB) error {
	query := fmt.Sprintf("INSERT INTO us_stocks.ticker (symbol) VALUES ('%s');", symbol)
	insert, err := db.Query(query)

	if err != nil {
		return err
	}
	defer insert.Close()

	return nil
}
