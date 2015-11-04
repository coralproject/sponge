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
	"log" // To Do. It needs to use our logging system.

	configuration "github.com/coralproject/sponge/config"
	"github.com/coralproject/sponge/utils"
	"gopkg.in/mgo.v2"
)

/* Implementing the Local DB Connections */

// MongoDB has the connection and db to local database
type MongoDB struct {
	connection string
	session    *mgo.Session
}

// global variables related to configuration
var config = *configuration.New()               // Reads the configuration file
var credential = config.GetCredentials("local") // Gets the credentials

/* Exported Functions */

// NewLocalDB gets the connection's string to the mongo database returned into a MongoDB struct type
// Method required by localDB.Interface
func NewLocalDB() *MongoDB {
	// Get mongodb connection string
	return &MongoDB{connection: mongoDBConnection()}
}

// Add imports data into the database
// Method required by localDB.Interface
// m has to be already initialized with a connection
func (m MongoDB) Add(d utils.Data, dry bool) error {

	// Connect to the local Database
	db := m.session.DB(credential.Database)

	errLogin := db.Login(credential.Username, credential.Password)

	if errLogin != nil {
		log.Fatal("Error when authenticating the database. ", errLogin)
		return errLogin
	}
	defer db.Logout()

	valComments := make([]interface{}, len(d.Comments))
	for i, v := range d.Comments {
		// To Do. Convert __id into commentid to create better ObjectIds in Mongodb. v[_id] = v[commentid]
		valComments[i] = v
	}

	// To Do . For each record returned
	//     Check to ensure the document isn't already added
	//     If not, add the document and kick off import actions
	//     If it's there, update the document
	//     Update the log collection

	var errInsert error
	// We are going to import d into one collection (Get the name of the collection from the strategy configuration file)
	// It has to be batch insert to be efficient
	if !dry {
		errInsert = db.C("comments").Insert(valComments...)
	} else {
		fmt.Println("Not inserting anything into local database... ")
		errInsert = nil
	}

	if errInsert != nil {
		log.Fatal("Error when inserting data into the new collection. ", errInsert)
		return errInsert
	}

	return nil
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
