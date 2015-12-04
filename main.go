/*
Package main

Import external source database into local source

*/
package main

import (
	"fmt"
	"log"

	"github.com/coralproject/sponge/fiddler"
	"github.com/coralproject/sponge/source"
)

func main() {

	// Connect to external source
	mysql, err := source.NewSource("mysql") // To Do. 1. Needs to ensure maximum rate limit is not reached
	if err != nil {
		log.Printf("Error when connecting to external source: %s", err)
	}
	// // Get All the tables from the MySQL
	// tables := mysql.GetTables()
	//
	// wg := sync.WaitGroup{}
	//
	// //var data utils.Data
	// for _, modelName := range tables {
	// 	wg.Add(1)
	// 	go func(modelName string) {
	// 		defer wg.Done()

	modelName := "User"
	// Get the data
	log.Printf("### Getting data from external source.\n")
	data, err := mysql.GetData(modelName)
	if err != nil {
		log.Printf("Error when getting mysql data from table %s. Error: %s", modelName, err.Error())
		return
	}

	//Transform the data
	log.Printf("### Transforming data to the coral schema.\n")
	dataCoral, err := fiddler.Transform(modelName, data)
	if err != nil {
		log.Printf(err.Error())
		return
	}

	fmt.Println(dataCoral)

	// // Send it to Shelf
	// log.Printf("### Pushing data into collection %s. ### \n", modelName)
	// err = shelf.Add(modelName, dataCoral)
	// if err != nil {
	// 	log.Printf("Error when pushing data into the shelf, collection %s. \n", modelName)
	// 	return
	// }

	// 	}(modelName)
	// }
	// wg.Wait()
}
