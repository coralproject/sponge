/*
Package localDB implements a way to get data into our local database.

Possible local databases:
* MongoDB

*/
package localDB

import (
	"database/sql"

	"github.com/coralproject/mod-data-import/config"
)

/* Implementing the Local DB Connections */

//////////// MongoDB local ////////////

// MongoDB has the connection and db to local database
type MongoDB struct {
	connection string
	database   *sql.DB
}

////// Not exported functions //////

func mongoDBConnection(c config.Credentials) string {
	return c.Username + ":" + c.Password + "@" + c.Host + "/" + c.Database
}

////// Exported Functions //////

// NewLocalDB initialize the connection to the local database
func NewLocalDB() (*MongoDB, error) {
	return MongoDB{}, nil
}

// Push imports data into the database
func Push(Data) error {
	return nil
}
