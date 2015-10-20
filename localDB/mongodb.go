/*
Package localDB implements a way to get data into our local mongo database.

Possible local databases:
* MongoDB

Code Principles

Check list at ...

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

//////////// MongoDB local ////////////

// MongoDB has the connection and db to local database
type MongoDB struct {
	connection string
	session    *mgo.Session
}

////// Not exported functions //////

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

////// Exported Functions //////

// NewLocalDB initialize the connection to the local database
func NewLocalDB() (*MongoDB, error) {

	c, err := config.GetCredentials()

	if err != nil {
		log.Fatal("Error when trying to create new local db", err)
	}

	connection, err := mongoDBConnection(c)
	if err != nil {
		log.Fatal("Error when trying to get credentials to connect to mongodb. ", err)
	}

	// Get MySQL connection string
	return &MongoDB{connection: connection}, err
}

// Open a connection to the mongodb
func (m MongoDB) Open() (*mgo.Session, error) {
	session, err := mgo.Dial(m.connection)
	if err != nil {
		log.Fatal("Error when trying to connect to the mongo database. ", err)
		return nil, err
	}

	m.session = session

	return session, nil
}

// Close closes the db
func (m MongoDB) Close(session *mgo.Session) error {
	session.Close()
	return nil
}

// Add imports data into the database
func (m MongoDB) Add(d utils.Data) error {

	session, err := m.Open()
	if err != nil {
		log.Fatal("Error when connecting to MongoDB.", err)
		return err
	}
	defer session.Close()

	db := session.DB("coral")

	errLogin := db.Login("gaba", "gabita")
	if errLogin != nil {
		log.Fatal("Error when authenticating the database. ", errLogin)
		return errLogin
	}
	defer db.Logout()

	commentsCollection := db.C("Comments")

	errInsert := commentsCollection.Insert(d)
	if errInsert != nil {
		log.Fatal("Error when inserting data into the new collection. ", errInsert)
		return errInsert
	}

	return nil
}
