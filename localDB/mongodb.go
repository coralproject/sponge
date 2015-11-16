/*
Package localDB implements a way to get data into our local mongo database.

Possible local databases:
* MongoDB

Code Principles

* Secure. For MongoDB look at security checklist: https://docs.mongodb.org/master/administration/security-checklist/
* Intentional

*/
package localDB

import (
	"fmt"
	"log" // To Do. It needs to use our logging systema

	configuration "github.com/coralproject/sponge/config"
	"github.com/coralproject/sponge/models"
	"gopkg.in/mgo.v2"
)

/* Implementing the Local DB Connections */

// MongoDB has the connection and db to local database
type MongoDB struct {
	connection string
	session    *mgo.Session
}

// global variables related to configuration
var config = *configuration.New()                // Reads the configuration file
var credential = config.GetCredential("mongodb") // Gets the credentials

/* Exported Functions */

// NewLocalDB gets the connection's string to the mongo database returned into a MongoDB struct type
// Method required by localDB.Interface
func NewLocalDB() *MongoDB {
	// Get mongodb connection string
	return &MongoDB{connection: mongoDBConnection()}
}

// Add imports data into the collection collection in mongodb
// m has to be already initialized with a connection
func (m MongoDB) Add(collection string, data []models.Model, dry bool) error {

	err := m.Open()
	if err != nil {
		fmt.Println("Error when opening connection to Mongodb. ", err)
	}

	// Connect to the local Database and close connection when done
	db := m.session.DB(credential.Database)

	errLogin := db.Login(credential.Username, credential.Password)
	if errLogin != nil {
		log.Fatal("Error when authenticating the database. ", errLogin)
		return errLogin
	}
	defer db.Logout()

	// Push Data
	var errI error
	if !dry {

		// To Do: We need to find a better way to send the data... []interface{} and []model.Model are different type...
		new := make([]interface{}, len(data))
		for i, v := range data {
			new[i] = v
		}
		// INSERT Collection
		errI = db.C(collection).Insert(new...)
		if errI != nil {
			log.Fatal("Error when inserting data into the new collection. ", errI)
		}
		n, _ := db.C(collection).Count()
		fmt.Printf("The collection %s contains %d records.\n", collection, n)
	} else {
		log.Println("Running DRY: Not inserting anything into local database... ")
		errI = nil
	}

	return errI
}

// Open a connection to the mongodb
// We are setting m's session field so m has to be a pointer.
func (m *MongoDB) Open() error {
	var err error

	m.session, err = mgo.Dial(m.connection)

	if err != nil {
		log.Fatal("Error when trying to connect to the mongo database. ", err)
	}

	return err
}

// Close closes the db
func (m MongoDB) Close() error {

	if m.session != nil {
		m.session.Close()
	}

	return nil
}

/* Not exported functions */

// Returns the connection string
func mongoDBConnection() string {
	//mongodb: //[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
	connection := credential.Username + ":" + credential.Password + "@" + credential.Host + "/" + credential.Database
	return connection
}
