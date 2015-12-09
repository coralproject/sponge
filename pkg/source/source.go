/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API

*/
package source

import configuration "github.com/coralproject/sponge/config"

// Global configuration variables that holds all the configuration from the config file
var config = *configuration.New() // Reads the configuration file

// Sourcer is where the data is coming from (mysql, api)
type Sourcer interface {
	GetData(string) ([]map[string]interface{}, error)
	GetTables() ([]string, error)
}

// New returns a new Source struct with the connection string in it
func New(d string) (Sourcer, error) {

	switch d {
	case "mysql":
		// Get MySQL connection string
		return MySQL{Connection: connectionMySQL(), Database: nil}, nil
	case "mongodb":
		// Get MongoDB connection string
		return MongoDB{Connection: connectionMongoDB(), Database: nil}, nil
	}

	return nil, notFoundError{dbms: d}
}
