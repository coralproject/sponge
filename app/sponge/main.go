/*
Package main
Import external source database into local source and transform it
*/
package main

import (
	"os"
	"time"

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

	//Get All the tables's names that we have in the strategy json file
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
			//continue
		}

		//Transform the data row by row
		log.User("main", "main", "# Transforming data to the coral schema.\n")
		// Loop on all the data

		// initialize benchmarking for current table
		start := time.Now()
		blockStart := time.Now()
		blockSize := int64(1000) // number of documents between each report
		documents := int64(0)
		totalDocuments := int64(len(data))

		for _, row := range data {

			// output benchmarking for each block of documents
			if documents%blockSize == 0 && documents > 0 {

				// calculate stats
				percentComplete := float64(documents) / float64(totalDocuments) * float64(100)
				msSinceStart := time.Since(start).Nanoseconds() / int64(1000000)
				msSinceBlock := time.Since(blockStart).Nanoseconds() / int64(1000000)
				timeRemaining := int64(float64(time.Since(start).Seconds()) / float64(percentComplete) * float64(100))

				// log stats
				log.User("import", "", "%v%% (%v/%v imported) %vms, %vms avg - last %v in %vms, %vms avg -- est time remaining %vs\n", int64(percentComplete), documents, totalDocuments, msSinceStart, msSinceStart/documents, blockSize, msSinceBlock, msSinceBlock/blockSize, int64(timeRemaining))
				blockStart = time.Now()

			}
			documents = documents + 1

			// transform the row
			newRow, err := fiddler.TransformRow(row, modelName)
			if err != nil {
				log.Error("main", "main", err, "Error when transforming the row %s.", row)
			}

			// send the row to pillar
			err = coral.AddRow(newRow, modelName)
			if err != nil {
				log.Error("main", "main", err, "Error when adding the row %s.", row)
			}
		}
	}
	log.Dev("shutdown", "main", "Complete")
}
