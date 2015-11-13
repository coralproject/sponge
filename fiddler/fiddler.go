/*
Package fiddler implements a way to get data into our local database source.

Possible local sources:
* MongoDB
* Service Layer Shelf

Code Principles

* Secure. For MongoDB look at security checklist: https://docs.mongodb.org/master/administration/security-checklist/
* Intentional

*/
package fiddler

import (
	configuration "github.com/coralproject/sponge/config"
	"github.com/coralproject/sponge/models"
)

////// Structures  //////

// LocalDB is an interface to a local database source
type LocalDB interface {
	NewLocalDB() (*LocalDB, error)
	Add(string, []models.Model) error // []models.Model] is the structure we are using to recieve the data to be added
}

// global variables related to configuration
var config = *configuration.New() // Reads the configuration file
