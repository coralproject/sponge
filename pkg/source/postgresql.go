/*
Package source implements a way to get data from external PostgreSQL sources.
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
	// import postgresql driver
	_ "github.com/lib/pq"
)

/* Implementing the Sources */

// PostgreSQL is the struct that has the connection string to the external postgresql database
type PostgreSQL struct {
	Connection string
	Database   *sql.DB
}

/* Exported Functions */

// GetData returns the raw data from the tableName
func (p PostgreSQL) GetData(entityname string, options *Options) ([]map[string]interface{}, error) { //offset int, limit int, orderby string, q string

	// Get the corresponding entity to the modelName
	tableName := strategy.GetEntityForeignName(entityname)
	tableFields := strategy.GetEntityForeignFields(entityname) // []map[string]string

	// open a connection
	db, err := p.open()
	if err != nil {
		log.Error(uuid, "postgresql.getdata", err, "Connecting to postgresql database.")
		return nil, err
	}
	defer p.close(db)

	// Fields for that external source table
	f := make([]string, 0, len(tableFields))
	for _, field := range tableFields {
		if field != nil {
			f = append(f, field["foreign"].(string))
		}
	}

	fields := strings.Join(f, ", ")
	if options.orderby == "" {
		options.orderby = strategy.GetOrderBy(entityname)
	}

	// Get only the fields that we are going to use
	var where string
	if options.query != "" {
		where = fmt.Sprintf("where %s ", options.query)
	}

	query := fmt.Sprintf("SELECT %s from %s %s order by %s OFFSET %v LIMIT %v", fields, tableName, where, options.orderby, options.offset, options.limit)

	data, err := gosqljson.QueryDbToMapJSON(db, "lower", query)
	if err != nil {
		log.Error(uuid, "postgresql.getdata", err, "Running SQL query.")
		return nil, err
	}

	byt := []byte(data)

	var dat []map[string]interface{}
	err = json.Unmarshal(byt, &dat)
	if err != nil {
		log.Error(uuid, "postgresql.getdata", err, "Unmarshalling the result of the query.")
		return nil, err
	}

	return dat, nil
}

// GetQueryData returns the raw data from the tableName
func (p PostgreSQL) GetQueryData(entityname string, options *Options, ids []string) ([]map[string]interface{}, error) { //offset int, limit int, orderby string

	// Get the corresponding table to the modelName
	tableName := strategy.GetEntityForeignName(entityname)
	tableFields := strategy.GetEntityForeignFields(entityname) // []map[string]string

	// open a connection
	db, err := p.open()
	if err != nil {
		log.Error(uuid, "postgresql.getquerydata", err, "Error connecting to postgresql database.")
		return nil, err
	}
	defer p.close(db)

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
	if len(options.orderby) == 0 {
		options.orderby = strategy.GetOrderBy(entityname)
	}

	var queryWhere string
	// if we are quering specifics recrords
	if len(ids) > 0 {
		idField := strategy.GetIDField(entityname)
		queryWhere = fmt.Sprintf("where %s in (%s)", idField, strings.Join(ids, ", "))
	}

	// Get only the fields that we are going to use
	// the query string . To Do. Select only the stuff you are going to use
	query := strings.Join([]string{"SELECT", fields, "from", tableName, queryWhere, "order by", options.orderby, "limit", fmt.Sprintf("%v", options.offset), ", ", fmt.Sprintf("%v", options.limit)}, " ")

	data, err := gosqljson.QueryDbToMapJSON(db, "lower", query)
	if err != nil {
		log.Error(uuid, "postgresql.getquerydata", err, "Running SQL query.")
		return nil, err
	}

	byt := []byte(data)

	var dat []map[string]interface{}
	err = json.Unmarshal(byt, &dat)
	if err != nil {
		log.Error(uuid, "postgresql.getquerydata", err, "Unmarshalling the query.")
		return nil, err
	}

	return dat, nil
}

// IsAPI is a func from the Sourcer interface to check if the external source is api or database
func (p PostgreSQL) IsAPI() bool {
	return false
}

//////* Not exported functions *//////

// ConnectionPostgresSQL returns the connection string
func connectionPostgreSQL() string {
	credentialD, ok := credential.(str.CredentialDatabase)
	if !ok {
		err := fmt.Errorf("Error asserting type CredentialDatabase from interface Credential.")
		log.Error(uuid, "postgresql.getquerydata", err, "Asserting Type CredentialDatabase")
		return ""
	}
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", credentialD.Username, credentialD.Password, credentialD.Host, credentialD.Database)

}

// Open gives back  DB
func (p *PostgreSQL) open() (*sql.DB, error) {

	database, err := sql.Open("postgres", p.Connection)
	if err != nil {
		return nil, err
	}

	if err = database.Ping(); err != nil {
		return nil, err
	}

	p.Database = database

	return database, nil
}

// Close closes the db
func (p PostgreSQL) close(db *sql.DB) error {
	return db.Close()
}
