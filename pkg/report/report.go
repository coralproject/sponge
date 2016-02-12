/*
Package report

CSV errors file with this fields:
table, id, row, "what went wrong"

*/

package report

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/ardanlabs/kit/log"
)

const (
	filePathDefault = "failed_import.csv"
)

var (
	records  [][]string
	filePath string
	uuid     string
)

// Init initialize the records to be recorded
func Init(u string, errorsfile string) {

	uuid = u

	// Initialize with the headers
	records = [][]string{
		{"table", "id", "row", "note", "error"},
	}

	if errorsfile == "" {
		filePath = filePathDefault
	} else {
		filePath = errorsfile
	}
}

// Record adds a new record to the report
func Record(model string, id interface{}, row map[string]interface{}, note string, e error) {

	// convertthe row to string
	// srow := ""
	// for key, value := range row {
	// 	srow = srow + "/" + key + ":" + value
	// }

	srow, err := json.Marshal(row)
	if err != nil {
		log.Error(uuid, "report.record", e, "Marshalling the report's row.")
	}

	n := len(records)

	var original [][]string
	original = make([][]string, n)
	copy(original, records)

	records = make([][]string, n+1)
	copy(records[:n], original)
	records[n] = []string{model, fmt.Sprintf("%v", id), fmt.Sprintf("%v", string(srow)), note, fmt.Sprint(e)}

	Write()

}

// GetRecords returns the actual records that I'm recording
func GetRecords() [][]string {
	return records[1:len(records)]
}

// Write writes the report to disk
func Write() {
	// remove existing file
	//os.Remove(filePath)

	// only write the file if there is any report to write
	if len(records) > 1 {
		outfile, err := os.Create(filePath)
		if err != nil {
			log.Fatal(uuid, "report.write", "Creating or opening the file %s.", filePath)
		}
		defer outfile.Close()

		w := csv.NewWriter(outfile)

		for _, record := range records {
			if err := w.Write(record); err != nil {
				log.Error(uuid, "report.write", err, "Writing to the CSV file %s.", filePath)
			}
		}

		// Write any buffered data to the underlying writer (standard output).
		w.Flush()

		if err := w.Error(); err != nil {
			log.Error(uuid, "report.write", err, "Writing to CSV file")
		}
	} else {
		log.User(uuid, "report.write", "There is nothing to write into the report.")
	}
}

//* This functions are for the old report that is already save in disk. *//

// ReadReport gets the data that needs to be imported from the report already in disk and return the records in a way that can be easily read it
func ReadReport(filePath string) ([]map[string]interface{}, error) {
	// Read the CSV file
	outfile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(uuid, "report.read", "Unable to open file %s.", filePath)
	}
	defer outfile.Close()

	// Read the file
	f := csv.NewReader(outfile)
	f.FieldsPerRecord = 5
	r, err := f.ReadAll()
	if err != nil {
		log.Error(uuid, "report.read", err, "Fails at reading the report %s.", filePath)
	}

	// Get into results
	// [
	// 	{
	//		table: xxx,
	//		ids: [x,x,xx,x,xxx]
	//	},
	// 	{
	// 		table: x,
	// 		ids: []
	//  }
	// ]
	var results []map[string]interface{}
	for _, row := range r {
		if row[0] == "table" {
			continue
		}
		table := row[0]
		id := row[1]

		results = addRecord(results, table, id)
	}

	return results, nil
}

//results
//[{
// "table": "comment"
// "ids": []
// }]
// adds the table, id to the results
func addRecord(results []map[string]interface{}, table string, id string) []map[string]interface{} {

	n := len(results)

	// copy results to a temporal map
	var original []map[string]interface{}
	original = make([]map[string]interface{}, n)
	copy(original, results)

	if id == "" { // add the whole table to the results

		// increment our map
		results = make([]map[string]interface{}, n+1)
		copy(results[:n], original)

		// item to be added
		item := map[string]interface{}{"table": table, "ids": []string{}}

		results[n] = item
		return results
	}

	// search table in the original
	pos := sort.Search(n, func(i int) bool {
		return (original[i]["table"] == table)
	}) // returns the position where the value is

	if pos == n { // not found
		// increment our map
		results = make([]map[string]interface{}, n+1)
		copy(results[:n], original)

		item := map[string]interface{}{"table": table, "ids": []string{id}}
		results[n] = item
	} else {
		ids := original[pos]["ids"].([]string)
		results[pos]["ids"] = append(ids, id)
	}

	return results
}
