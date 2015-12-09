/*
Package main
Import external source database into local source and transform it
*/
package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/coralproject/shelf/pkg/cfg"
	"github.com/coralproject/shelf/pkg/db/mongo"
	"github.com/coralproject/sponge/pkg/fiddler"
	"github.com/coralproject/sponge/pkg/source"
)

func init() {
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	if err := cfg.Init("SPONGE"); err != nil {
		log.Error("startup", "init", err, "Initializing config")
		os.Exit(1)
	}

	err := mongo.InitMGO()
	if err != nil {
		log.Error("startup", "init", err, "Initializing MongoDB")
		os.Exit(1)
	}
}

func main() {

	// Connect to external source
	mysql, err := source.New("mysql") // To Do. 1. Needs to ensure maximum rate limit is not reached
	if err != nil {
		log.Printf(err.Error())
		return
	}

	// Get All the tables from the MySQL
	tables, err := mysql.GetTables()
	if err != nil {
		log.Printf(err.Error())
		return
	}

	wg := sync.WaitGroup{}

	//var data utils.Data
	for _, modelName := range tables {
		wg.Add(1)
		go func(modelName string) {
			defer wg.Done()

			// Get the data
			log.Printf("### Getting data from external source.\n")
			data, err := mysql.GetData(modelName)
			if err != nil {
				log.Printf(err.Error())
				return
			}

			//Transform the data
			log.Printf("### Transforming data to the coral schema.\n")
			dataCoral, err := fiddler.Transform(modelName, data)
			if err != nil {
				log.Printf(err.Error())
				return
			}

			// var context interface{}
			// var db *db.DB
			log.Printf("### Pushing data into collection %s. ### \n", modelName)
			switch modelName {
			case "User":
				//err = comment.AddUsers(context, db, dataCoral)
				fmt.Printf("Ready to Add Users %s", dataCoral)
			case "Comment":
				//err = comment.AddComments(context, db, dataCoral)
				fmt.Printf("Ready to Add Comments %s", dataCoral)
			}

			if err != nil {
				log.Printf("Error when pushing data into the shelf, collection %s. \n", modelName)
			}

		}(modelName)
	}
	wg.Wait()
}
