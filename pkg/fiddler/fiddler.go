/*
Package fiddler transform, through a strategy file, data from external source into our local coral schema.

*/
package fiddler

import (
	"encoding/json"
	"time"

	"github.com/coralproject/sponge/pkg/log"
	str "github.com/coralproject/sponge/strategy"
)

// global variables related to strategy
var strategy = str.New() // Reads the strategy file

const longForm = "2015-11-02 12:26:05" // date format. To Do: it needs to be defined in the strategy file for the publisher

// Transform from external source data into the coral schema
func Transform(modelName string, data []map[string]interface{}) ([]byte, error) {
	var dataCoral []byte

	table := strategy.GetTables()[modelName]

	// Loop on all the data
	for _, row := range data {

		newRow, err := transformRow(row, table.Fields)
		if err != nil {
			return nil, err
		}

		// add to comments
		dataCoral = append(dataCoral[:], newRow[:]...)
	}

	return dataCoral, nil
}

// Convert a row into the comment coral structure
func transformRow(row map[string]interface{}, fields []map[string]string) ([]byte, error) {
	// "fields": [
	// 	{
	// 		"foreign": "commentid",
	// 		"local": "CommentID",
	// 		"relation": "Identity",
	// 		"type": "int"
	// 	},
	// 	... ]

	newRow := make(map[string]interface{})
	// Loop on the fields for the translation
	for _, f := range fields {
		// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
		newValue := transformField(row[f["foreign"]], f["relation"])
		if newValue != nil {
			newRow[f["local"]] = newValue
		}
	}

	// Convert to Json
	jrow, err := json.Marshal(newRow)
	if err != nil {
		log.Error("transform", "transformCommentrow", err, "Transform Comment")
		return nil, err
	}

	return jrow, nil
}

//Here we transform the record into what we want (based on the configuration in the strategy)
// 1. convert types (values are all strings) into the struct
func transformField(oldValue interface{}, relation string) interface{} {

	var newValue interface{}

	switch relation {
	case "Identity":
		newValue = oldValue
	case "ParseTimeDate":
		newValue, _ = time.Parse(longForm, oldValue.(string))
	}

	return newValue
}
