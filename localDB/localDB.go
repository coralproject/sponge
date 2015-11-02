/*
Package localDB implements a way to get data into our local database.

Possible local databases:
* MongoDB

*/
package localDB

import "github.com/coralproject/sponge/utils"

////// Structures  //////

// LocalDB is an interface to the different local databases available
type LocalDB interface {
	NewLocalDB() (*LocalDB, error)
	Add(utils.Data, bool) error // Data is where we are pushing the data to, the boolean is to check if run it dry or not
}
