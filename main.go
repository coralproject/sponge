/*
Package main

Import source database into mongodb

*/
package main

import (
	// It is only being used when defining the NullStrings in the struct

	"fmt"
	"log" // It should be our logger

	"github.com/coralproject/mod-data-import/localDB"
	"github.com/coralproject/mod-data-import/source"
	"github.com/coralproject/mod-data-import/utils"
)

func main() {
	// Connects into mysql database and retrieve one row
	var mysql *source.MySQL
	var errMy error

	mysql, errMy = source.NewSource()
	if errMy != nil {
		log.Fatal("Error when creating new source ", errMy)
	}

	var d utils.Data

	// Get the last data from external source
	d = mysql.GetNewData()
	if d.Error != nil {
		log.Fatal("Error when querying the external database", d.Error)
	}

	fmt.Println("Pull out this comment from the external source: \n ", d.Comments[0])

	// Get that comment into a MongoDB collection called "Comments"

	mongo, errMo := localDB.NewLocalDB()
	if errMo != nil {
		log.Fatal("Error when creating new local db. ", errMo)
	}

	errMo = mongo.Add(d)
	if errMo != nil {
		log.Fatal("Error when inserting data into local db. ", errMo)
	}
}
