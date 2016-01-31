/*
Package source implements a way to get data from external MongoDB sources.
*/
package source

import (
	"fmt"

	"github.com/coralproject/sponge/pkg/log"
	"gopkg.in/mgo.v2"
)

/* Implementing the Sources */

// MongoDB is the struct that has the connection string to the external mysql database
type MongoDB struct {
	Connection string
	Database   *mgo.Session
}

/* Exported Functions */

// GetData returns the raw data from the tableName
func (m MongoDB) GetData(coralTableName string, limit int, offset int, orderby string) ([]map[string]interface{}, error) {

	var dat []map[string]interface{}

	// Get the corresponding table to the modelName
	collectionName := strategy.GetTableForeignName(coralTableName)
	// tableFields := config.GetTableFields(modelName) // map[string]string

	// open a connection
	session, err := m.initSession()
	if err != nil {
		log.Error("Importing", "GetData", err, "Init mongo session")
		return nil, err
	}
	defer m.closeSession(session)

	cred := mgo.Credential{
		Username: credentialMongo.Username,
		Password: credentialMongo.Password,
	}

	err = session.Login(&cred)
	if err != nil {
		log.Error("Importing", "GetData", err, "Login mongo session")
		return nil, err
	}

	db := session.DB(credentialMongo.Database)
	col := db.C(collectionName)

	fmt.Println("TO DO ", col)
	// err = col.Find(dat).All(dat)
	// if err != nil {
	// 	log.Error("Importing", "GetData", err, "Get collection")
	// 	return nil, err
	// }

	return dat, nil
}

// GetQueryData needs to be implemented for mongodb to implement the sourcer interface
func (m MongoDB) GetQueryData(string, int, int, string, []string) ([]map[string]interface{}, error) {
	return nil, nil
}

//////* Not exported functions *//////

// ConnectionMongoDB returns the connection string
func connectionMongoDB() string {
	return credentialMongo.Username + ":" + credentialMongo.Password + "@" + "/" + credentialMongo.Database
}

// Open gives back a pointer to the DB
func (m *MongoDB) initSession() (*mgo.Session, error) {

	database, err := mgo.Dial(m.Connection)
	if err != nil {
		log.Error("Importing", "initSession", err, "Dial into session")
		return nil, err
	}

	m.Database = database

	return database, nil
}

// Close closes the db
func (m MongoDB) closeSession(session *mgo.Session) error {
	session.Close()
	return nil
}
