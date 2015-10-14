/*
Package localDB implements a way to get data into our local database.

Possible local databases:
* MongoDB

*/
package localDB

import "database/sql"

////// Structures  //////

// Data is where we are pushing the data to
type Data struct {
	DB    sql.DB
	error error
}

// LocalDB is an interface to the different local databases available
type LocalDB interface {
	NewLocalDB() (*LocalDB, error)
	Push(Data) error
}

// Possible Errors
// - TimeOutReader
// - WrongConnection
// - DataErrorReader
