/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API?

*/
package source

import "github.com/coralproject/mod-data-import/utils"

////// Structures  //////

// Data is a struct that has all the db rows and error field

// Source is where the data is coming from (mysql, api)
type Source interface {
	NewSource() (*Source, error)
	GetNewData() utils.Data // Data is a struct that has all the db rows and error field
}

// Possible Errors
// - TimeOutReader
// - WrongConnection
// - DataErrorReader
