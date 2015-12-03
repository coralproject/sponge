/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API

*/
package source

import configuration "github.com/coralproject/sponge/config"

// Source is where the data is coming from (mysql, api)
type Source interface {
	GetData(string) ([]map[string]interface{}, error)
	GetTables() []string
}

// Global configuration variables that holds all the configuration from the config file
var config = *configuration.New() // Reads the configuration file

// NewSource returns a new Source struct with the connection string in it
func NewSource(d string) (Source, error) {

	if d == "mysql" {
		// Get MySQL connection string
		return MySQL{Connection: ConnectionMySQL(), Database: nil}, nil
	}
	// if d == "mongodb" {
	// 	// Get MySQL connection string
	// 	return MongoDB{Connection: ConnectionMongoDB(), Database: nil}, nil
	// }

	return nil, notFoundError{dbms: d}
}
