package localDB_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/coralproject/mod-data-import/localDB"
	"github.com/coralproject/mod-data-import/utils"
)

func TestNewLocalDB(t *testing.T) {
	t.Skip("It shows what to test.")

	//func NewLocalDB() (*MongoDB, error) {

	// Error when there is no configuration file
	// Error when there is no field for mongodb connection in configuration file
	// Fail if connection string has no specific pattern

	// The connection string has always be of the form user:
	//c.Username + ":" + c.Password + "@" + c.Host + "/" + c.Database
}

func TestAdd(t *testing.T) {
	t.Skip("It shows what to test.")
}

/* ON HOW TO USE THIS PACKAGE */

// ExampleMongoDB on how to use the MongoDB
func ExampleMongoDB() {

	// Connects into mongo database
	mongo, errMo := localDB.NewLocalDB()
	if errMo != nil {
		log.Fatal("Error when creating new local db. ", errMo)
	}

	var d utils.Data
	d.Comments = []Comment{{commentid: 1}, {commentid: 2}}

	// Inserts all the documents into the collection Comments
	errMo = mongo.Add(d)
	if errMo != nil {
		log.Fatal("Error when inserting data into local db. ", errMo)
	}

	fmt.Println("Done.")
}
