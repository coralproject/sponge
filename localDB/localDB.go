/*
Package localDB implements a way to get data into our local database.

Possible local databases:
* MongoDB

*/
package localDB

import "github.com/coralproject/mod-data-import/utils"

////// Structures  //////

// LocalDB is an interface to the different local databases available
type LocalDB interface {
	NewLocalDB() (*LocalDB, error)
	Add(utils.Data) error // Data is where we are pushing the data to
}
