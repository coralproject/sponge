package source_test

import (
	"fmt"
	"log"
	"testing"
)

// NewSource returns a new MySQL struct
// Signature: NewSource() (*MySQL, error)
// It depends on the credentials
func TestNewSource(t *testing.T) {
	t.Skip("Not Implemented")

	// fake Credentials

	mysql, err := NewSource()

	// It needs to returns a variable of type
	// type MySQL struct {
	// 	connection string
	// 	database   *sql.DB
	// }

	// mysql.connection has to be type  username:password@/database
	// mysql.database has to be not open

}

func TestGetNewData(t *testing.T) {
	t.Skip("Not Implemented")

	// fake Credentials and fake NewSource
	d := GetNewData()

	// d has to be type utils.Data

}

/* ON HOW TO USE THIS PACKAGE */

// ExampleMySQL on how to use the MySQL
func ExampleMySQL() {

	// Creates a new mysql source
	m, err := NewSource()
	if err != nil {
		log.Fatal("Error when creating new source ", err)
	}

	// Opens a connection to the database (it uses configuration package with a json configuration file)
	db, err := m.Open()
	if err != nil {
		log.Fatal("Error when connection to database with ", m.connection, err)
	}
	// Closes connection when func finishes
	defer m.Close(db)

	// Returns rows related to that query
	d := m.Get(db, "SELECT * FROM nyt_comments LIMIT 1")
	for d.Next() {
		var comment string
		if d.Scan(&comment); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("The comment is '%s'.", comment)
	}
}
