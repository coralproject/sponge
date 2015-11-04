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

	/* Arguments for the command line */
	// I want to be able to run the program dry (no insert into local db)
	var dry bool
	flag.BoolVar(&dry, "dry", false, "a bool") // To Do

	/* Connects into mysql database and retrieve all data */

	var mysql *source.MySQL
	// To Do. 1. Needs to ensure maximum rate limit is not reached
	mysql = source.NewSource()

	// Print message
	fmt.Println("Connecting to external source to get updated data: ...")

	// To Do 2. Determine which slice of data to get next
	// To Do 3. Use the strategy to request the slice (either db query or api call)

	/* EXTRACT DATA */

	// We are using d to store data from the source
	var d utils.Data

	// Get the last data from external source (right now is getting all the data from table comments on external source, it needs to look at the strategy to know which table to bring in)
	d = mysql.GetNewData()
	if d.Error != nil {
		log.Fatal("Error when querying the external database", d.Error)
	}

	// Print message
	fmt.Printf("Got %d rows from external source. \n", len(d.Comments))

	// Print message
	fmt.Println("Connecting to local database to insert data: ... ")

	// the database we are importing into
	mongo := localDB.NewLocalDB()

	/* Open a session for MongoDB */
	err := mongo.Open()
	if err != nil {
		log.Fatal("Error when connecting to MongoDB.", err)
	}

	/* Close the session when it is done */
	defer mongo.Close()

	/* INSERT DATA */

	// To Do: do this for all the collections/tables in the configuration

	// Inserts all the documents into the collection Comments
	fmt.Printf("Inserting %d comments...\n", len(d.Comments))
	errMo := mongo.Add(d, false) // To Do: It should have dry variable instead of false.
	if errMo != nil {
		log.Fatal("Error when inserting data into local db. ", errMo)
	}

	fmt.Println("Done.")

}
