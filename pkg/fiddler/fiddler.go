/*
Package fiddler transform, through a strategy file, data from external source into our local coral schema.

*/
package fiddler

import (
	"encoding/json"
	"time"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/strategy"
)

// global variables related to strategy
var strategy = str.New() // Reads the strategy file

const longForm = "2015-11-02 12:26:05" // date format. To Do: it needs to be defined in the strategy file for the publisher

// Transform from external source data into the coral schema
func Transform(modelName string, data []map[string]interface{}) ([]byte, error) {

	var d []map[string]interface{}

	table := strategy.GetTables()[modelName]

	// Loop on all the data
	for _, row := range data {

		newRow, err := transformRow(row, table.Fields)
		if err != nil {
			return nil, err
		}
		// append a row to the stream
		d = appendRow(d, newRow)
	}

	// Convert to Json
	dataCoral, err := json.Marshal(d)
	if err != nil {
		log.Error("transform", "Transform", err, "Transform Data")
		return nil, err
	}

	return dataCoral, nil
}

// Convert a row into the comment coral structure
func transformRow(row map[string]interface{}, fields []map[string]string) ([]map[string]interface{}, error) { //([]byte, error) {
	// "fields": [
	// 	{
	// 		"foreign": "commentid",
	// 		"local": "CommentID",
	// 		"relation": "Identity",
	// 		"type": "int"
	// 	},
	// 	... ]

	var newRow []map[string]interface{}
	var source []map[string]interface{}
	newRow = make([]map[string]interface{}, 1)
	newRow[0] = make(map[string]interface{})
	// Loop on the fields for the translation
	for _, f := range fields {
		// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
		newValue := transformField(row[f["foreign"]], f["relation"], f["local"])
		if newValue != nil {

			if f["relation"] != "Source" {
				newRow[0][f["local"]] = newValue // newvalue could be string or time.Time or int
			} else { // special case when I'm looking into a source relationship
				// {
				//	"source":
				//				[
				//					{ "asset_id": xxx},
				//				]
				// }
				// append a field to the slice source, newValue's example: { "asset_id": xxx }
				source = appendField(source, newValue)
			}
		}
	}
	newRow[0]["source"] = source

	return newRow, nil
}

//Here we transform the record into what we want (based on the configuration in the strategy)
// 1. convert types (values are all strings) into the struct
func transformField(oldValue interface{}, relation string, local string) interface{} {

	switch relation {
	case "Identity":
		return oldValue
	case "Source": // this is dirty! look at this again please
		var newValue []map[string]interface{}
		newValue = make([]map[string]interface{}, 1)
		newValue[0] = make(map[string]interface{})
		newValue[0][local] = oldValue
		return newValue
	case "ParseTimeDate":
		var newValue time.Time
		newValue, _ = time.Parse(longForm, oldValue.(string))
		return newValue
	}

	return nil
}

// appends an item to []item
func appendRow(items []map[string]interface{}, item []map[string]interface{}) []map[string]interface{} {
	n := len(items)
	total := len(items) + 1
	if total > cap(items) {
		newSize := total*3/2 + 1
		newItems := make([]map[string]interface{}, total, newSize)
		copy(newItems, items)
		items = newItems
	}

	items = items[:total]
	copy(items[n:], item)
	items = items[:total]

	return items
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
