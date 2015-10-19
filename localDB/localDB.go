/*
Package localDB implements a way to get data into our local database.

Possible local databases:
* MongoDB

*/
package localDB

import "github.com/coralproject/mod-data-import/models"

////// Structures  //////

// Data is where we are pushing the data to
type Data struct {
	Comments []models.Comment // Look for some kind of structure where to put this data into
	error    error
}

// LocalDB is an interface to the different local databases available
type LocalDB interface {
	NewLocalDB() (*LocalDB, error)
	Add(Data) error
}
