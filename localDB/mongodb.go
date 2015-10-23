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

	"github.com/coralproject/mod-data-import/config"
	"github.com/coralproject/mod-data-import/utils"
	"gopkg.in/mgo.v2"
)

/* Implementing the Local DB Connections */

// MongoDB has the connection and db to local database
type MongoDB struct {
	connection string
	session    *mgo.Session
}

/* Exported Functions */

// NewLocalDB gets the connection's string to the mongo database returned into a MongoDB struct type
// Method required by localDB.Interface
func NewLocalDB() (*MongoDB, error) {

	c, err := config.GetCredentials()
	if err != nil {
		log.Fatal("Error when trying to get credentials. ", err)
	}

	connection, err := mongoDBConnection(c)
	if err != nil {
		log.Fatal("Error when trying to get credentials to connect to mongodb. ", err)
	}

	// Get mongodb connection string
	return &MongoDB{connection: connection}, err
}

// Add imports data into the database
// Method required by localDB.Interface
func (m MongoDB) Add(d utils.Data, dry bool) error {

	err := m.open()
	if err != nil {
		log.Fatal("Error when connecting to MongoDB.", err)
		return err
	}
	defer m.close()

	db := m.session.DB("coral")

	errLogin := db.Login("gaba", "gabita")
	if errLogin != nil {
		log.Fatal("Error when authenticating the database. ", errLogin)
		return errLogin
	}
	defer db.Logout()

	valComments := make([]interface{}, len(d.Comments))
	for i, v := range d.Comments {
		// To Do. Convert __id into commentid to create better ObjectIds in Mongodb. v[__id] = v[commentid]
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

/* Not exported functions */

func mongoDBConnection(credentials []config.Credential) (string, error) {

	// look at the credentials related to mongodb
	for i := 0; i < len(credentials); i++ {
		if credentials[i].Adapter == "mongodb" {
			c := credentials[i]
			//mongodb: //[username:password@]host1[:port1][,host2[:port2],...[,hostN[:portN]]][/[database][?options]]
			connection := c.Username + ":" + c.Password + "@" + c.Host + "/" + c.Database
			return connection, nil
		}
	}

	err := fmt.Errorf("Error when trying to get the connection string for mongodb.")

	return "", err
}

// Open a connection to the mongodb
func (m MongoDB) open() error {
	session, err := mgo.Dial(m.connection)
	if err != nil {
		log.Fatal("Error when trying to connect to the mongo database. ", err)
		return err
	}

	m.session = session

	return nil
}

// Close closes the db
func (m MongoDB) close() error {
	m.session.Close()
	return nil
}
