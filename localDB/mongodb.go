/*
Package localDB implements a way to get data into our local mongo database.

Possible local databases:
* MongoDB

*/
package localDB

import (
	"fmt"
	"log"

	"github.com/coralproject/mod-data-import/config"
)

/* Implementing the Local DB Connections */

//////////// MongoDB local ////////////

// MongoDB has the connection and db to local database
type MongoDB struct {
	connection string
	//database   *bson
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

// Push imports data into the database
func Push(Data) error {
	return nil
}
