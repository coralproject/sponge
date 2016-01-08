package report

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"testing"
)

//func Record(model string, id interface{}, row map[string]interface{}, note string, e error) {
func TestRecord(t *testing.T) {

	model := "comment"
	id := "1"
	row := map[string]interface{}{"id": "1", "body": "comment"}
	note := "This is a note"
	e := errors.New("an error")

	Record(model, id, row, note, e)

	records := GetRecords()

	// records should have one item
	if len(records) != 1 {
		t.Fatalf("got %v records, want 1", len(records))
	}

	if records[0][0] != model {
		t.Fatalf("got %s, it should be %s", records[0][0], model)
	}

	if records[0][1] != id {
		t.Fatalf("got %s, it should be %s", records[0][1], id)
	}

	srow := ""
	for key, value := range row {
		srow = srow + "/" + key + ":" + value.(string)
	}

	if records[0][2] != srow {
		t.Fatalf("got %s, it should be %s", records[0][2], srow)
	}

	if records[0][3] != note {
		t.Fatalf("got %s, it should be %s", records[1][3], note)
	}

	if records[0][4] != fmt.Sprint(e) {
		t.Fatalf("got %s, it should be %s", records[1][3], fmt.Sprint(e))
	}

}

//func Write() {
func TestWrite(t *testing.T) {

	model := "comment"
	id := "1"
	row := map[string]interface{}{"id": "1", "body": "comment"}
	note := "This is a note"
	e := errors.New("an error")

	Record(model, id, row, note, e)

	Write()

	// test the file was written
	filePath := "failed_import.csv"
	outfile, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("unable to open file %s.", filePath)
	}
	defer outfile.Close()

	// test how many rows it wrote
	f := csv.NewReader(outfile)
	f.FieldsPerRecord = 5
	r, err := f.ReadAll()
	if err != nil {
		t.Fatalf("fails at reading the report %s.", filePath)
	}

	if len(r) != 2 { //headers and first row
		t.Fatalf("got %v, it should be 2", len(r))
	}
}

func TestReadReport(t *testing.T) {

	model := "comment"
	id := "1"
	row := map[string]interface{}{"id": "1", "body": "comment"}
	note := "This is a note"
	e := errors.New("an error")

	Record(model, id, row, note, e)

	Write()

	records, err := ReadReport()
	if err != nil {
		t.Fatalf("got an error when reading report")
	}

	if len(records) != 1 {
		t.Fatalf("got %v, it should be 1", len(records))
	}
}
