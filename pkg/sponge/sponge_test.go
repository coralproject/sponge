package sponge

import "testing"

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
