/*
Package main
Import external source database into local source and transform it
*/
package main

import (
	"flag"
	"os"
	"time"

	"github.com/ardanlabs/kit/cfg"
	"github.com/coralproject/sponge/pkg/coral"
	"github.com/coralproject/sponge/pkg/fiddler"
	"github.com/coralproject/sponge/pkg/log"
	"github.com/coralproject/sponge/pkg/report"
	"github.com/coralproject/sponge/pkg/source"
)

// Limit on query
var limitFlag int
var offsetFlag int

// Order by field
var orderbyFlag string

// Import from report on failed records (or not)
var importFailedFlag bool

const (
	limitDefault        = 9999999999
	offsetDefault       = 0
	orderbyDefault      = ""
	importFailedDefault = false
)

// Initialize log, get flag variables, initialize report
func init() {

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	flag.IntVar(&limitFlag, "limit", limitDefault, "-limit= Number of rows that we are going to import at a time")
	flag.IntVar(&offsetFlag, "offset", offsetDefault, "-offset= Offset for the sql query")
	flag.StringVar(&orderbyFlag, "orderby", orderbyDefault, "-orderby= Order by field of the query on external source")
	flag.BoolVar(&importFailedFlag, "onlyfails", importFailedDefault, "-onlyfails Import only the failed documents recorded in report")

	flag.Parse()

	// Initialize the report and write it down at the end (it does not create the file until the end)
	report.Init()
}

func process(modelName string, data []map[string]interface{}) {
	//Transform the data row by row
	log.User("main", "process", "# Transforming data to the coral schema.\n")
	log.User("main", "process", "# And importing %v documents.", len(data))
	// Loop on all the data}

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
			log.User("main", "process", "%v%% (%v/%v imported) %vms, %vms avg - last %v in %vms, %vms avg -- est time remaining %vs\n", int64(percentComplete), documents, totalDocuments, msSinceStart, msSinceStart/documents, blockSize, msSinceBlock, msSinceBlock/blockSize, int64(timeRemaining))
			blockStart = time.Now()

		}
		documents = documents + 1

		// transform the row
		newRow, err := fiddler.TransformRow(row, modelName)
		id := fiddler.GetID(modelName)
		if err != nil {
			log.Error("main", "process", err, "Error when transforming the row %s.", row)

			//RECORD to report about failing transformation
			report.Record(modelName, row[id], row, "Failing transform data", err)
		}

		// send the row to pillar
		err = coral.AddRow(newRow, modelName)
		if err != nil {
			log.Error("main", "process", err, "Error when adding the row %s.", row)
			//RECORD to report about failing adding row to coral db
			report.Record(modelName, row[id], row, "Failing add row to coral", err)
		}
	}
}

func main() {

	log.Dev("main", "main", "Start")

	// Connect to external source
	log.User("main", "main", "### Connecting to external database...")

	mysql, err := source.New("mysql") // To Do. 1. Needs to ensure maximum rate limit is not reached
	if err != nil {
		log.Error("main", "main", err, "Connect to external MySQL")
		return
	}

	if importFailedFlag { // import only what is in the report of failed imported
		log.User("main", "main", "Reading file of data to import.")

		// get the data that needs to be imported
		rowsToImport, err := report.ReadReport() //[]map[string]interface{}
		if err != nil {
			log.Error("main", "main", err, "Getting the rows that will be imported")
		}

		var data []map[string]interface{}
		for _, row := range rowsToImport {
			table := row["table"].(string)
			if len(row["ids"].([]string)) < 1 {
				// Get the data
				log.User("main", "main", "### Reading data from table '%s'. \n", table)
				data, err = mysql.GetData(table, offsetFlag, limitFlag, orderbyFlag)
			} else {
				log.User("main", "main", "### Reading data from table '%s', quering '%s'. \n", table, row["ids"])
				data, err = mysql.GetQueryData(table, offsetFlag, limitFlag, orderbyFlag, row["ids"].([]string))
			}
			if err != nil {
				report.Record(table, row["ids"], row, "Failing getting data", err)
			}

			// transform and get data into pillar
			process(table, data)
		}
	} else { // import everything that is in the strategy
		//Get All the tables's names that we have in the strategy json file
		tables, err := mysql.GetTables()
		if err != nil {
			log.Error("main", "main", err, "Get external MySQL tables")
			return
		}

		for _, modelName := range tables {

			// Get the data
			log.User("main", "main", "### Reading data from table '%s'. \n", modelName)
			data, err := mysql.GetData(modelName, offsetFlag, limitFlag, orderbyFlag)
			if err != nil {
				log.Error("main", "main", err, "Get external MySQL data")
				//RECORD to report about failing modelName
				report.Record(modelName, "", nil, "Failing getting data", err)
				continue
			}

			//transform and send to pillar
			process(modelName, data)
		}

	}

	// Write report on failures (if any)
	report.Write()

	log.Dev("shutdown", "main", "Complete")
}
