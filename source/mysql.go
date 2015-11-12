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

	"github.com/coralproject/sponge/models"
	"github.com/coralproject/sponge/utils"

	_ "github.com/go-sql-driver/mysql" // Check if this can be imported not blank. To Do.
)

// Global configuration variables that holds the credentials for mysql
var credential = config.GetCredential("mysql")

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
	//connection := mysqlConnection(config)

	// Get MySQL connection string
	return &MySQL{Connection: connection(), Database: nil}
}

// GetTables gets all the tables names from this data source
func (m *MySQL) GetTables() map[string]string {
	return config.Strategy.Tables
}

// GetData returns the raw data from the tableName
func (m *MySQL) GetData(tableName string, modelName string) utils.Data {

	var d utils.Data
	d.Type = modelName

	db, err := m.open()
	if err != nil {
		d.Error = err
		return d
	}
	defer m.close(db)

	queryString := strings.Join([]string{"SELECT * from", tableName}, " ")

	// Returns data into a map that is a json structure
	d.Rows, err = runQuery(db, modelName, queryString)
	if err != nil {
		d.Error = err
	}

	return d
}

// Run the query on the db
func runQuery(db *sql.DB, model string, query string) ([]models.Model, error) {

	var m models.Model
	var ms []models.Model

	// Creates a Model of the type model. A model struct.
	m, err := models.New(model)
	if err != nil {
		return nil, modelError{model: model}
	}

	sd, err := db.Query(query)
	if err != nil {
		return nil, queryError{query: query}
	}
	defer sd.Close()

	ms, err = m.Transform(sd)

	return ms, err
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

// // Get returns data from the query to the db
// func (m MySQL) get(db *sql.DB, query string) (*sql.Rows, error) {
//
// 	// LOOK INTO config.Strategy to see which is the strategy to follow
// 	d, err := db.Query(query)
// 	if err != nil {
// 		log.Fatal("Error when quering the DB ", err)
// 		return nil, queryError{query: query}
// 	}
//
// 	// To Do: it needs to return DATA type
// 	return d, nil
// }
