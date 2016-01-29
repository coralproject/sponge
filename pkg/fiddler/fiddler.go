// Package fiddler transform, through a strategy file, data from external source into our local coral schema.
//
// It gets a map[string]interface{} as a row and the coral's model that is going to convert it to.
// With that model goes to the strategy file to see what is the transformation that needs to do.
//
package fiddler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/pkg/strategy"
)

// global variables related to strategy
var (
	strategy   str.Strategy
	dateLayout string
)

// Init initialize needed variables
func Init() {

	str.Init()
	strategy = str.New() // Reads the strategy file
}

// GetID returns the identifier for modelName
func GetID(modelName string) string {
	return strategy.GetIDField(modelName)
}

// GetCollections give the names of all the collections in the strategy file
func GetCollections() []string {
	tables := strategy.GetTables() // map[string]Table
	keys := []string{}
	for k := range tables {
		keys = append(keys, k)
	}
	return keys
}

// TransformRow transform a row of data into the coral schema
func TransformRow(row map[string]interface{}, modelName string) ([]byte, error) {

	table := strategy.GetTables()[modelName]

	if table.Local == "" {
		return nil, errors.New("No table found in the strategy file.")
	}

	newRow, err := transformRow(modelName, row, table.Fields)
	if err != nil {
		log.Error("transform", "TransformRow", err, "Transform Data")
		return nil, err
	}
	// Convert to Json
	dataCoral, err := json.Marshal(newRow)
	if err != nil {
		log.Error("transform", "TransformRow", err, "Transform Data")
		return nil, err
	}

	return dataCoral, err
}

// Convert a row into the comment coral structure

func transformRow(modelName string, row map[string]interface{}, fields []map[string]string) (map[string]interface{}, error) {

	var err error
	// newRow will hold the transformed row
	var newRow map[string]interface{}
	newRow = make(map[string]interface{})

	// source is being used only for the special ocation when the fields relatsionship is source
	var source map[string]interface{}
	source = make(map[string]interface{})

	// Loop on the fields for the transformation
	for _, f := range fields {

		dateLayout = strategy.GetDateTimeFormat(modelName, f["local"])

		// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
		newValue, err := transformField(row[strings.ToLower(f["foreign"])], f["relation"], f["local"])
		if err != nil {
			log.Error("fiddler", "transformRow", err, "Transforming field %s.", f["foreign"])
		}

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

	return newRow, err
}

//Here we transform the record into what we want (based on the configuration in the strategy)
// 1. convert types (values are all strings) into the struct
func transformField(oldValue interface{}, relation string, local string) (interface{}, error) {

	var err error

	if oldValue != nil {
		switch relation {
		case "Identity":
			return oldValue, nil
		case "Source":
			return oldValue, nil
		case "ParseTimeDate":
			return parseDate(oldValue.(string))
		}
		err = fmt.Errorf("Type of transformation %s not found for %v.", relation, oldValue)
		// } else {
		// 	err = fmt.Errorf("Empty value for %s field.", local)
	}

	return nil, err
}

func parseDate(value string) (string, error) {

	// on format https://golang.org/pkg/time/#Parse
	// date layout is the representation of 2006 Mon Jan 2 15:04:05 in the desired format. https://golang.org/pkg/time/#pkg-constants

	if value == "" {
		return "", nil
	}

	dt, err := time.Parse(dateLayout, value)
	if err != nil {
		log.Error("fiddler", "parseDate", err, "Parsing date %s.", value)
	}

	dtRFC3339 := dt.Format(time.RFC3339)

	return dtRFC3339, err
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
