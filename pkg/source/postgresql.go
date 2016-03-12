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

// GetTables gets all the tables names from this data source
func (m PostgreSQL) GetTables() ([]string, error) {
	keys := make([]string, len(strategy.Map.Tables))

	for k, val := range strategy.Map.Tables {
		keys[val.Priority] = k
	}
	return keys, nil
}

// GetData returns the raw data from the tableName
func (m PostgreSQL) GetData(coralTableName string, offset int, limit int, orderby string) ([]map[string]interface{}, error) {

	// Get the corresponding table to the modelName
	tableName := strategy.GetTableForeignName(coralTableName)
	tableFields := strategy.GetTableForeignFields(coralTableName) // []map[string]string

	// open a connection
	db, err := m.open()
	if err != nil {
		log.Error(uuid, "source.getdata", err, "Connecting to postgresql database.")
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
	if orderby == "" {
		orderby = strategy.GetOrderBy(coralTableName)
	}

	// Get only the fields that we are going to use
	// the query string . To Do. Select only the stuff you are going to use

	query := fmt.Sprintf("SELECT %s from %s order by %s OFFSET %v LIMIT %v", fields, tableName, orderby, offset, limit)

	data, err := gosqljson.QueryDbToMapJSON(db, "lower", query)
	if err != nil {
		log.Error(uuid, "source.getdata", err, "Running SQL query.")
		return nil, err
	}

	byt := []byte(data)

	var dat []map[string]interface{}
	err = json.Unmarshal(byt, &dat)
	if err != nil {
		log.Error(uuid, "source.getdata", err, "Unmarshalling the result of the query.")
		return nil, err
	}

	return dat, nil
}

// GetQueryData returns the raw data from the tableName
func (m PostgreSQL) GetQueryData(coralTableName string, offset int, limit int, orderby string, ids []string) ([]map[string]interface{}, error) {

	// Get the corresponding table to the modelName
	tableName := strategy.GetTableForeignName(coralTableName)
	tableFields := strategy.GetTableForeignFields(coralTableName) // []map[string]string

	// open a connection
	db, err := m.open()
	if err != nil {
		log.Error(uuid, "source.getquerydata", err, "Error connecting to postgresql database.")
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
	if len(orderby) == 0 {
		orderby = strategy.GetOrderBy(coralTableName)
	}

	var queryWhere string
	// if we are quering specifics recrords
	if len(ids) > 0 {
		idField := strategy.GetIDField(coralTableName)
		queryWhere = fmt.Sprintf("where %s in (%s)", idField, strings.Join(ids, ", "))
	}

	// Get only the fields that we are going to use
	// the query string . To Do. Select only the stuff you are going to use
	query := strings.Join([]string{"SELECT", fields, "from", tableName, queryWhere, "order by", orderby, "limit", fmt.Sprintf("%v", offset), ", ", fmt.Sprintf("%v", limit)}, " ")

	data, err := gosqljson.QueryDbToMapJSON(db, "lower", query)
	if err != nil {
		log.Error(uuid, "source.getdata", err, "Running SQL query.")
		return nil, err
	}

	byt := []byte(data)

	var dat []map[string]interface{}
	err = json.Unmarshal(byt, &dat)
	if err != nil {
		log.Error(uuid, "source.getdata", err, "Unmarshalling the query.")
		return nil, err
	}

	return dat, nil
}

//////* Not exported functions *//////

// ConnectionPostgresSQL returns the connection string
func connectionPostgreSQL() string {
	//return fmt.Sprintf("%s:%s@%s:%s/%s", credential.Username, credential.Password, credential.Host, credential.Port, credential.Database)
	//return fmt.Sprintf("%s?user=%s&password=%s&host=%s&port=%s&sslmode=disable", credential.Database, credential.Username, credential.Password, credential.Host, credential.Port)
	//return fmt.Sprintf("user=%s dbname=%s sslmode=verify-full", credential.Username, credential.Database)

	//db, err := sql.Open("postgres", "postgres://pqgotest:password@localhost/pqgotest?sslmode=verify-full")
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", credential.Username, credential.Password, credential.Host, credential.Database)

}

// Open gives back  DB
func (m *PostgreSQL) open() (*sql.DB, error) {

	database, err := sql.Open("postgres", m.Connection)
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
func (m PostgreSQL) close(db *sql.DB) error {
	return db.Close()
}
