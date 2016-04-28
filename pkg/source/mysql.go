/*
Package source implements a way to get data from external MySQL sources.
*/
package source

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/pkg/strategy"
	"github.com/gabelula/gosqljson"
	// importing mysql
	_ "github.com/go-sql-driver/mysql"
)

/* Implementing the Sources */

// MySQL is the struct that has the connection string to the external mysql database
type MySQL struct {
	Connection string
	Database   *sql.DB
}

/* Exported Functions */

// GetData returns the raw data from that entity
func (m MySQL) GetData(entityname string, options *Options) ([]map[string]interface{}, error) { //offset int, limit int, orderby string, q string

	// Get the corresponding table to the modelName
	tableName := strategy.GetEntityForeignName(entityname)
	tableFields := strategy.GetEntityForeignFields(entityname) // []map[string]string

	// open a connection
	db, err := m.open()
	if err != nil {
		log.Error(uuid, "mysql.getdata", err, "Connecting to mysql database.")
		return nil, err
	}
	defer m.close(db)

	// Fields for that external source table
	f := make([]string, 0, len(tableFields))
	for _, field := range tableFields {
		if field != nil {
			f = append(f, field["foreign"].(string))
		}
	}

	fields := strings.Join(f, ", ")
	if options.Orderby == "" {
		options.Orderby = strategy.GetOrderBy(entityname)
	}

	// Get only the fields that we are going to use
	// the query string . To Do. Select only the stuff you are going to use
	//query := strings.Join([]string{"SELECT", fields, "from", tableName, "order by", orderby, "limit", fmt.Sprintf("%v", offset), ", ", fmt.Sprintf("%v", limit)}, " ")
	var where string
	if options.Query != "" {
		where = fmt.Sprintf("where %s ", options.Query)
	}
	query := fmt.Sprintf("SELECT %s from %s %s order by %s limit %v, %v", fields, tableName, where, options.Orderby, options.Offset, options.Limit)

	data, err := gosqljson.QueryDbToMapJSON(db, "lower", query)
	if err != nil {
		log.Error(uuid, "mysql.getdata", err, "Running SQL query.")
		return nil, err
	}

	byt := []byte(data)

	var dat []map[string]interface{}
	err = json.Unmarshal(byt, &dat)
	if err != nil {
		log.Error(uuid, "mysql.getdata", err, "Unmarshalling the result of the query.")
		return nil, err
	}

	return dat, nil
}

// GetQueryData returns the raw data from the table based on the ids
func (m MySQL) GetQueryData(entityname string, options *Options, ids []string) ([]map[string]interface{}, error) { //offset int, limit int, orderby string

	// Get the corresponding entity to the entityname
	tableName := strategy.GetEntityForeignName(entityname)
	tableFields := strategy.GetEntityForeignFields(entityname) // []map[string]string

	// open a connection
	db, err := m.open()
	if err != nil {
		log.Error(uuid, "mysql.getquerydata", err, "Error connecting to mysql database.")
		return nil, err
	}
	defer m.close(db)

	// Fields for that external source table
	f := make([]string, 0, len(tableFields))
	for _, field := range tableFields {
		if field != nil {
			f = append(f, field["foreign"].(string))
		}
	}

	// all the fields
	fields := strings.Join(f, ", ")

	// if we are ordering by
	if len(options.Orderby) == 0 {
		options.Orderby = strategy.GetOrderBy(entityname)
	}

	var queryWhere string
	// if we are quering specifics recrords
	if len(ids) > 0 {
		idField := strategy.GetIDField(entityname)
		queryWhere = fmt.Sprintf("where %s in (%s)", idField, strings.Join(ids, ", "))
	}

	// Get only the fields that we are going to use
	// the query string . To Do. Select only the stuff you are going to use
	query := strings.Join([]string{"SELECT", fields, "from", tableName, queryWhere, "order by", options.Orderby, "limit", fmt.Sprintf("%v", options.Offset), ", ", fmt.Sprintf("%v", options.Limit)}, " ")

	data, err := gosqljson.QueryDbToMapJSON(db, "lower", query)
	if err != nil {
		log.Error(uuid, "mysql.getquerydata", err, "Running SQL query.")
		return nil, err
	}

	byt := []byte(data)

	var dat []map[string]interface{}
	err = json.Unmarshal(byt, &dat)
	if err != nil {
		log.Error(uuid, "mysql.getquerydata", err, "Unmarshalling the query.")
		return nil, err
	}

	return dat, nil
}

// IsWebService returns true only if the implementation of Sourcer is an API
func (m MySQL) IsWebService() bool {
	return false
}

//////* Not exported functions *//////

// ConnectionMySQL returns the connection string
func connectionMySQL() string {
	credentialD, ok := credential.(str.CredentialDatabase)
	if !ok {
		err := fmt.Errorf("Error asserting type CredentialDatabase from interface Credential.")
		log.Error(uuid, "mysql.connectionMySQL", err, "Asserting Type CredentialDatabase")
		return ""
	}
	return fmt.Sprintf("%s:%s@(%s:%s)/%s", credentialD.Username, credentialD.Password, credentialD.Host, credentialD.Port, credentialD.Database)
}

// Open gives back  DB
func (m *MySQL) open() (*sql.DB, error) {

	database, err := sql.Open("mysql", m.Connection)
	if err != nil {
		return nil, err
	}

	if err = database.Ping(); err != nil {
		return nil, err
	}

	m.Database = database

	return database, nil
}

// Close closes the db
func (m MySQL) close(db *sql.DB) error {
	return db.Close()
}
