/*
Package fiddler implements a way to get data into our local database source.

Eventually this will be in https://github.com/CoralProject/fiddler

*/
package fiddler

import (
	"log"

	configuration "github.com/coralproject/sponge/config"
	"github.com/coralproject/sponge/models"
)

/* Implementing the Local DB Connection to Shelf */

// ShelfDB has the connection to Shelf Service Layer
type ShelfDB struct {
	credential configuration.Credential
}

/* Exported Functions */

// NewLocalDB gets the connection's string to the shelf instance
// Method required by localDB.Interface
func (s *ShelfDB) NewLocalDB() *ShelfDB { // To Do
	// Get mongodb connection string
	return &ShelfDB{credential: config.GetCredential("shelfdb")} // Gets the credentials}
}

// Authenticate into service layer
func (s *ShelfDB) Authenticate(apiKey string) error {
	//Authenticate against service layer
	// Error when failing
	return nil
}

// Add imports data into the collection collection in shelf
// m has to be already initialized with a connection
func (s ShelfDB) Add(modelName string, data []models.Model) error {

	err := s.Authenticate("s.credential.apiKey")
	if err != nil {
		log.Fatal(err)
		return err
	}

	// Call service layer endpoint to add data to model modelName
	return nil
}
