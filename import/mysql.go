/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API?

*/
package source

import (
	"database/sql"
	"log"

	"github.com/coralproject/mod-data-import/config"
	_ "github.com/go-sql-driver/mysql"
)

/* Implementing the Sources */

//////////// MYSQL SOURCE ////////////

// MySQL is the struct that has the connection string to the external mysql database
type MySQL struct {
	connection string
	database   *sql.DB
}

////// Not exported functions //////

func mysqlConnection(c config.Credentials) string {
	return c.Username + ":" + c.Password + "@" + c.Host + "/" + c.Database
}

////// Exported Functions //////

func NewSource() (*MySQL, error) {
	c, err := config.GetCredentials()

	if err != nil {
		log.Fatal("Error when trying to create new source ", err)
	}

	// Get MySQL connection string
	return &MySQL{connection: mysqlConnection(c)}, err
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
func (m MySQL) Reader(db *sql.DB) Data {
	var d Data
	var err error

	d.rows, err = db.Query("SELECT * FROM nyt_comments")
	if err != nil {
		log.Fatal("Error when quering the DB ", err)
	}

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

	// d := m.Reader()
	// for d.rows.Next() {
	// 	var comment string
	// 	if d.rows.Scan(&comment); err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("The comment is '%s'.", comment)
	// }
}

// Possible Errors
// - TimeOutReader
// - WrongConnection
// - DataErrorReader
