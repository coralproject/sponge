/*
Package source implements a way to get data from external sources.

External possible sources:
* MySQL
* API?

*/
package source

import (
	"log"

	"github.com/coralproject/mod-data-import/config"
	"github.com/jinzhu/gorm"
)

////// Structures  //////

// Data is a struct that has all the db rows and error field
type Data struct {
	rows  []byte // Move into appropiate structure
	error error
}

// Source is where the data is coming from (mysql, api)
type Source interface {
	Open() (*gorm.DB, error)
	Close(db *gorm.DB) error
	Reader() Data
	Updated() bool // To Do: Look for a better name - this func check if there's new data in the source db
}

/* Implementing the Sources */

//////////// MYSQL SOURCE ////////////

// MySQL is the struct that has the connection string to the external mysql database
type MySQL struct {
	connection string
}

////// Not exported functions //////

func newMySQL() *MySQL {
	c := config.GetCredentials()
	// Get MySQL connection string
	return &MySQL{connection: mysqlConnection(c)}
}

func mysqlConnection(c config.Credentials) string {
	return c.Username + ":" + c.Password + "@" + c.Host + "/" + c.Database
}

////// Exported Functions //////

// Open gives back a pointer to the DB
func (m MySQL) Open() (*gorm.DB, error) {

	db, err := gorm.Open("MySQL", m.connection)
	if err != nil {
		log.Fatal("Could not connect to MySQL database with ", m.connection, err)
		return nil, err
	}

	err = db.DB().Ping()
	if err != nil {
		log.Fatal("Could not connect to the database with ", m.connection, err)
		return nil, err
	}

	return &db, nil
}

// Close closes the db
func (m MySQL) Close(db *gorm.DB) error {
	return db.Close()
}

// Reader returns the data needed from mysql db
func (m MySQL) Reader() Data {
	var d Data
	return d
}

// ExampleMySQL on how to use the MySQL
func ExampleMySQL() {
	m := newMySQL()

	db, err := m.Open()
	defer m.close(db)

	d := db.Reader()

}

//////////// Disquis API SOURCE ////////////

// DisquisAPI is the struct that has the connection to the Disquis API
type DisquisAPI struct {
	connection string
}

// Reader on the API brings back the data needed from the API
func (a DisquisAPI) Reader() Data {
	var d Data
	return d
}

// Possible Errors
// - TimeOutReader
// - WrongConnection
// - DataErrorReader
