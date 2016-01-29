/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API

*/
package source

import (
	"errors"

	str "github.com/coralproject/sponge/pkg/strategy"
)

// global variables related to strategy
var strategy str.Strategy

// Init initialize the needed variables
func Init() {
	strategy = str.New()
	// logLevel := func() int {
	// 	ll, err := cfg.Int("LOGGING_LEVEL")
	// 	if err != nil {
	// 		return log.USER
	// 	}
	// 	return ll
	// }
	//
	// log.Init(os.Stderr, logLevel)
}

// Sourcer is where the data is coming from (mysql, api)
type Sourcer interface {
	GetData(string, int, int, string) ([]map[string]interface{}, error) //tableName, offset, limit, orderby
	GetQueryData(string, int, int, string, []string) ([]map[string]interface{}, error)
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

	return nil, errors.New("Configuration not found for source database.")
}
