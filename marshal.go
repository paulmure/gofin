package gofin

import "database/sql"

// An Inserter handles inserting financial data into a local database.
type Inserter interface {
	// Insert handles inserting a single row of market data into
	// the local database.
	Insert(db *sql.DB) error
}
