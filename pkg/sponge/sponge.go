/*Package sponge imports external source database into local source, transform it and send it to the coral system (called pillar).*/
package sponge

import (
	"time"

	"github.com/coralproject/sponge/pkg/coral"
	"github.com/coralproject/sponge/pkg/fiddler"
	"github.com/coralproject/sponge/pkg/log"
	"github.com/coralproject/sponge/pkg/report"
	"github.com/coralproject/sponge/pkg/source"
)

// Import gets data, transform it and send it to pillar
func Import(limit int, offset int, orderby string, table string, importonlyfailed string, errorsfile string) {

	// Initialize the report and write it down at the end (it does not create the file until the end)
	report.Init(errorsfile)
	// Connect to external source
	log.User("main", "import", "### Connecting to external database...")

	// Initialize the source
	source.Init()
	mysql, err := source.New("mysql") // To Do. 1. Needs to ensure maximum rate limit is not reached
	if err != nil {
		log.Error("sponge", "import", err, "Connect to external MySQL")
	}

	fiddler.Init()
	coral.Init()

	if importonlyfailed != "" { // import only what is in the report of failed importeda
		importOnlyFailedRecords(mysql, limit, offset, orderby, importonlyfailed)
	} else { // import everything that is in the strategy
		if table != "" {
			importTable(mysql, limit, offset, orderby, table)
		} else {
			importAll(mysql, limit, offset, orderby)
		}
	}

}

// CreateIndex will read the strategy file and create index that are mentioned there for each collection
func CreateIndex(collection string) {

	log.User("main", "createindex", "###  Create Index.")

	if collection == "" {
		//create index for everybody

		// get data from strategy file
		tables := fiddler.GetCollections()
		// for each table
		for t := range tables {
			log.User("main", "createindex", "### Index for collection %s.", tables[t])
			coral.CreateIndex(tables[t])
		}
		return
	}

	log.User("main", "createindex", "### Index for collection %s.", collection)
	//create index only for collection
	coral.CreateIndex(collection)
}

// Import gets data from report on failed import, transform it and send it to pillar
func importOnlyFailedRecords(mysql source.Sourcer, limit int, offset int, orderby string, importonlyfailed string) {

	log.User("sponge", "importOnlyFailedRecords", "### Reading file of data to import.")

	// get the data that needs to be imported
	rowsToImport, err := report.ReadReport(importonlyfailed) //[]map[string]interface{}
	if err != nil {
		log.Error("sponge", "importOnlyFailedRecords", err, "Getting the rows that will be imported")
	}

	var data []map[string]interface{}
	for _, row := range rowsToImport {
		table := row["table"].(string)
		if len(row["ids"].([]string)) < 1 {
			// Get the data
			log.User("sponge", "importOnlyFailedRecords", "### Reading data from table '%s'. \n", table)
			data, err = mysql.GetData(table, offset, limit, orderby)
		} else {
			log.User("sponge", "importOnlyFailedRecords", "### Reading data from table '%s', quering '%s'. \n", table, row["ids"])
			data, err = mysql.GetQueryData(table, offset, limit, orderby, row["ids"].([]string))
		}
		if err != nil {
			report.Record(table, row["ids"], row, "Failing getting data", err)
		}

		// transform and get data into pillar
		process(table, data)
	}
}

// Import gets ALL data, transform it and send it to pillar
func importAll(mysql source.Sourcer, limit int, offset int, orderby string) {

	log.User("sponge", "importAll", "### Reading tables to import from strategy file.")

	//Get All the tables's names that we have in the strategy json file
	tables, err := mysql.GetTables()
	if err != nil {
		log.Error("sponge", "importAll", err, "Get external MySQL tables")
		return
	}
	for _, modelName := range tables {

		// Get the data
		log.User("sponge", "importAll", "### Reading data from table '%s'. \n", modelName)
		data, err := mysql.GetData(modelName, offset, limit, orderby)
		if err != nil {
			log.Error("sponge", "importAll", err, "Get external MySQL data")
			//RECORD to report about failing modelName
			report.Record(modelName, "", nil, "Failing getting data", err)
			continue
		}

		//transform and send to pillar
		process(modelName, data)
	}
}

// ImportTable gets ony data related to table, transform it and send it to pillar
func importTable(mysql source.Sourcer, limit int, offset int, orderby string, modelName string) {

	// Get the data
	log.User("sponge", "importTable", "### Reading data from table '%s'. \n", modelName)
	data, err := mysql.GetData(modelName, offset, limit, orderby)
	if err != nil {
		log.Error("sponge", "importAll", err, "Get external MySQL data")
		//RECORD to report about failing modelName
		report.Record(modelName, "", nil, "Failing getting data", err)
		return
	}

	//transform and send to pillar
	process(modelName, data)
}

func process(modelName string, data []map[string]interface{}) {
	//Transform the data row by row
	log.User("sponge", "process", "# Transforming data to the coral schema.\n")
	log.User("sponge", "process", "# And importing %v documents.", len(data))
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
			log.User("sponge", "process", "%v%% (%v/%v imported) %vms, %vms avg - last %v in %vms, %vms avg -- est time remaining %vs\n", int64(percentComplete), documents, totalDocuments, msSinceStart, msSinceStart/documents, blockSize, msSinceBlock, msSinceBlock/blockSize, int64(timeRemaining))
			blockStart = time.Now()

		}
		documents = documents + 1

		// transform the row
		newRow, err := fiddler.TransformRow(row, modelName)
		id := fiddler.GetID(modelName)
		if err != nil {
			log.Error("sponge", "process", err, "Error when transforming the row %s.", row)

			//RECORD to report about failing transformation
			report.Record(modelName, row[id], row, "Failing transform data", err)
		}

		// To Do: acquire meta-data
		/*
		   hit API
		   sponge.API.GetData(row)
		   store result in newrow.metadata
		*/

		log.Dev("sponge", "process", "Transform: %v -> %v", row, string(newRow))

		// send the row to pillar
		err = coral.AddRow(newRow, modelName)
		if err != nil {
			log.Error("sponge", "process", err, "Error when adding the row %s.", row)
			//RECORD to report about failing adding row to coral db
			report.Record(modelName, row[id], row, "Failing add row to coral", err)
		}
	}
}
