/*
Package main
Import external source database into local source and transform it
*/
package main

import (
	"fmt"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"

	"github.com/coralproject/sponge/pkg/coral"
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
}

func main() {

	log.Dev("startup", "main", "Start")

	// Connect to external source
	mysql, err := source.New("mysql") // To Do. 1. Needs to ensure maximum rate limit is not reached
	if err != nil {
		log.Error("startup", "main", err, "Connect to external MySQL")
		return
	}

	// Get All the tables's names that we have in the strategy json file
	tables, err := mysql.GetTables()

	if err != nil {
		log.Error("startup", "main", err, "Get external MySQL tables")
		return
	}

	for _, modelName := range tables {
		// Get the data
		log.User("main", "main", "### Getting data '%s' from external source.\n", modelName)
		data, err := mysql.GetData(modelName)
		if err != nil {
			log.Error("main", "main", err, "Get external MySQL data")
			continue
		}

		//Transform the data row by row
		log.User("main", "main", "# Transforming data to the coral schema.\n")
		// Loop on all the data
		for _, row := range data {
			newRow, err := fiddler.TransformRow(row, modelName)
			if err != nil {
				log.Error("main", "main", err, "Error when transforming the row %s.", row)
			}

			fmt.Println("####### ADD ROW: ", string(newRow))
			// send the row to pillar
			err = coral.AddRow(newRow, modelName)
			if err != nil {
				log.Error("main", "main", err, "Error when adding the row %s.", row)
			}
		}
	}
	log.Dev("shutdown", "main", "Complete")
}
