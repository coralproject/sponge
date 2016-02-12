/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API

*/
package source

import (
	"fmt"

	str "github.com/coralproject/sponge/pkg/strategy"
)

// global variables related to strategy
var (
	strategy str.Strategy
	uuid     string
)

// Global configuration variables that holds the credentials for the foreign database connection
var credential str.CredentialDatabase

// Init initialize the needed variables
func Init(u string) string {

	uuid = u
	str.Init(uuid)

	strategy = str.New()

	credential = strategy.GetCredential(strategy.Map.Foreign, "foreign")

	return strategy.Map.Foreign
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

	return nil, fmt.Errorf("Configuration not found for source database %s.", d)
}

// GetTables gets all the tables names from this data source
func GetTables() ([]string, error) {

	keys := make([]string, len(strategy.Map.Tables))

	for k, val := range strategy.Map.Tables {
		keys[val.Priority] = k
	}
	return keys, nil
}
