/*
Package source implements a way to get data from external MySQL sources.

External possible sources:
* MySQL
* API

*/
package source

import (
	"database/sql"
	"strings"

	configuration "github.com/coralproject/sponge/config"

	_ "github.com/go-sql-driver/mysql" // Check if this can be imported not blank. To Do.
)

// Global configuration variables that holds the credentials for mysql
var credential = config.GetCredential("mysql").(configuration.CredentialDatabase)

/* Implementing the Sources */

// MySQL is the struct that has the connection string to the external mysql database
type MySQL struct {
	Connection string
	Database   *sql.DB
}

/* Exported Functions */

// NewSource returns a new Mysql struct with the connection string in it
// Method required by source.Interface
func NewSource() *MySQL {

	// Get MySQL connection string
	return &MySQL{Connection: connection(), Database: nil}
}

// GetTables gets all the tables names from this data source
func (m *MySQL) GetTables() []string {
	keys := []string{}
	for k := range config.Strategy.Tables {
		keys = append(keys, k)
	}
	return keys
}

// GetData returns the raw data from the tableName
func (m *MySQL) GetData(modelName string) (*sql.Rows, error) {

	// Get the corresponding table to the modelName
	tableName := config.GetTableName(modelName)
	tableFields := config.GetTableFields(modelName) // map[string]string

	// open a connection
	db, err := m.open()
	if err != nil {
		return nil, err
	}
	defer m.close(db)

	// Fields for that table
	f := make([]string, 0, len(tableFields))
	for _, value := range tableFields {
		if value != "" {
			f = append(f, value)
		}
	}

	fields := strings.Join(f, ", ")

	// Get only the fields that we are going to use
	// the query string . To Do. Select only the stuff you are going to use
	query := strings.Join([]string{"SELECT", fields, "from", tableName}, " ")
	sd, err := db.Query(query)
	if err != nil {
		return nil, queryError{query: query, error: err}
	}

	return sd, nil
}

//////* Not exported functions *//////

// Returns the connection string
func connection() string {
	return credential.Username + ":" + credential.Password + "@" + "/" + credential.Database
}

// Open gives back a pointer to the DB
func (m *MySQL) open() (*sql.DB, error) {

	database, err := sql.Open("mysql", m.Connection)
	if err != nil {
		return nil, connectError{connection: m.Connection}
	}

	if err = database.Ping(); err != nil {
		return nil, connectError{connection: m.Connection}
	}

	m.Database = database

	return database, nil
}

// Close closes the db
func (m MySQL) close(db *sql.DB) error {
	return db.Close()
}
