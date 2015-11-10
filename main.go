/*
Package main

Import source database into mongodb

*/
package main

import (
	// It is only being used when defining the NullStrings in the struct

	"flag"
	"log"
	"sync"
	// It should be our logger

	"github.com/coralproject/sponge/localDB"
	"github.com/coralproject/sponge/source"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

// Import from tablename to collectionName
func retrieves(collectionName string, tableName string, mysql *source.MySQL, mongo *localDB.MongoDB, dry bool) {

	log.Printf("Connecting to external source to get updated data for %s. \n", tableName)
	data := mysql.GetData(tableName, collectionName)
	if data.Error != nil {
		log.Fatalf("Error when getting mysql data from table %s. ", tableName)
	}

	log.Printf("Connecting to local database to get push data into collection %s. \n", collectionName)
	err := mongo.Add(collectionName, data.Rows, dry) //push data into mongodb local collection <-- go routine
	if err != nil {
		log.Fatalf("It was not able to push data into Mongodb. Error: %s.", err)
	}
}

func main() {

	/* Arguments for the command line */
	// I want to be able to run the program dry (no insert into local db)
	var dry bool
	flag.BoolVar(&dry, "dry", false, "a bool") // To Do

	log.Printf("Starting main with DRY in %t.", dry)

	/* The external data source */

	var mysql *source.MySQL
	mysql = source.NewSource() // To Do. 1. Needs to ensure maximum rate limit is not reached

	/* The local data source */
	var mongo *localDB.MongoDB
	mongo = localDB.NewLocalDB()

	/* EXTRACT DATA */

	log.Printf("Getting data from MYSQL. ")
	// Get All the tables from the MySQL
	tables := mysql.GetTables()

	wg := sync.WaitGroup{}
	//var data utils.Data
	for collectionName, tableName := range tables {
		wg.Add(1)
		go func(collectionName string, tableName string) {
			defer wg.Done()
			log.Printf("Pushing data %s into collection %s.", tableName, collectionName)
			retrieves(collectionName, tableName, mysql, mongo, dry)
		}(collectionName, tableName)
	}
	wg.Wait()

}
