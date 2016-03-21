package report

import (
	"errors"
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
			return log.DEV
		}
		return ll
	}

	log.Init(os.Stderr, logLevel, log.Ldefault)
}

func teardown() {

}

func TestMain(m *testing.M) {

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

//func Record(model string, id interface{}, row map[string]interface{}, note string, e error) {
func TestRecord(t *testing.T) {

	model := "comments"
	id := "1"
	row := map[string]interface{}{"id": "1", "body": "This is a comment."}
	n := "This is a note."
	e := errors.New("an error")

	u := uuidimported.New()
	dbname := "test.db"

	Init(u, dbname)

	Record(model, id, row, n, e)

	records, err := GetRecords(model)
	//map[1:{/id:1/body:This is a comment. This is a note. an error}]

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

	if r.Row != "/id:1/body:This is a comment." {
		t.Fatalf("got %v, it should be /id:1/body:This is a comment.", r.Row)
	}

	if r.Error != "an error" {
		t.Fatalf("got %v, it should be %v", r.Error, e)
	}

}

func TestReadReport(t *testing.T) {

	model := "comment"
	id := "1"
	row := map[string]interface{}{"id": "1", "body": "comment"}
	note := "This is a note"
	e := errors.New("an error")

	u := uuidimported.New()
	dbname := "test2.db"

	Init(u, dbname)

	Record(model, id, row, note, e)

	records, err := ReadReport(dbname)
	// map[comment:[1]]
	if err != nil {
		t.Fatalf("got an error when reading report")
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
