/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API

*/
package source

import (
	"database/sql"

	configuration "github.com/coralproject/sponge/config"
)

// Source is where the data is coming from (mysql, api)
type Source interface {
	NewSource() *Source
	GetData(string) (*sql.Rows, error)
}

// Global configuration variables that holds all the configuration from the config file
var config = *configuration.New() // Reads the configuration file
