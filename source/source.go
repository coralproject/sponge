/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API?

*/
package source

import "github.com/coralproject/mod-data-import/models"

////// Structures  //////

// Data is a struct that has all the db rows and error field
type Data struct {
	//Rows     *sql.Rows // Move into appropiate structure because this has to work for API too (no database/sql)
	Comments []models.Comment
	Error    error
}

// Source is where the data is coming from (mysql, api)
type Source interface {
	NewSource() (*Source, error)
	GetNewData() Data
}

// Possible Errors
// - TimeOutReader
// - WrongConnection
// - DataErrorReader
