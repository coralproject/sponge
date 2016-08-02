package report

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
	uuidimported "github.com/pborman/uuid"
)

func setup() {

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.NONE
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)
}

func teardown() {
	// drop the test.db
	dbname := "test.db"
	e := os.Remove(dbname)
	if e != nil {
		fmt.Printf("Error when removing %s: %v\n\n", dbname, e)
	}

	// drop the test2.db
	dbname = "test2.db"
	e = os.Remove(dbname)
	if e != nil {
		fmt.Printf("Error when removing %s: %v\n\n", dbname, e)
	}
}

func TestMain(m *testing.M) {

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

//func Record(model string, id interface{},  note string, e error) {
func TestRecord(t *testing.T) {

	var id string

	model := "comments"
	id = "1"
	n := "This is a note."
	e := errors.New("an error")

	u := uuidimported.New()
	dbname := "test.db"

	Init(u, dbname)

	Record(model, id, n, e)

	records, err := GetRecords(model) //map[1:{1 This is a note. an error}]
	if err != nil {
		t.Fatalf("got error %v, wanted nil", err)
	}

	// records should have one item
	if len(records) != 1 {
		t.Fatalf("got %v records, want 1", len(records))
	}

	r, ok := records[id].(Note)
	if !ok {
		t.Fatalf("%v should be a map", records[id])
	}

	if r.Details != n {
		t.Fatalf("got %v, it should be %v", r.Details, n)
	}

	if r.Error != "an error" {
		t.Fatalf("got %v, it should be %v", r.Error, e)
	}

}

func TestReadReport(t *testing.T) {

	var id string
	model := "comment"
	details := "This is a note"
	e := errors.New("an error")

	id = "1"

	u := uuidimported.New()
	dbname := "test2.db"

	Init(u, dbname)

	Record(model, id, details, e)

	records, err := ReadReport(dbname)
	if err != nil {
		t.Fatalf("got an error when reading report %v.", err)
	}

	// records should have one item
	if len(records) != 1 {
		t.Fatalf("got %v records, want 1", len(records))
	}

	if len(records[model]) != 1 {
		t.Fatalf("got %v records, want 1", len(records[model]))
	}

	if records[model][0] != "1" {
		t.Fatalf("got %v records, want 1", records[model][0])
	}

}
