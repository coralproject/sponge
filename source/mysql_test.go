package source_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
)

type FakeMySQL struct {
	Connection string
	Database   *sql.DB
}

// NewSource returns a new MySQL struct
// Signature: NewSource() (*MySQL, error)
// It depends on the credentials
func NewSource() (*FakeMySQL, error) {
	var f *FakeMySQL
	f = {
		Connection: "",
		Database: nil
	}
	return f, nil
}

func TestGetNewData(t *testing.T) {

	expectedType := type(utils.Data)
	expectedData := [
		{
			CommentID: 1,
			AssetID: 1,
			StatusID : 0,
			CommentTitle: "Titulo 1",
			CommentBody: "Body 1 "
		}
	]

	// fake Credentials and fake NewSource
	d := GetNewData()

	// d has to be type utils.Data
	if d.type != expectedType {
		t.Fatalf("Expected type to be %s but it was %s", expectedType, d.type)
	}

	// d has to has 2 Comments
	if len(d.Comments) != len(expectedData) {
		f.Fatalf("Expected only %d comments but there were %d.", len(expectedData), len(d.Comments))
	}

	// d.error has to be nil
	if d.Error != nil {
		f.Fatalf("Expected to have no error but there was this one: %s .", d.Error)
	}
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
