/*
Package main

Import source database into mongodb

*/
package main

import (
	// It is only being used when defining the NullStrings in the struct

	"flag"
	"fmt"
	"log" // It should be our logger

	"github.com/coralproject/sponge/localDB"
	"github.com/coralproject/sponge/source"
	"github.com/coralproject/sponge/utils"
)

func main() {

	// I want to be able to run the program dry (no insert into local db)
	dry := flag.Bool("dry", false, "a bool")

	// To Do: Get Strategy with configuration's fields for this phase 1 (tier 1)

	// Connects into mysql database and retrieve all row
	var mysql *source.MySQL

	/* Syncronization Loop */

	// To Do. 1. Needs to ensure maximum rate limit is not reached

	var errMy error

	mysql, errMy = source.NewSource()
	if errMy != nil {
		log.Fatal("Error when creating new source ", errMy)
	}

	fmt.Println("Connecting to external source to get updated data: ...")

	// To Do 2. Determine which slice of data to get next
	// To Do 3. Use the strategy to request the slice (either db query or api call)

	var d utils.Data

	// Get the last data from external source (right now is getting all the data from table comments on external source, it needs to look at the strategy to know which table to bring in)
	d = mysql.GetNewData()
	if d.Error != nil {
		log.Fatal("Error when querying the external database", d.Error)
	}

	fmt.Printf("Got %d rows from external source. \n", len(d.Comments))

	// Connects into mongo database and inserts everything
	fmt.Println("Connecting to local database to insert data: ... ")

	mongo, errMo := localDB.NewLocalDB()
	if errMo != nil {
		log.Fatal("Error when creating new local db. ", errMo)
	}

	// Inserts all the documents into the collection Comments

	fmt.Printf("Inserting %d comments...\n", len(d.Comments))
	errMo = mongo.Add(d, *dry)
	if errMo != nil {
		log.Fatal("Error when inserting data into local db. ", errMo)
	}

	fmt.Println("Done.")
}
