/*
Package main
Import external source database into local source and transform it
*/
package main

import (
	"flag"
	"os"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"

	"github.com/coralproject/sponge/pkg/coral"
	"github.com/coralproject/sponge/pkg/fiddler"
	"github.com/coralproject/sponge/pkg/source"
)

var limitFlag int

func init() {
	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	flag.IntVar(&limitFlag, "count", 10, "Number of rows that we are going to import at a time")
}

func main() {

	log.Dev("startup", "main", "Start")

	// Connect to external source
	log.User("main", "main", "### Connecting to external database...")
	mysql, err := source.New("mysql") // To Do. 1. Needs to ensure maximum rate limit is not reached
	if err != nil {
		log.Error("startup", "main", err, "Connect to external MySQL")
		return
	}

	//Get All the tables's names that we have in the strategy json file
	tables, err := mysql.GetTables()
	if err != nil {
		log.Error("startup", "main", err, "Get external MySQL tables")
		return
	}

	for _, modelName := range tables {

		// Get the data
		log.User("main", "main", "### Reading data from table '%s'. \n", modelName)
		data, err := mysql.GetData(modelName, limitFlag)
		if err != nil {
			log.Error("main", "main", err, "Get external MySQL data")
			//continue
		}

		total := len(data)
		countRow := 0
		//Transform the data row by row
		log.User("main", "main", "# Transforming data to the coral schema.\n")
		log.User("main", "main", "# And importing %v documents.", len(data))
		// Loop on all the data
		for _, row := range data {
			countRow = countRow + 1
			newRow, err := fiddler.TransformRow(row, modelName)
			if err != nil {
				log.Error("main", "main", err, "Error when transforming the row %s.", row)
			}

			// send the row to pillar
			log.User("main", "main", "# %v/%v documents completed.", countRow, total)
			err = coral.AddRow(newRow, modelName)
			if err != nil {
				log.Error("main", "main", err, "Error when adding the row %s.", row)
			}
		}
	}
	log.User("main", "main", "### Complete on %v seconds")

	log.Dev("shutdown", "main", "Complete")
}
