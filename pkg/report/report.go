/*
Package report

CSV errors file with this fields:
table, id, row, "what went wrong"

*/

package report

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/coralproject/sponge/pkg/log"
)

const (
	filePath = "failed_import.csv"
)

var records [][]string

// Init initialize the records to be recorded
func Init() {
	records = [][]string{
		{"table", "id", "row", "note", "error"},
	}

}

// Record adds a new record to the report
func Record(model string, id interface{}, row map[string]interface{}, note string, e error) {

	// convertthe row to string
	srow := ""
	for key, value := range row {
		srow = srow + "/" + key + ":" + value.(string)
	}

	n := len(records)

	var original [][]string
	original = make([][]string, n)
	copy(original, records)

	records = make([][]string, n+1)
	copy(records[:n], original)
	records[n] = []string{model, id.(string), srow, note, fmt.Sprint(e)}

}

// Write writes the report to disk
func Write() {
	if len(records) > 1 {
		outfile, err := os.Create(filePath)
		if err != nil {
			log.Fatal("report", "Write", "Unable to open output")
		}
		defer outfile.Close()

		w := csv.NewWriter(outfile)

		for _, record := range records {
			if err := w.Write(record); err != nil {
				log.Error("report", "Write", err, "Writing to CSV file")
			}
		}

		// Write any buffered data to the underlying writer (standard output).
		w.Flush()

		if err := w.Error(); err != nil {
			log.Error("report", "Write", err, "Writing to CSV file")
			fmt.Println(records)
		}
	} else {
		log.User("report", "Write", "No fails attempts.")
	}
}

// Read gets the data that needs to be imported from already data
func Read() ([][]string, error) {
	return records, nil
}
