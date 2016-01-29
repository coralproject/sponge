// Package coral is a package to interact with pillar endpoints
//
// Coral System
//
// Pillar is the service to get data into the coral system.  Each endpoint is documented in https://github.com/coralproject/pillar/tree/master/server
// In this package we are using:
//
// Assets
//
// The endpoint is setup in the ASSET_URL environment variable. It receives one document per POST.
// The structure of the json we are sending is at https://github.com/coralproject/pillar/blob/master/server/model/model.go
//
// Users
//
// The endpoint is setup in the USER_URL environment variable. It receives one document per POST.
// The structure of the json we are sending is at https://github.com/coralproject/pillar/blob/master/server/model/model.go
//
// Comments
//
// The endpoint is setup in the COMMENT_URL environment variable. It receives one document per POST.
// The structure of the json we are sending is at https://github.com/coralproject/pillar/blob/master/server/model/model.go
//
// CreateIndex
//
// The endpoint is setup in the CREATE_INDEX_URL environment variable. It receives information about what indexes to create.
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
//
package coral

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/sponge/pkg/strategy"
)

const (
	retryTimes int    = 3
	methodGet  string = "GET"
	methodPost string = "POST"
)

type restResponse struct {
	status  string
	header  http.Header
	payload string
}

var (
	endpoints map[string]string // model -> endpoint
)

// Init initialization of logs and strategy
func Init() {

	strategy.Init()
	endpoints = strategy.New().GetPillarEndpoints()
}

// AddRow send the row to pillar based on which collection is
func AddRow(data []byte, tableName string) error {

	var err error
	if _, ok := endpoints[tableName]; ok {

		err = doRequest(methodPost, endpoints[tableName], bytes.NewBuffer(data))
		if err != nil {
			log.Error("coral", "Addrow", err, "Sending request to PILLAR with %v.", bytes.NewBuffer(data))
		}
	} else {
		err = fmt.Errorf("No %s in the endpoints.", tableName)
	}

	return err
}

// CreateIndex calls the service to create index
func CreateIndex(collection string) error {

	var err error

	// get index
	s := strategy.New()
	is := s.GetIndexBy(collection) // []map[string]interface{}

	// get Endpoint
	createIndexURL := s.GetPillarEndpoints()["index"]

	indexes := make([]map[string]interface{}, len(is))
	for i := range is {
		indexes[i] = make(map[string]interface{})
		indexes[i]["target"] = collection

		// indexes[i]["index"] = map[string]interface{}{
		// 	"name":     is[i]["name"].(string),
		// 	"key":      is[i]["keys"],
		// 	"unique":   is[i]["unique"].(string),
		// 	"dropdups": is[i]["dropdups"].(string),
		// }

		indexes[i]["index"] = is[i]

		var data []byte
		data, err = json.Marshal(indexes[i])
		if err != nil {
			log.Error("coral", "CreateIndex", err, "Creating Index.")
		}

		err = doRequest(methodPost, createIndexURL, bytes.NewBuffer(data))
		if err != nil {
			log.Error("coral", "CreateIndex", err, "Creating Index.")
		}

	}

	return err
}

func doRequest(method string, urlStr string, payload io.Reader) error {

	var err error
	request, err := http.NewRequest(method, urlStr, payload)
	if err != nil {
		log.Error("coral", "doRequest", err, "New request")
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	var response *http.Response

	// only retry if there is a network Errorf

	// Retry retryTimes times if it fails to do the request
	for i := 0; i < retryTimes; i++ {
		response, err = client.Do(request)
		if err != nil {
			log.Error("coral", "doRequest", err, "Processing request")
		} else {
			defer response.Body.Close()
			if response.StatusCode != 200 {
				err = fmt.Errorf("Not succesful status code: %s.", response.Status)
				// wait and retry to do the request
				time.Sleep(250 * time.Millisecond)
			} else {
				break
			}
		}
	}
	return err
}
