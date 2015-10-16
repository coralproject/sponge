/*
Package source implements a way to get data from external MySQL sources.

External possible sources:
* MySQL
* API?

*/
package source

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/coralproject/mod-data-import/config"
	_ "github.com/go-sql-driver/mysql" // Check if this can be imported not blak. To Do.
)

/* Implementing the Sources */

//////////// MYSQL SOURCE ////////////

// MySQL is the struct that has the connection string to the external mysql database
type MySQL struct {
	connection string
	database   *sql.DB
}

////// Not exported functions //////

func mysqlConnection(credentials []config.Credential) (string, error) {
	// look at the credentials related to mysql
	for i := 0; i < len(credentials); i++ {
		if credentials[i].Adapter == "mysql" {
			c := credentials[i]
			connection := c.Username + ":" + c.Password + "@" + "/" + c.Database
			return connection, nil
		}
	}

	err := fmt.Errorf("Error when trying to get the connection string for mysql.")

	return "", err
}

////// Exported Functions //////

// NewSource returns a new connection
func NewSource() (*MySQL, error) {

	c, err := config.GetCredentials()
	if err != nil {
		log.Fatal("Error when trying to create new source ", err)
	}

	connection, err := mysqlConnection(c)
	if err != nil {
		log.Fatal("Error when trying to get credentials to connect to mysql. ", err)
	}

	// Get MySQL connection string
	return &MySQL{connection: connection}, err
}

// Open gives back a pointer to the DB
func (m MySQL) Open() (*sql.DB, error) {

	db, err := sql.Open("mysql", m.connection)
	if err != nil {
		log.Fatal("Could not connect to MySQL database with ", m.connection, err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Could not connect to the database with ", m.connection, err)
		return nil, err
	}

	m.database = db

	return db, nil
}

// Close closes the db
func (m MySQL) Close(db *sql.DB) error {
	return db.Close()
}

// Reader returns the data needed from mysql db
func (m MySQL) Reader(db *sql.DB) *sql.Rows {

	// LOOK INTO config.Strategy to see which is the strategy to follow
	d, err := db.Query("SELECT * from nyt_comments LIMIT 1")
	if err != nil {
		log.Fatal("Error when quering the DB ", err)
	}

	// To Do: it needs to return DATA type
	return d
}

// ExampleMySQL on how to use the MySQL
func ExampleMySQL() {

	m, err := NewSource()
	if err != nil {
		log.Fatal("Error when creating new source ", err)
	}

	db, err := m.Open()
	if err != nil {
		log.Fatal("Error when connection to database with ", m.connection, err)
	}
	defer m.Close(db)

	d := m.Reader(db)
	for d.Next() {
		var comment string
		if d.Scan(&comment); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The comment is '%s'.", comment)
	}
}

// GetConnection returns the connection string
func (m MySQL) GetConnection() string {
	return m.connection
}

// Possible Errors
// - TimeOutReader
// - WrongConnection
// - DataErrorReader
