/*
Package source implements a way to get data from external MongoDB sources.
*/
package source

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/pkg/strategy"
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
func (m MongoDB) GetData(coralTableName string, offset int, limit int, orderby string, query string) ([]map[string]interface{}, bool, error) { //(*sql.Rows, error) {

	var data []map[string]interface{}
	var notFinish = false

	// Get the corresponding entity to the coral collection name
	collectionName := strategy.GetEntityForeignName(coralTableName)
	fields := strategy.GetEntityForeignFields(coralTableName) //[]]map[string]string

	// open a connection
	session, err := m.initSession()
	if err != nil {
		log.Error(uuid, "source.getdata", err, "Initializing mongo session.")
		return nil, notFinish, err
	}
	defer m.closeSession(session)

	credentialD, ok := credential.(str.CredentialDatabase)
	if !ok {
		err = fmt.Errorf("Error asserting type CredentialDatabase from interface Credential.")
		log.Error(uuid, "source.getdata", err, "Asserting Type CredentialDatabase")
		return nil, notFinish, err
	}
	cred := mgo.Credential{
		Username: credentialD.Username,
		Password: credentialD.Password,
	}

	err = session.Login(&cred)
	if err != nil {
		log.Error(uuid, "source.getdata", err, "Login mongo session.")
		return nil, notFinish, err
	}

	db := session.DB(credentialD.Database)
	col := db.C(collectionName)

	//Get all the fields that we are going to get from the document { field: 1}
	fieldsToGet := make(map[string]bool)
	//var fieldsNames []string
	for _, f := range fields {
		fieldsToGet[f["foreign"].(string)] = true
		//fieldsNames = append(fieldsNames, f["local"])
	}

	var mquery map[string]interface{}
	if query != "" {
		err = json.Unmarshal([]byte(query), &mquery)
		if err != nil {
			log.Error(uuid, "mongo.getdata", err, "Unmarshalling query %v", query)
			return nil, notFinish, err
		}
	}

	//.Select(fieldsToGet) <--- SOME FIELDS ARE NOT THE RIGHT ONES TO DO THE SELECT. For example: context.object.0.uri
	err = col.Find(mquery).Limit(limit).All(&data)
	if err != nil {
		log.Error(uuid, "mongo.getdata", err, "Getting collection %s.", collectionName)
		return nil, notFinish, err
	}

	flattenData, err := normalizeData(data)
	if err != nil {
		log.Error(uuid, "source.getdata", err, "Normalizing data from mongo to fit into fiddler.")
		return nil, notFinish, err
	}

	return flattenData, notFinish, nil
}

// GetQueryData needs to be implemented for mongodb to implement the sourcer interface
func (m MongoDB) GetQueryData(coralTableName string, offset int, limit int, orderby string, ids []string) ([]map[string]interface{}, error) {

	var d []map[string]interface{}
	var err error

	// if we are quering specifics recrords
	if len(ids) > 0 {
		idField := strategy.GetIDField(coralTableName)

		for i, j := range ids {
			ids[i] = fmt.Sprintf("\"%s\"", j)
		}
		query := fmt.Sprintf("{\"%s\": {\"$in\":[ %v ] } }", idField, strings.Join(ids, ", "))

		d, _, err = m.GetData(coralTableName, offset, limit, orderby, query)
	} else {
		err = fmt.Errorf("No ids to get.")
	}

	return d, err
}

func (m MongoDB) IsAPI() bool {
	return false
}

//////* Not exported functions *//////

// ConnectionMongoDB returns the connection string
func connectionMongoDB() string {
	credentialD, ok := credential.(str.CredentialDatabase)
	if !ok {
		log.Error(uuid, "source.getdata", fmt.Errorf("Error asserting type CredentialDatabase from interface Credential."), "Asserting Type CredentialDatabase")
		return ""
	}
	return fmt.Sprintf("%s:%s@/%s", credentialD.Username, credentialD.Password, credentialD.Database)
}

// Open gives back a pointer to the DB
func (m *MongoDB) initSession() (*mgo.Session, error) {

	database, err := mgo.Dial(m.Connection)
	if err != nil {
		log.Error(uuid, "source.initsession", err, "Dial into session.")
		return nil, err
	}

	m.Database = database

	return database, nil
}

// Close closes the db
func (m MongoDB) closeSession(session *mgo.Session) {
	session.Close()
}
