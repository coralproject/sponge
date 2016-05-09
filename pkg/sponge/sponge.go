/*Package sponge imports external source database into local source, transform it and send it to the coral system (called pillar).*/
package sponge

import (
	"fmt"
	"strings"
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
	options  source.Options
)

// Init initialize the packages that are going to be used by sponge
func Init(u string) error {

	uuid = u
	var err error

	// Initialize the source
	foreignSource, err := source.Init(uuid)
	if err != nil {
		log.Error(uuid, "sponge.import", err, "Initialization of Source")
		return err
	}
	dbsource, err = source.New(foreignSource) // To Do. 1. Needs to ensure maximum rate limit is not reached
	if err != nil {
		log.Error(uuid, "sponge.import", err, "Connect to external Database")
		return err
	}

	fiddler.Init(uuid)
	coral.Init(uuid)

	return err
}

// AddOptions adds flags to the sponge
func AddOptions(limit int, offset int, orderby string, query string, types string, importonlyfailed bool, reportOnFailedRecords bool, reportdbfile string) {
	options = source.Options{
		Limit:                 limit,
		Offset:                offset,
		Orderby:               orderby,
		Query:                 query,
		Types:                 types,
		Importonlyfailed:      importonlyfailed,
		ReportOnFailedRecords: reportOnFailedRecords,
		Reportdbfile:          reportdbfile,
	}
}

