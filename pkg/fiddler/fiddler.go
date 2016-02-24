// Package fiddler transform, through a strategy file, data from external source into our local coral schema.
//
// It gets a map[string]interface{} as a row and the coral's model that is going to convert it to.
// With that model goes to the strategy file to see what is the transformation that needs to do.
//
package fiddler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ardanlabs/kit/log"
	str "github.com/coralproject/sponge/pkg/strategy"
)

// global variables related to strategy
var (
	strategy   str.Strategy
	dateLayout string
	uuid       string
)

// Init initialize needed variables
func Init(u string) {

	uuid = u

	str.Init(uuid)
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
func TransformRow(row map[string]interface{}, modelName string) (interface{}, []map[string]interface{}, error) { // id row, transformation, error

	var newRows []map[string]interface{}
	var err error

	table := strategy.GetTables()[modelName]
	idField := GetID(modelName)
	id := row[idField]

	if table.Local == "" {
		return "", nil, fmt.Errorf("No table %s found in the strategy file.", table)
	}

	// if has an array field type array
	if strategy.HasArrayField(table) {
		newRows, err = transformRowWithArrayField(modelName, row, table.Fields)
		if err != nil {
			log.Error(uuid, "fiddler.transformRow", err, "Transform the row into several coral documents.")
		}
	} else {
		newRow, err := transformRow(modelName, row, table.Fields)
		if err != nil {
			log.Error(uuid, "fiddler.transformRow", err, "Transform the row into coral.")
			return id, nil, err
		}
		newRows = append(newRows, newRow)
	}

	return id, newRows, err
}

// Convert a row into the comment coral structure
func transformRow(modelName string, row map[string]interface{}, fields []map[string]interface{}) (map[string]interface{}, error) {

	var err error
	// newRow will hold the transformed row
	var newRow map[string]interface{}
	newRow = make(map[string]interface{})

	// source is being used only for the special ocation when the fields relatsionship is source
	var source map[string]interface{}
	source = make(map[string]interface{})

	// Loop on the fields for the transformation
	for _, f := range fields {

		dateLayout = strategy.GetDateTimeFormat(modelName, f["local"].(string))

		// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
		newValue, err := transformField(row[strings.ToLower(f["foreign"].(string))], f["relation"].(string), f["local"].(string))
		if err != nil {
			log.Error(uuid, "fiddler.transformRow", err, "Transforming field %v.", f["foreign"])
		}

		switch f["relation"] {
		case "Source": // { 	"source": { "asset_id": xxx}, }
			source[f["local"].(string)] = newValue

		default: // Identity
			newRow[f["local"].(string)] = newValue
		}

		if source != nil && len(source) > 0 {
			newRow["source"] = source
		}
	}
	return newRow, err
}

// when we are calling this func we are sure that the strategy has a field with an array
func transformRowWithArrayField(modelName string, row map[string]interface{}, fields []map[string]interface{}) ([]map[string]interface{}, error) {
	var err error
	var newRows []map[string]interface{}

	// newRow will hold the transformed row
	var newRow map[string]interface{}
	newRow = make(map[string]interface{})

	// source is being used only for the special ocation when the fields relatsionship is source
	var source map[string]interface{}
	source = make(map[string]interface{})

	// Loop on the fields for the transformation
	for _, f := range fields {
		foreign := f["foreign"].(string)
		relation := f["relation"].(string)

		switch f["relation"] {
		case "Loop":
			newRows = transformArrayFields(foreign, f["fields"], row)
		case "Source": // { 	"source": { "asset_id": xxx}, }
			local := f["local"].(string)
			// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
			newValue, err := transformField(row[strings.ToLower(foreign)], relation, local)
			if err != nil {
				log.Error(uuid, "fiddler.transformRow", err, "Transforming field %s.", f["foreign"])
			}
			source[local] = newValue

		case "Constant":
			local := f["local"].(string)
			newRow[local] = f["value"]

		default: // Identity or ParseTimeDate
			local := f["local"].(string)
			dateLayout = strategy.GetDateTimeFormat(modelName, local)

			// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
			newValue, err := transformField(row[strings.ToLower(foreign)], relation, local)
			if err != nil {
				log.Error(uuid, "fiddler.transformRow", err, "Transforming field %s.", f["foreign"])
			}
			newRow[local] = newValue
		}

		if source != nil && len(source) > 0 {
			newRow["source"] = source
		}
	}
	// add newRow to all rows in newRows
	for i := range newRows {

		for key := range newRow {
			//fmt.Printf("Adding value %v for key %v. \n\n", newRow[key], key)
			if key == "source" {
				for k := range newRow[key].(map[string]interface{}) {
					newRows[i][key].(map[string]interface{})[k] = newRow[key].(map[string]interface{})[k]
				}
			} else {
				newRows[i][key] = newRow[key]
			}
		}
	}

	return newRows, err
}

// "fields": [
// 	{
// 		"foreign": "published",
// 		"local": "date",
// 		"relation": "Identity"
// 	},
// 	{
// 		"foreign": "actor.id",
// 		"local" : "userid",
// 		"relation": "Source"
// 	}
// ],

//
// "object" : {
// 	"likes" : [
// 		{
// 			"actor" : {
// 				...}
// 			}
// 			]
// 		}
//
// row[object.likes.i.published , objects.likes.i.actor.id]
//
// object.likes.0.actor.
// object.likes.1.actor.
// object.likes.2.actor.

func transformArrayFields(foreign string, fields interface{}, row map[string]interface{}) []map[string]interface{} {
	var newRows []map[string]interface{}
	// The transformation for one row with arrays is multiple rows

	// We are getting each row into newRow
	newRow := make(map[string]interface{})
	source := make(map[string]interface{})

	// While still have more rows to add
	finish := false
	i := 0
	for !finish {
		for _, f := range fields.([]interface{}) { // loop through all the fields that we need to create the row

			field := f.(map[string]interface{})
			lastfield := field["foreign"].(string)

			fi := foreign + "." + strconv.Itoa(i) + "." + lastfield // this one is the field in the source document

			// if that row has data on fi
			if row[fi] != nil {
				// transform that specific field
				newvalue, err := transformField(row[fi], field["relation"].(string), field["local"].(string))
				if err != nil {
					log.Error(uuid, "fiddler.transformRow", err, "Transforming field %s.", field["foreign"])
				}
				switch field["relation"] {
				case "Source":
					source[field["local"].(string)] = newvalue
				default:
					newRow[field["local"].(string)] = newvalue
				}
				if source != nil && len(source) > 0 {
					newRow["source"] = source
				}
			} else {
				finish = true //I'm assuming that we are done when one of the fields.i.whatever has not data
				break
			}
		}
		if len(newRow) > 0 { // Add the row only if we got any
			newRows = append(newRows, newRow)
			newRow = make(map[string]interface{}) // initialize the row to get the next one
			source = make(map[string]interface{})
		}

		i++ // NEXT POSSIBLE ROW
	}
	return newRows
}

//Here we transform the record into what we want (based on the configuration in the strategy)
// 1. convert types (values are all strings) into the struct
func transformField(oldValue interface{}, relation string, local string) (interface{}, error) {

	var tfield interface{}
	var err error

	if oldValue != nil {
		switch relation {
		case "Identity":
			return oldValue, err
		case "Source":
			return oldValue, err
		case "ParseTimeDate":
			switch v := oldValue.(type) {
			case string:
				return parseDate(oldValue.(string))
			case time.Time:
				return v.Format(time.RFC3339), nil
			default:
				return "", fmt.Errorf("Type of data %v not recognizable.", v)
			}
		}
	}
	err = fmt.Errorf("Type of transformation %s not found for %v.", relation, oldValue)

	return tfield, err
}

func parseDateLayout(value string) (time.Time, error) {
	var err error
	var dt time.Time

	if value != "" {
		dt, err = time.Parse(dateLayout, value)
	}
	return dt, err
}

func parseDate(value interface{}) (string, error) {

	// on format https://golang.org/pkg/time/#Parse
	// date layout is the representation of 2006 Mon Jan 2 15:04:05 in the desired format. https://golang.org/pkg/time/#pkg-constants
	var err error
	var dt time.Time

	switch v := value.(type) {
	case string:
		dt, err = parseDateLayout(value.(string))
	case time.Time:
		dt = value.(time.Time)
	default:
		err = fmt.Errorf("Type of data %v not recognizable.", v)
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
