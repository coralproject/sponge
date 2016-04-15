// Package coral is a package to interact with pillar endpoints
//
// Coral System
//
// Pillar is the service to get data into the coral system.  Each endpoint is documented in https://github.com/coralproject/pillar/tree/master/server
// In this package we are using:
//
// Assets
//
// The endpoint is setup in the translations file. It receives one document per POST.
// The structure of the json we are sending is at the package models in https://github.com/coralproject/pillar/
//
// Users
//
// The endpoint is setup in the translations file. It receives one document per POST.
// The structure of the json we are sending is at the package models in https://github.com/coralproject/pillar/
//
// Comments
//
// The endpoint is setup in the translations file. It receives one document per POST.
// The structure of the json we are sending is at the package models in https://github.com/coralproject/pillar/
//
// CreateIndex
//
// The endpoint is setup in the translations file. It receives one document per POST.
// In the strategy file, for each collection, we are getting
// "Index": [{
// 	"name": "asset-url",
// 	"key": "asseturl",
// 	"unique": true,
// 	"dropdups": true
// }],
// And we want to send to Pillar
// 	 {
// 			"target": "asset",
// 			"index": {
// 					"name": "asset-url",
// 					"key": ["url"],
// 					"unique": true,
// 					"dropdups": true
// 			}
// 	 },
//
package coral

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/sponge/pkg/strategy"
	"github.com/coralproject/sponge/pkg/webservice"
)

const (
	retryTimes int    = 3
	methodGet  string = "GET"
	methodPost string = "POST"
)

var (
	endpoints map[string]string // model -> endpoint
)

var (
	uuid string
	str  strategy.Strategy
)

// Init initialization of logs and strategy
func Init(u string) {

	var err error
	uuid = u

	strategy.Init(uuid)
	str, err = strategy.New()
	if err != nil {
		log.Error(uuid, "coral.init", err, "Reading the streategy file.")
	}
	endpoints = str.GetPillarEndpoints()
}

// AddRow send the row to pillar based on which collection is
func AddRow(data map[string]interface{}, tableName string) error {

	var err error
	if _, ok := endpoints[tableName]; ok {

		d, err := json.Marshal(data)
		if err != nil {
			log.Error(uuid, "coral.Addrow", err, "Marshalling %v.", data)
		}

		userAgent := fmt.Sprintf("Sponge Publisher %s.", str.Name)
		_, err = webservice.DoRequest(uuid, userAgent, methodPost, endpoints[tableName], bytes.NewBuffer(d))
		if err != nil {
			log.Error(uuid, "coral.Addrow", err, "Sending request to PILLAR with %v.", data)
		}
	} else {
		err = fmt.Errorf("No information about %s in the available endpoints.", tableName)
	}

	return err
}

// CreateIndex calls the service to create index
func CreateIndex(collection string) error {

	var err error

	// get index
	is := str.GetIndexBy(collection) // []map[string]interface{}

	if is == nil {
		err = fmt.Errorf("%s does not exist", collection)
	}

	// get Endpoint
	createIndexURL := str.GetPillarEndpoints()["index"]

	indexes := make([]map[string]interface{}, len(is))
	for i := range is {
		indexes[i] = make(map[string]interface{})
		indexes[i]["target"] = collection
		indexes[i]["index"] = is[i]

		var data []byte
		data, err = json.Marshal(indexes[i])
		if err != nil {
			log.Error(uuid, "coral.createIndex", err, "Marshal index information.")
		}

		userAgent := fmt.Sprintf("Sponge Publisher %s.", str.Name)
		_, err = webservice.DoRequest(uuid, userAgent, methodPost, createIndexURL, bytes.NewBuffer(data))
		if err != nil {
			log.Error(uuid, "coral.createIndex", err, "Sending request to create Index to Pillar.")
		}

	}

	return err
}
