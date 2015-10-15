/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API?

*/
package source

import "database/sql"

////// Structures  //////

// Data is a struct that has all the db rows and error field
type Data struct {
	Rows  *sql.Rows // Move into appropiate structure
	Error error
}

// Source is where the data is coming from (mysql, api)
type Source interface {
	NewSource() (*Source, error)
	Open() (*sql.DB, error)
	Close(db *sql.DB) error
	GetNewData() Data
}

// Possible Errors
// - TimeOutReader
// - WrongConnection
// - DataErrorReader
