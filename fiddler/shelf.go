/*
Package fiddler implements a way to get data into our local database source.

Eventually this will be in https://github.com/CoralProject/fiddler

*/
package fiddler

import (
	"log"
	"net/http"

	configuration "github.com/coralproject/sponge/config"
	"github.com/coralproject/sponge/models"
	"github.com/coralproject/sponge/utils"
)

/* Implementing the Local DB Connection to Shelf */

// Shelf has the configuration to Shelf Service Layer
type Shelf struct {
	credential configuration.CredentialAPI
}

/* Exported Functions */

// NewLocalDB gets the connection's string to the shelf instance
// Method required by localDB.Interface
func (s *Shelf) NewLocalDB() *Shelf { // To Do

	var credential configuration.CredentialAPI
	credential = config.GetCredential("shelf").(configuration.CredentialAPI)

	// Get mongodb connection string
	return &Shelf{credential: credential} // Gets the credentials}
}

// Authenticate into service layer
func (s *Shelf) Authenticate() (*utils.Resource, error) {
	//Authenticate against service layer
	//	API, err := BasicAuth{s.credential.Username, s.credential.Password}

	var r *utils.Resource

	basicauth := utils.BasicAuth{s.credential.Username, s.credential.Password}

	APIURL, err := s.credential.GetAuthenticationEndpoint()
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	// Get new Resource
	//r.Res()

	// Authenticate on APIURL
	//r.Post()

	return r, nil
}

// Add imports data into the collection collection in shelf
// m has to be already initialized with a connection
func (s Shelf) Add(modelName string, data []models.Model) error {

	r, err := s.Authenticate()
	if err != nil {
		log.Fatal(err)
		return err
	}

	APIURL, err := s.credential.GetEndpoint(modelName)
	if err != nil {
		log.Fatal(err)
	}

	// Build the request
	req, err := http.NewRequest("GET", APIURL, nil)
	if err != nil {
		return err
	}

	// Send the request via a client
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	// Call service layer endpoint to add data to model modelName
	return nil
}
