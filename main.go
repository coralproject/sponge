/*
Package main

Import external source database into local source

*/
package main

import (
	"fmt"
	"log"
	"sync" // It should be our logger

	"github.com/coralproject/sponge/fiddler"
	"github.com/coralproject/sponge/source"
)

//* Errors used in this package *//

// When trying to export data from modelName
type dataExportError struct {
	modelName string
}

func (e dataExportError) Error() string {
	return fmt.Sprintf("Error when getting mysql data from table %s.", e.modelName)
}

// When trying to import data into shelfdb
type dataImportError struct {
	error string
}

func (e dataImportError) Error() string {
	return fmt.Sprintf("It was not able to push data into Mongodb. Error: %s.", e.error)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
}

// Import into modelName
func importit(modelName string, mysql *source.MySQL, fiddler *fiddler.ShelfDB) {

	// Get data from mysql
	log.Printf("Connecting to external source to get updated data for %s. \n", modelName)
	data := mysql.GetData(modelName)
	if data.Error != nil {
		log.Fatalf("Error when getting mysql data from table %s. ", modelName)
	}

	// To Do
	// // Send data into the service layer
	log.Printf("Connecting to local database to get push data into collection %s. \n", modelName)
	err := fiddler.Add(modelName, data.Rows) //push data into the service layer
	if err != nil {
		//log.Fatalf("It was not able to push data into Shelfdb. Error: %s.", err)
		log.Fatal(err.(dataExportError))
	}
}

func main() {

	/* Arguments for the command line */
	// I want to be able to run the program dry (no insert into local db) <--- not sure about this feature
	// var dry bool
	// flag.BoolVar(&dry, "dry", false, "a bool") // To Do
	//log.Printf("Starting main with DRY in %t.", dry)

	/* The external data source */
	var mysql *source.MySQL
	mysql = source.NewSource() // To Do. 1. Needs to ensure maximum rate limit is not reached

	/* The local data source */
	var fiddler *fiddler.ShelfDB
	fiddler = fiddler.NewLocalDB()

	/* EXTRACT DATA */

	log.Printf("Getting data from MYSQL. ")
	// Get All the tables from the MySQL
	tables := mysql.GetTables()

	wg := sync.WaitGroup{}

	fmt.Println(tables)

	//var data utils.Data
	for _, modelName := range tables {
		wg.Add(1)
		go func(modelName string) {
			defer wg.Done()
			log.Printf("### Pushing data into collection %s. ### \n", modelName)

			importit(modelName, mysql, fiddler)
		}(modelName)
	}
	wg.Wait()

}
