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

	var err error

	uuid = u

	str.Init(uuid)
	strategy, err = str.New() // Reads the strategy file
	if err != nil {
		log.Error(uuid, "fiddler.init", err, "Reading the streategy file.")
	}
}

// GetID returns the identifier for modelName
func GetID(modelName string) string {
	return strategy.GetIDField(modelName)
}

// GetCollections give the names of all the collections in the strategy file
func GetCollections() []string {
	tables := strategy.GetEntities() // map[string]Table
	keys := []string{}
	for k := range tables {
		keys = append(keys, k)
	}
	return keys
}

// TransformRow transform a row of data into the coral schema
func TransformRow(row map[string]interface{}, coralName string) (interface{}, []map[string]interface{}, error) { // id row, transformation, error

	var newRows []map[string]interface{}
	var err error

	table := strategy.GetEntities()[coralName]
	idField := GetID(coralName)
	id := row[idField]

	if table.Local == "" {
		return "", nil, fmt.Errorf("No table %v found in the strategy file.", table)
	}

	// if has an array field type array
	if strategy.HasArrayField(table) {
		newRows, err = transformRowWithArrayField(coralName, row, table.Fields)
		if err != nil {
			log.Error(uuid, "fiddler.transformRow", err, "Transform the row into several coral documents.")
		}
	} else {
		newRow, err := transformRow(coralName, row, table.Fields)
		if err != nil {
			log.Error(uuid, "fiddler.transformRow", err, "Transform the row into coral.")
			return id, nil, err
		}
		newRows = append(newRows, newRow)
	}

	return id, newRows, err
}

