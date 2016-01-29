package sponge

import (
	"fmt"
	"os"
	"testing"

	"github.com/ardanlabs/kit/cfg"
	"github.com/ardanlabs/kit/log"
)

func setup() {

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.DEV
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	// MOCK STRATEGY CONF
	strategyConf := "../../tests/strategy_test.json"
	e := os.Setenv("STRATEGY_CONF", strategyConf) // IS NOT REALLY SETTING UP THE VARIABLE environment FOR THE WHOLE PROGRAM :(
	if e != nil {
		fmt.Println("It could not setup the mock strategy conf variable")
	}

	// Initialize coral
	Init()
}

func teardown() {
}

func TestMain(m *testing.M) {

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

// Signature: process(modelName string, data []map[string]interface{})
func TestProcess(t *testing.T) {

	modelName := "comment"
	var data []map[string]interface{}

	// mock up pillar

	process(modelName, data)

	// check data is sent to pillar with the right transformations

}

// func TestImportAll(t *testing.T) {
//
// }

//
// func TestImportFailedRecordsWholeTable(t *testing.T) {
//
// }
//
// func TestImportFailedRecordsOneRecord(t *testing.T) {
//
// }
//
// func TestImportFailedRecordsTwoRecords(t *testing.T) {
//
// }
//
// func TestImportFailedRecordsTwoRecordsSeveralTables(t *testing.T) {
//
// }
//
// func TestProcess(t *testing.T) {
//
// }
//
// func ExampleProcess() {
//
// }
