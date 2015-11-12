/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API

*/
package source

import (
	"fmt"

	configuration "github.com/coralproject/sponge/config"
	"github.com/coralproject/sponge/utils"
)

// Source is where the data is coming from (mysql, api)
type Source interface {
	NewSource() *Source
	GetNewData() utils.Data // Data is a struct that has all the db rows and error field
}

// Global configuration variables that holds all the configuration from the config file
var config = *configuration.New() // Reads the configuration file

//* Errors used in this package *//

// When trying to connect to the database with the connection string
type connectError struct {
	connection string
}

func (e connectError) Error() string {
	return fmt.Sprintf("Error when connecting to database with %s.", e.connection)
}

// When trying to query the database with the query string
type queryError struct {
	query string
}

func (e queryError) Error() string {
	return fmt.Sprintf("Error when quering the database with %s.", e.query)
}

// When trying to create a new model... <-- To Do
type modelError struct {
	model string
}

func (e modelError) Error() string {
	return fmt.Sprintf("Error when trying to create a new model %s.", e.model)
}