// Convert a row into the comment coral structure
func transformRow(coralName string, row map[string]interface{}, fields []map[string]interface{}) (map[string]interface{}, error) {

	var err error
	// newRow will hold the transformed row
	var newRow map[string]interface{}
	newRow = make(map[string]interface{})

	// source is being used only for the special ocation when the fields relatsionship is source
	var source map[string]interface{}
	source = make(map[string]interface{})

	// metadata is being used to send metadata
	var metadata map[string]interface{}
	metadata = make(map[string]interface{})

	// Loop on the fields for the transformation
	for _, f := range fields {

		dateLayout = strategy.GetDateTimeFormat(coralName, f["local"].(string))

		foreign, ok := f["foreign"].(string)
		if !ok {
			return nil, fmt.Errorf("%v not expected type", f["foreign"])
		}
		foreign = strings.ToLower(foreign)

		relation, ok := f["relation"].(string)
		if !ok {
			return nil, fmt.Errorf("%v not expected type", f["relation"])
		}
		relation = strings.ToLower(relation)
		local, ok := f["local"].(string)
		if !ok {
			return nil, fmt.Errorf("%v not expected type", f["local"])
		}
		local = strings.ToLower(local)

		// If this is a required field but the value we got is null then skip this document
		required, ok := f["required"]
		if ok && required == "true" && row[foreign] == nil {
			return nil, fmt.Errorf("Required field %s is null.", foreign)
		}

		// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
		newValue, err := transformField(row[foreign], relation, local, coralName)
		if err != nil {
			log.Error(uuid, "fiddler.transformRow", err, "Transforming field %v, value %v.", foreign, row[foreign])
			return nil, err
		}

		// We are not adding fields which have an empty value
		if newValue == nil {
			break
		}

		switch f["relation"] {
		case "Source": // { 	"source": { "asset_id": xxx}, }
			source[local] = newValue

		case "Metadata": // { "metadata": { "source": xxx , "markers": [xxx]} }
			metadata[local] = newValue

		default: // Identity or SubDocument or Status or Constant
			newRow[local] = newValue
		}

		if source != nil && len(source) > 0 {
			newRow["source"] = source
		}

		if metadata != nil && len(metadata) > 0 {
			newRow["metadata"] = metadata
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

	var metadata map[string]interface{}
	metadata = make(map[string]interface{})

	// Loop on the fields for the transformation
	for _, f := range fields {
		foreign, ok := f["foreign"].(string)
		if !ok {
			return nil, fmt.Errorf("%v not expected type", f["foreign"])
		}
		foreign = strings.ToLower(foreign)
		relation, ok := f["relation"].(string)
		if !ok {
			return nil, fmt.Errorf("%v not expected type", f["relation"])
		}
		relation = strings.ToLower(relation)

		switch relation {
		case "loop":
			newRows, err = transformArrayFields(foreign, f["fields"], row, modelName)
			if err != nil {
				log.Error(uuid, "fiddler.transformRow", err, "Transforming field %s.", f["foreign"])
			}
		case "source": // { 	"source": { "asset_id": xxx}, }
			local := f["local"].(string)
			local = strings.ToLower(local)
			// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
			newValue, err := transformField(row[foreign], relation, local, modelName)
			if err != nil {
				log.Error(uuid, "fiddler.transformRow", err, "Transforming field %s.", f["foreign"])
			}
			source[local] = newValue

		case "metadata": // { 	"metadata": { "source": xxx},  "markers": [ xxx ]}
			local := f["local"].(string)
			local = strings.ToLower(local)
			// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
			newValue, err := transformField(row[foreign], relation, local, modelName)
			if err != nil {
				log.Error(uuid, "fiddler.transformRow", err, "Transforming field %s.", f["foreign"])
			}
			metadata[local] = newValue

		case "constant":
			local, ok := f["local"].(string)
			if !ok {
				return nil, fmt.Errorf("%v not expected type", f["local"])
			}

			newRow[local] = f["value"]

		default: // Identity or ParseTimeDate
			local, ok := f["local"].(string)
			if !ok {
				return nil, fmt.Errorf("%v not expected type", f["local"])
			}
			dateLayout = strategy.GetDateTimeFormat(modelName, local)

			// convert field f["foreign"] with value row[f["foreign"]] into field f["local"], whose relationship is f["relation"]
			newValue, err := transformField(row[foreign], relation, local, modelName)
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

			if key == "source" || key == "metadata" {
				nr, ok := newRow[key].(map[string]interface{})
				if !ok {
					return nil, fmt.Errorf("%v not expected type", newRow[key])
				}

				for k := range nr {
					nrsi, ok := newRows[i][key].(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("%v not expected type", newRows[i][key])
					}
					nrsi[k] = nr[k]
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

func transformArrayFields(foreign string, fields interface{}, row map[string]interface{}, modelName string) ([]map[string]interface{}, error) {
	var newRows []map[string]interface{}
	// The transformation for one row with arrays is multiple rows

	// We are getting each row into newRow
	newRow := make(map[string]interface{})
	source := make(map[string]interface{})
	metadata := make(map[string]interface{})

	// While still have more rows to add
	finish := false
	i := 0
	for !finish {
		fis, ok := fields.([]interface{})
		if !ok {
			log.Error(uuid, "fiddler.transformArrayFields", fmt.Errorf("%v not expected type", fields), "Not expected interface{} type")
		}

		for _, f := range fis { // loop through all the fields that we need to create the row

			field, ok := f.(map[string]interface{})
			if !ok {
				log.Error(uuid, "fiddler.transformArrayFields", fmt.Errorf("%v not expected type", f), "Not expected type")
			}
			lastfield, ok := field["foreign"].(string)
			if !ok {
				log.Error(uuid, "fiddler.transformArrayFields", fmt.Errorf("%v not expected type", field["foreign"]), "Not expected type")
			}
			lastfield = strings.ToLower(lastfield)

			fi := foreign + "." + strconv.Itoa(i) + "." + lastfield // this one is the field in the source document

			relation, ok := field["relation"].(string)
			if !ok {
				log.Error(uuid, "fiddler.transformArrayFields", fmt.Errorf("%v not expected type", field["relation"]), "Not expected type")
			}
			relation = strings.ToLower(relation)
			local, ok := field["local"].(string)
			if !ok {
				log.Error(uuid, "fiddler.transformArrayFields", fmt.Errorf("%v not expected type", field["local"]), "Not expected type")
			}
			local = strings.ToLower(local)

			// if that row has data on fi
			if row[fi] != nil {
				// transform that specific field
				newvalue, err := transformField(row[fi], relation, local, modelName)
				if err != nil {
					log.Error(uuid, "fiddler.transformRow", err, "Transforming field %s.", field["foreign"])
				}
				switch field["relation"] {
				case "Source":
					source[local] = newvalue
				case "Metadata":
					metadata[local] = newvalue
				default:
					newRow[local] = newvalue
				}
				if source != nil && len(source) > 0 {
					newRow["source"] = source
				}
				if metadata != nil && len(metadata) > 0 {
					newRow["metadata"] = metadata
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
			metadata = make(map[string]interface{})
		}

		i++ // NEXT POSSIBLE ROW
	}
	return newRows, nil
}

//Here we transform the record into what we want (based on the configuration in the strategy)
// 1. convert types (values are all strings) into the struct
//func transformField(oldValue interface{}, relation string, local string, coralName string) (interface{}, error) {
// // 1. convert types into the struct
func transformField(oldValue interface{}, relation string, coralName string, foreignfield string) (interface{}, error) {

	var newValue interface{}
	var err error

	if oldValue != nil {
		switch strings.ToLower(relation) {
		case "identity":
			return oldValue, err
		case "source":
			return oldValue, err
		case "metadata":
			return oldValue, err
		case "status":
			ov, ok := oldValue.(string)
			if !ok {
				return nil, fmt.Errorf("%v not expected type string", oldValue)
			}
			return strategy.GetStatus(coralName, ov), err
		case "subdocument":
			//oldvalue is [map[_id:terrence-mccoy name:Terrence McCoy url:http://www.washingtonpost.com/people/terrence-mccoy twitter:@terrence_mccoy]]
			// check all the other fields in the subcollection
			// oldValue is an array of values []map[string]interface{}
			// look through all the fields and do the transformation one by one
			fields := strategy.GetFieldsForSubDocument(coralName, foreignfield)

			// oldValue is interface{} and it is []map[string]interface{} when subdocument
			ovSlice := oldValue.([]interface{}) // cast into []interface{}
			ov := make([]map[string]interface{}, len(ovSlice))
			for i := range ov {
				ov[i] = ovSlice[i].(map[string]interface{})
			}

			return transformSubDocumentField(ov, relation, fields)
		case "parsetimedate":
			switch v := oldValue.(type) {
			case string:
				return parseDate(v)
			case time.Time:
				return v.Format(time.RFC3339), nil
			case float64:
				return time.Unix(int64(v), 0).Format(time.RFC3339), nil
			default:
				return "", fmt.Errorf("Type of data %v not recognizable.", v)
			}
		}
		err = fmt.Errorf("Type of transformation %s not found for %v.", relation, oldValue)
	}

	return newValue, err
}

func parseDateLayout(value string) (time.Time, error) {
	var err error
	var dt time.Time

	if value != "" {
		dt, err = time.Parse(dateLayout, value)
	}
	return dt, err
}

func transformSubDocumentField(oldValue []map[string]interface{}, relation string, fields []map[string]interface{}) (interface{}, error) {

	var err error

	newValue := make([]map[string]interface{}, len(oldValue))
	//stringType := reflect.TypeOf("string").Elem()

	for i := range oldValue { // oldValue is an array of documents
		r := oldValue[i] // when it is a subdocument we know that oldValue is really a slice of maps
		newValue[i] = make(map[string]interface{})
		for _, f := range fields { // each field is a map with local, foreign fields

			if r[f["foreign"].(string)] != nil {
				newValue[i][f["local"].(string)] = r[f["foreign"].(string)]
			}
		}
	}

	return newValue, err
}
func parseDate(value interface{}) (string, error) {

	// on format https://golang.org/pkg/time/#Parse
	// date layout is the representation of 2006 Mon Jan 2 15:04:05 in the desired format. https://golang.org/pkg/time/#pkg-constants
	var err error
	var dt time.Time

	switch v := value.(type) {
	case string:
		vs, ok := value.(string)
		if !ok {
			return "", fmt.Errorf("%v not expected type", vs)
		}
		dt, err = parseDateLayout(vs)
	case time.Time:
		dt, ok := value.(time.Time)
		if !ok {
			return "", fmt.Errorf("%v not expected type", dt)
		}
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
