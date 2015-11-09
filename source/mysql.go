/*
Package source implements a way to get data from external MySQL sources.

External possible sources:
* MySQL
* API

*/
package source

import (
	"database/sql"
	"log"
	"strings"

	configuration "github.com/coralproject/sponge/config"
	"github.com/coralproject/sponge/models"
	"github.com/coralproject/sponge/utils"

	_ "github.com/go-sql-driver/mysql" // Check if this can be imported not blank. To Do.
)

// Global configuration variables that holds all the configuration from the config file
var config = *configuration.New() // Reads the configuration file
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
func (m *MySQL) GetData(tableName string) utils.Data {

	var d utils.Data

	db, err := m.open()
	if err != nil {
		log.Fatal("Error when connecting to database. ", err)
		d.Error = err
	}
	defer m.close(db)

	queryString := strings.Join([]string{"SELECT * from", tableName, "limit 10"}, " ")

	// Returns data into a map that is a json structure
	d.Rows, err = runQuery(db, tableName, queryString)
	d.Type = tableName
	if err != nil {
		log.Fatal("Error when quering the DB ", err)
	}
	return d
}

func runQuery(db *sql.DB, table string, query string) ([]models.Model, error) {

	var m models.Model
	var ms []models.Model

	m = utils.New(table)

	sd, err := db.Query(query)
	if err != nil {
		log.Fatal("Error when quering the DB ", err)
		return nil, err
	}
	defer sd.Close()

	ms, err = m.ProcessData(sd)

	return ms, err
}

//////* Not exported functions *//////

// Returns the connection string
func connection() string {
	return credential.Username + ":" + credential.Password + "@" + "/" + credential.Database
}

// Open gives back a pointer to the DB
func (m *MySQL) open() (*sql.DB, error) {

	var err error
	m.Database, err = sql.Open("mysql", m.Connection)
	if err != nil {
		log.Fatal("Could not connect to MySQL database with ", m.Connection, err)
		return nil, err
	}

	err = m.Database.Ping()
	if err != nil {
		log.Fatal("Could not connect to the database with ", m.Connection, err)
		return nil, err
	}

	return m.Database, nil
}

// Close closes the db
func (m MySQL) close(db *sql.DB) error {
	return db.Close()
}

// Get returns data from the query to the db
func (m MySQL) get(db *sql.DB, query string) *sql.Rows {

	// LOOK INTO config.Strategy to see which is the strategy to follow
	d, err := db.Query(query)
	if err != nil {
		log.Fatal("Error when quering the DB ", err)
	}

	// To Do: it needs to return DATA type
	return d
}
