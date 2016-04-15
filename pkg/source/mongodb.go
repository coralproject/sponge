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
func (m MongoDB) GetData(entityname string, options *Options) ([]map[string]interface{}, error) { // offset int, limit int, orderby string, query string) ([]map[string]interface{}, bool, error) { //(*sql.Rows, error) {

	var data []map[string]interface{}

	// Get the corresponding entity to the coral collection name
	//	entity := strategy.GetEntityForeignName(entityname)
	fields := strategy.GetEntityForeignFields(entityname) //[]]map[string]string

	// open a connection
	session, err := m.initSession()
	if err != nil {
		log.Error(uuid, "mongodb.getdata", err, "Initializing mongo session.")
		return nil, err
	}
	defer m.closeSession(session)

	credentialD, ok := credential.(str.CredentialDatabase)
	if !ok {
		err = fmt.Errorf("Error asserting type CredentialDatabase from interface Credential.")
		log.Error(uuid, "mongodb.getdata", err, "Asserting Type CredentialDatabase")
		return nil, err
	}
	cred := mgo.Credential{
		Username: credentialD.Username,
		Password: credentialD.Password,
	}

	err = session.Login(&cred)
	if err != nil {
		log.Error(uuid, "mongodb.getdata", err, "Login mongo session.")
		return nil, err
	}

	db := session.DB(credentialD.Database)
	col := db.C(entityname)

	//Get all the fields that we are going to get from the document { field: 1}
	fieldsToGet := make(map[string]bool)
	//var fieldsNames []string
	for _, f := range fields {
		ff, ok := f["foreign"].(string)
		if !ok {
			err := fmt.Errorf("Error asserting type String from field.")
			log.Error(uuid, "mongodb.getdata", err, "Type asserting %v into string.", f["foreign"])
		}
		fieldsToGet[ff] = true
		//fieldsNames = append(fieldsNames, f["local"])
	}

	var mquery map[string]interface{}
	if options.Query != "" {
		err = json.Unmarshal([]byte(options.Query), &mquery)
		if err != nil {
			log.Error(uuid, "mongodb.getdata", err, "Unmarshalling query %v", options.Query)
			return nil, err
		}
	}

	//.Select(fieldsToGet) <--- I'm not using Select because SOME FIELDS IN THE TRANSLATION FILE ARE NOT THE RIGHT ONES TO DO THE SELECT. For example: context.object.0.uri
	err = col.Find(mquery).Limit(options.Limit).All(&data)
	if err != nil {
		log.Error(uuid, "mongodb.getdata", err, "Getting collection %s.", entityname)
		return nil, err
	}

	flattenData, err := normalizeData(data)
	if err != nil {
		log.Error(uuid, "mongodb.getdata", err, "Normalizing data from mongo to fit into fiddler.")
		return nil, err
	}

	return flattenData, nil
}

// GetQueryData needs to be implemented for mongodb to implement the sourcer interface
func (m MongoDB) GetQueryData(entity string, options *Options, ids []string) ([]map[string]interface{}, error) { //offset int, limit int, orderby string

	var d []map[string]interface{}
	var err error

	// if we are quering specifics recrords
	if len(ids) > 0 {
		idField := strategy.GetIDField(entity)

		for i, j := range ids {
			ids[i] = fmt.Sprintf("\"%s\"", j)
		}
		options.Query = fmt.Sprintf("{\"%s\": {\"$in\":[ %v ] } }", idField, strings.Join(ids, ", "))

		d, err = m.GetData(entity, options)
	} else {
		err = fmt.Errorf("No ids to get.")
	}

	return d, err
}

// IsWebService is used to check what is that sourcerer interface
func (m MongoDB) IsWebService() bool {
	return false
}

//////* Not exported functions *//////

// ConnectionMongoDB returns the connection string
func connectionMongoDB() string {
	credentialD, ok := credential.(str.CredentialDatabase)
	if !ok {
		log.Error(uuid, "mongodb.connectionMongoDB", fmt.Errorf("Error asserting type CredentialDatabase from interface Credential."), "Asserting Type CredentialDatabase")
		return ""
	}
	return fmt.Sprintf("%s:%s@/%s", credentialD.Username, credentialD.Password, credentialD.Database)
}

// Open gives back a pointer to the DB
func (m *MongoDB) initSession() (*mgo.Session, error) {

	database, err := mgo.Dial(m.Connection)
	if err != nil {
		log.Error(uuid, "mongodb.initsession", err, "Dial into session.")
		return nil, err
	}

	m.Database = database

	return database, nil
}

// Close closes the db
func (m MongoDB) closeSession(session *mgo.Session) {
	session.Close()
}
