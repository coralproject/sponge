package fiddler_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/coralproject/sponge/fiddler"
	"github.com/coralproject/sponge/models"
	"github.com/coralproject/sponge/utils"
)

// Stubs - AHHHHH! Blocked on how to write stubs
// type fakeConfig struct {
// 	Name        string
// 	Strategy    config.Strategy
// 	Credentials []config.Credential
// }
//
// func New() (*fakeConfig, error) {
// 	fake :=  fakeConfig{} //{Name: "test", Strategy: config.Strategy{}, Credentials: []config.Credential{}}
// 	return fakeC, nil
// }
//
// func GetCredentials() ([]config.Credential, error) {
// 	credentials := []config.Credential
// 	return credentials, nil
// }

// Testing

func TestNewLocalDB(t *testing.T) {

	// Stubs config.GetCredentials and mongoDBConnection()

	mongo := fiddler.NewLocalDB()
	// Expect mongodb Type MongoDB struct
	var typeMongoDB = reflect.TypeOf((*fiddler.MongoDB)(nil)).Elem()
	if reflect.TypeOf(*mongo) == typeMongoDB {
		t.Error("Expect a mongoDB struct to be type ", typeMongoDB)
	}
}

func TestNewLocalDBBadConfig(t *testing.T) {
	// Error when there is no configuration file
	// config.GetCredentials gives an error
	t.Skip("Pass.")
}

func TestNewLocalDBBadMongoDbConnection(t *testing.T) {
	// Error when there is no field for mongodb connection in configuration file
	// mongoDBConnection(c) gives an error
	t.Skip("Pass.")
}

func TestAdd(t *testing.T) {
	t.Skip("Pass.")
}

/* ON HOW TO USE THIS PACKAGE */

// ExampleMongoDB on how to use the MongoDB
func ExampleMongoDB() {

	// Connects into mongo database
	mongo := fiddler.NewLocalDB()

	var d utils.Data
	d.Comments = []models.Comment{{CommentID: 1}, {CommentID: 2}}

	// Inserts all the documents into the collection CommentsExample (dry in false)
	errMo := mongo.Add(d, false)
	if errMo != nil {
		log.Fatal("Error when inserting data into local db. ", errMo)
	}

	fmt.Println("Done.")
}
