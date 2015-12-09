/*
Package source implements a way to get data from external MongoDB sources.
*/
package source

import (
	"log"

	"gopkg.in/mgo.v2"

	configuration "github.com/coralproject/sponge/config"
	_ "github.com/go-sql-driver/mysql" // Check if this can be imported not blank. To Do.
)

// Global configuration variables that holds the credentials for mysql
var credentialMongo = config.GetCredential("mongodb", "source").(configuration.CredentialDatabase)

/* Implementing the Sources */

// MongoDB is the struct that has the connection string to the external mysql database
type MongoDB struct {
	Connection string
	Database   *mgo.Session
}

/* Exported Functions */

// GetTables gets all the tables names from this data source
func (m MongoDB) GetTables() ([]string, error) {
	keys := []string{}
	for k := range config.Strategy.Tables {
		keys = append(keys, k)
	}
	return keys, nil
}

// GetData returns the raw data from the tableName
func (m MongoDB) GetData(modelName string) ([]map[string]interface{}, error) { //(*sql.Rows, error) {

	var dat []map[string]interface{}

	// Get the corresponding table to the modelName
	collectionName := config.GetTableName(modelName)
	// tableFields := config.GetTableFields(modelName) // map[string]string

	// open a connection
	session, err := m.initSession()
	if err != nil {
		return nil, err
	}
	defer m.closeSession(session)

	cred := mgo.Credential{
		Username: credentialMongo.Username,
		Password: credentialMongo.Password,
	}

	errLogin := session.Login(&cred)
	if errLogin != nil {
		log.Fatal("Error when authenticating the database. ", errLogin)
		return nil, errLogin
	}

	db := session.DB(credentialMongo.Database)
	col := db.C(collectionName)

	err = col.Find(dat).All(dat)
	if err != nil {
		return nil, err
	}

	return dat, nil
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
		return nil, &connectError{connection: m.Connection}
	}

	m.Database = database

	return database, nil
}

// Close closes the db
func (m MongoDB) closeSession(session *mgo.Session) error {
	session.Close()
	return nil
}
