/*
Package main

Import source database into mongodb

*/
package main

import (
	// It is only being used when defining the NullStrings in the struct

	"flag"
	"fmt"
	"log"
	// It should be our logger

	"github.com/coralproject/sponge/localDB"
	"github.com/coralproject/sponge/source"
)

func main() {

	/* Arguments for the command line */
	// I want to be able to run the program dry (no insert into local db)
	var dry bool
	flag.BoolVar(&dry, "dry", false, "a bool") // To Do

	/* The external data source */

	var mysql *source.MySQL
	mysql = source.NewSource() // To Do. 1. Needs to ensure maximum rate limit is not reached

	/* The local data source */
	var mongo *localDB.MongoDB
	mongo = localDB.NewLocalDB()

	/* EXTRACT DATA */

	// Get All the tables from the MySQL
	tables := mysql.GetTables()

	//var data utils.Data
	for collectionName, tableName := range tables {
		fmt.Printf("Connecting to external source to get updated data for %s. \n", tableName)
		data := mysql.GetData(tableName, collectionName)
		if data.Error != nil {
			fmt.Println("Error when getting data from ", tableName)
		}

		fmt.Printf("Connecting to local database to get push data into %s. \n", collectionName)
		err := mongo.Add(collectionName, data.Rows, dry) //push data into mongodb local collection <-- go routine
		if err != nil {
			log.Fatal("It was not able to push data into Mongodb. ", err)
		}
	}

}
