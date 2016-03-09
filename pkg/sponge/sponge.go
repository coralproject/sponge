/*Package sponge imports external source database into local source, transform it and send it to the coral system (called pillar).*/
package sponge

import (
	"strings"
	"sync"
	"time"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/sponge/pkg/coral"
	"github.com/coralproject/sponge/pkg/fiddler"
	"github.com/coralproject/sponge/pkg/report"
	"github.com/coralproject/sponge/pkg/source"
)

const (
	// VersionNumber is the version for sponge
	VersionNumber = 0.1
)

var (
	dbsource source.Sourcer
	uuid     string
)

// Init initialize the packages that are going to be used by sponge
func Init(u string) error {

	uuid = u
	var err error

	// Initialize the source
	foreignSource := source.Init(uuid)
	dbsource, err = source.New(foreignSource) // To Do. 1. Needs to ensure maximum rate limit is not reached
	if err != nil {
		log.Error(uuid, "sponge.import", err, "Connect to external Database")
		return err
	}

	fiddler.Init(uuid)
	coral.Init(uuid)

	return err
}

// Import gets data, transform it and send it to pillar
func Import(limit int, offset int, orderby string, types string, importonlyfailed string, errorsfile string) {

	// Initialize the report and write it down at the end (it does not create the file until the end)
	report.Init(uuid, errorsfile)

	// Connect to external source
	log.User(uuid, "sponge.import", "### Connecting to external database...")

	if importonlyfailed != "" { // import only what is in the report of failed importeda
		importOnlyFailedRecords(dbsource, limit, offset, orderby, importonlyfailed)
	} else { // import everything that is in the strategy
		if types != "" {
			for _, t := range strings.Split(types, ",") {
				importType(dbsource, limit, offset, orderby, strings.Trim(t, " ")) // removes any extra space
			}
		} else {
			importAll(dbsource, limit, offset, orderby)
		}
	}

}

// CreateIndex will read the strategy file and create index that are mentioned there for each collection
func CreateIndex(collection string) {

	log.User(uuid, "sponge.createindex", "###  Create Index.")

	//create index for everybody
	if collection == "" {
		// get data from strategy file
		tables := fiddler.GetCollections()

		// for each table
		for t := range tables {
			log.User(uuid, "sponge.createindex", "### Create index for collection %s.", tables[t])
			coral.CreateIndex(tables[t])
		}
		return
	}

	log.User(uuid, "sponge.createindex", "### Create index for collection %s.", collection)
	//create index only for collection
	coral.CreateIndex(collection)
}

// Import gets data from report on failed import, transform it and send it to pillar
func importOnlyFailedRecords(dbsource source.Sourcer, limit int, offset int, orderby string, importonlyfailed string) {

	log.User(uuid, "sponge.importOnlyFailedRecords", "### Reading file of data to import.")

	// get the data that needs to be imported
	rowsToImport, err := report.ReadReport(importonlyfailed) //[]map[string]interface{}
	if err != nil {
		log.Error(uuid, "sponge.importOnlyFailedRecords", err, "Getting the rows that will be imported")
	}

	var data []map[string]interface{}
	for _, row := range rowsToImport {
		table := row["table"].(string)
		if len(row["ids"].([]string)) < 1 {
			// Get the data
			log.User(uuid, "sponge.importOnlyFailedRecords", "### Reading data for table '%s'. \n", table)
			data, err = dbsource.GetData(table, offset, limit, orderby)
		} else {
			log.User(uuid, "sponge.importOnlyFailedRecords", "### Reading data for table '%s', quering '%s'. \n", table, row["ids"])
			data, err = dbsource.GetQueryData(table, offset, limit, orderby, row["ids"].([]string))
		}
		if err != nil {
			report.Record(table, row["ids"], row, "Failing getting data", err)
		}

		// transform and get data into pillar
		process(table, data)
	}
}

// Import gets ALL data, transform it and send it to pillar
func importAll(dbsource source.Sourcer, limit int, offset int, orderby string) {

	log.User(uuid, "sponge.importAll", "### Reading tables to import from strategy file.")

	//Get All the tables's names that we have in the strategy json file
	tables, err := dbsource.GetTables()
	if err != nil {
		log.Error(uuid, "sponge.importAll", err, "Get external tables")
		return
	}
	for _, modelName := range tables {

		// Get the data
		log.User(uuid, "sponge.importAll", "### Reading data from table '%s'. \n", modelName)

		// get only some data at a time
		CHUNK := 1000
		i := offset
		n := CHUNK
		var wg sync.WaitGroup

		for i < limit { // TO DO: limit could be too big
			wg.Add(1)

			go func() {
				defer wg.Done()
				data, err := dbsource.GetData(modelName, i, n, orderby)
				if err != nil {
					log.Error(uuid, "sponge.importAll", err, "Get external data for table %s.", modelName)
					//RECORD to report about failing modelName
					report.Record(modelName, "", nil, "Failing to get data.", err)
					return
				}

				//transform and send to pillar
				process(modelName, data)
			}()

			i = n + 1
			n = i + n
		}

		wg.Wait()
	}
}

// ImportType gets ony data related to table, transform it and send it to pillar
func importType(dbsource source.Sourcer, limit int, offset int, orderby string, modelName string) {

	// Get the data
	log.User(uuid, "sponge.importTable", "### Reading data from table '%s'.", modelName)

	data, err := dbsource.GetData(modelName, offset, limit, orderby)
	if err != nil {
		log.Error(uuid, "sponge.importAll", err, "Get external data for table %s.", modelName)
		//RECORD to report about failing modelName
		report.Record(modelName, "", nil, "Failing to get data", err)
		return
	}

	// Transform and send to pillar
	process(modelName, data)

}

func process(modelName string, data []map[string]interface{}) {
	// Transform the data row by row
	log.User(uuid, "sponge.process", "# Transforming data to the coral schema.\n")
	log.User(uuid, "sponge.process", "# And importing %v documents.", len(data))

	// Initialize benchmarking for current table
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

			//log stats
			log.User(uuid, "sponge.process", "%v%% (%v/%v imported) %vms, %vms avg - last %v in %vms, %vms avg -- est time remaining %vs\n", int64(percentComplete), documents, totalDocuments, msSinceStart, msSinceStart/documents, blockSize, msSinceBlock, msSinceBlock/blockSize, int64(timeRemaining))
			blockStart = time.Now()

		}
		documents = documents + 1

		// transform the row
		id, newRows, err := fiddler.TransformRow(row, modelName)
		if err != nil {
			log.Error(uuid, "sponge.process", err, "Error when transforming the row %s.", row)
			//RECORD to report about failing transformation
			report.Record(modelName, id, row, "Failing transform data", err) // TO DO, needs to recalculate id
		}

		// To Do: acquire meta-data
		/*
		   hit API
		   sponge.API.GetData(row)
		   store result in newrow.metadata
		*/

		// Usually newRows only will have a document but in the case that we have subcollections
		// we may get more than one document from a transformation
		for _, newRow := range newRows {

			log.Dev(uuid, "sponge.process", "Transforming: %v into %v.", row, newRow)

			// send the row to pillar
			err = coral.AddRow(newRow, modelName)
			if err != nil {
				log.Error(uuid, "sponge.process", err, "Error when adding a row") // thae row %v to %s.", string(newRow), modelName)
				//RECORD to report about failing adding row to coral db
				report.Record(modelName, id, row, "Failing add row to coral", err)
			}
		}
	}
}
