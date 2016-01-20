/*
Package fiddler transform, through a strategy file, data from external source into our local coral schema.

*/
package fiddler

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/strategy"
)

// global variables related to strategy
var strategy str.Strategy

func init() {

	strategy = str.New() // Reads the strategy file

}

// GetID returns the identifier for modelName
func GetID(modelName string) string {
	return strategy.GetIDField(modelName)
}

// TransformRow transform a row of data into the coral schema
func TransformRow(row map[string]interface{}, modelName string) ([]byte, error) {

	table := strategy.GetTables()[modelName]

	if table.Local == "" {
		return nil, errors.New("No table found in the strategy file.")
	}

	newRow, err := transformRow(row, table.Fields)

	// Convert to Json
	dataCoral, err := json.Marshal(newRow)
	if err != nil {
		log.Error("transform", "TransformRow", err, "Transform Data")
		return nil, err
	}

	return dataCoral, err
}

// Convert a row into the comment coral structure
func transformRow(row map[string]interface{}, fields []map[string]string) (map[string]interface{}, error) { //([]byte, error) {
	// "fields": [
	// 	{
	// 		"foreign": "commentid",
	// 		"local": "CommentID",
	// 		"relation": "Identity",
	// 		"type": "int"
	// 	},
	// 	... ]

	// newRow will hold the transformed row
	var newRow map[string]interface{}
	newRow = make(map[string]interface{})

	// source is being used only for the special ocation when the fields relatsionship is source
	var source map[string]interface{}
	source = make(map[string]interface{})

	// Loop on the fields for the transformation
	for _, f := range fields {

		// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
		newValue := transformField(row[strings.ToLower(f["foreign"])], f["relation"], f["local"])

		if newValue != nil {

			if f["relation"] != "Source" {
				newRow[f["local"]] = newValue // newvalue could be string or time.Time or int
			} else { // special case when I'm looking into a source relationship
				// {
				//	"source":
				//					{ "asset_id": xxx},
				// }
				source[f["local"]] = newValue
			}
		}
	}

	if source != nil && len(source) > 0 {
		newRow["source"] = source
	}

	return newRow, nil
}

//Here we transform the record into what we want (based on the configuration in the strategy)
// 1. convert types (values are all strings) into the struct
func transformField(oldValue interface{}, relation string, local string) interface{} {

	if oldValue != nil {
		switch relation {
		case "Identity":
			return oldValue
		case "Source":
			return oldValue
		case "ParseTimeDate":
			var newValue time.Time
			// TODO: move time format into strategy file
			newValue, _ = time.Parse("2006-01-02", oldValue.(string))

			return newValue
		}
	}
	return nil
}

// source is [ { "asset_id": xxx}, { "comment_id": xxx} ]
// newItem is the format { "asset_id": xxx }
func appendField(source []map[string]interface{}, item interface{}) []map[string]interface{} {
	n := len(source)
	total := len(source) + 1
	if total > cap(source) {
		newSize := total*3/2 + 1
		newSource := make([]map[string]interface{}, total, newSize)
		copy(newSource, source)
		source = newSource
	}

	source = source[:total]
	copy(source[n:], item.([]map[string]interface{}))
	source = source[:total]

	return source
}