// Import gets data, transform it and send it to pillar
func Import() {

	// if there is a flag to start recording a report of failed records, then initialize it
	if options.ReportOnFailedRecords {
		report.Init(uuid, options.Reportdbfile)
	}

	//import only failed reportOnFailedRecords
	if options.Importonlyfailed {
		importOnlyFailedRecords() //dbsource, options)
		return
	}

	// import only the collections from the options
	if options.Types != "" {
		for _, t := range strings.Split(options.Types, ",") {
			importType(strings.Trim(t, " "))
		}
		return
	}

	// this is the func we are going to be running in daemon mode
	importAll()

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

//************ Not exported functions ************//

// Import gets data from report on failed import, transform it and send it to pillar
func importOnlyFailedRecords() { //dbsource source.Sourcer, limit int, offset int, orderby string, thisStrategy string, reportOnFailedRecords bool) {

	log.User(uuid, "sponge.importOnlyFailedRecords", "### Reading file of data to import.")

	// get the data that needs to be imported
	tables, err := report.ReadReport(options.Reportdbfile) //map[model]map[id]interface{}
	if err != nil {
		log.Error(uuid, "sponge.importOnlyFailedRecords", err, "Getting the rows that will be imported")
	}

	var data []map[string]interface{}

	for table, ids := range tables {

		if len(ids) < 1 { // only one ID
			// Get the data
			log.User(uuid, "sponge.importOnlyFailedRecords", "### Reading data for entity '%s'. \n", table)
			data, err = dbsource.GetData(table, &options) //options.offset, options.limit, options.orderby, options.query)
		} else {
			log.User(uuid, "sponge.importOnlyFailedRecords", "### Reading data for entity '%s', quering '%s'. \n", table, ids)
			data, err = dbsource.GetQueryData(table, &options, ids)
		}
		if err != nil && options.ReportOnFailedRecords {
			report.Record(table, ids, "Failing getting data", err)
		}

		// transform and get data into pillar
		process(table, data)
	}
}

// Import gets ALL data, transform it and send it to pillar
func importAll() {
	//dbsource source.Sourcer, limit int, offset int, orderby string, reportOnFailedRecords bool) {

	log.User(uuid, "sponge.importAll", "### Reading tables to import from strategy file.")

	//Get all the collections's names that we have in the strategy json file
	collections, err := source.GetEntities()

	if err != nil {
		log.Error(uuid, "sponge.importAll", err, "Get collections's names.")
		return
	}

	if dbsource.IsWebService() {
		importFromAPI(collections)
		return
	}

	importFromDB(collections)

}

func importFromAPI(collections []string) {

	pageAfter := ""
	log.User(uuid, "sponge.importFromAPI", "### Reading data from API. \n")

	api, ok := dbsource.(source.API)
	if !ok {
		err := fmt.Errorf("Error asserting sourcer into source.API.")
		log.Error(uuid, "sponge.importFromAPI", err, "Asserting type for source.API")
	}

	var err error
	var data []map[string]interface{}
	var nextPageAfter string

	for true {
		data, nextPageAfter, err = api.GetFireHoseData(pageAfter)
		if err != nil {
			log.Error(uuid, "sponge.importFromAPI", err, "Getting data from API")
			return
		}

		//fmt.Println("DEBUG data ", data)

		if data != nil {
			processAPI(collections, data)
			pageAfter = nextPageAfter
		}

		if data == nil {
			log.User(uuid, "sponge.importFromAPI", "Waiting 5 minutes for more data.")
			time.Sleep(5 * time.Minute) // sleep 5 minutes
		}
	}

}

func importFromDB(collections []string) {
	// var data []map[string]interface{}

	for _, name := range collections { // Reads through all the collections whose transformations are in the strategy configuration file

		foreignEntity := source.GetForeignEntity(name)

		log.User(uuid, "sponge.importAll", "### Reading data to import from %s into collection '%s'. \n", foreignEntity, name)

		// Get the data
		data, err := dbsource.GetData(foreignEntity, &options) //options.offset, options.limit, options.orderby, "")
		if err != nil {
			log.Error(uuid, "sponge.importAll", err, "Get external data for collection %s.", name)
			//RECORD to report about failing modelName
			if options.ReportOnFailedRecords {
				report.Record(name, "", "Failing to get data.", err)
			}
			continue
		}

		log.User(uuid, "sponge.importAll", "### Transforming data and sending it to Coral. \n")
		//transform and send to pillar the data
		process(name, data)
	}
}

// ImportType gets ony data related to table, transform it and send it to pillar
func importType(coralEntity string) { //dbsource source.Sourcer, limit int, offset int, orderby string, query string, modelName string, reportOnFailedRecords bool) {

	foreignEntity := source.GetForeignEntity(coralEntity)
	// Get the data
	log.User(uuid, "sponge.importTable", "### Reading data from table '%s'.", foreignEntity)

	data, err := dbsource.GetData(foreignEntity, &options) //options.offset, options.limit, options.orderby, options.query)
	if err != nil {
		log.Error(uuid, "sponge.importAll", err, "Get external data for table %s.", foreignEntity)
		//RECORD to report about failing modelName
		if options.ReportOnFailedRecords {
			report.Record(foreignEntity, "", "Failing to get data", err)
		}
		return
	}

	// Transform and send to pillar
	process(coralEntity, data)

}

func process(coralName string, data []map[string]interface{}) {

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
		id, newRows, err := fiddler.TransformRow(row, coralName)
		if err != nil {
			log.Error(uuid, "sponge.process", err, "Error when transforming the row %s.", row)
			//RECORD to report about failing transformation
			if options.ReportOnFailedRecords {
				report.Record(coralName, id, "Failing transform data", err)
			}
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
			err = coral.AddRow(newRow, coralName)
			if err != nil {
				log.Error(uuid, "sponge.process", err, "Error when adding a row") // thae row %v to %s.", string(newRow), modelName)
				//RECORD to report about failing adding row to coral db
				if options.ReportOnFailedRecords {
					report.Record(coralName, id, "Failing add row to coral", err)
				}
			}
		}
	}
}

func processAPI(collections []string, data []map[string]interface{}) {

	// Transform the data row by row
	log.User(uuid, "sponge.processAPI", "# Transforming data to the coral schema.\n")
	log.User(uuid, "sponge.processAPI", "# And importing %v documents.", len(data))

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

		for _, name := range collections { // over the same row I look at the different collections in the strategy file
			// transform the row
			id, newRows, err := fiddler.TransformRow(row, name)
			if err != nil {
				log.Error(uuid, "sponge.process", err, "Error when transforming the row %s.", row)
				//RECORD to report about failing transformation
				if options.ReportOnFailedRecords {
					report.Record(name, id, "Failing transform data", err)
				}
				break
			}

			// Usually newRows only will have a document but in the case that we have subcollections
			// we may get more than one document from a transformation
			for _, newRow := range newRows {

				log.Dev(uuid, "sponge.process", "Transforming: %v into %v.", row, newRow)

				// send the row to pillar
				err = coral.AddRow(newRow, name)
				if err != nil {
					log.Error(uuid, "sponge.process", err, "Error when adding a row") // thae row %v to %s.", string(newRow), modelName)

					//RECORD to report about failing adding row to coral db
					if options.ReportOnFailedRecords {
						report.Record(name, id, "Failing add row to coral", err)
					}
				}
			}
		}
	}
}
