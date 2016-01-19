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
package coral

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/ardanlabs/kit/cfg"
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

var endpoints map[string]string // model -> endpoint

func init() {

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)

	s := strategy.New()
	endpoints = s.GetPillarEndpoints()

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

	// Retry retryTimes times if it fails to do the request
	for i := 0; i < retryTimes; i++ {
		response, err = client.Do(request)
		if err != nil {
			log.Error("coral", "doRequest", err, "Processing request")
		} else {
			if response.StatusCode != 200 {
				err = fmt.Errorf("Not succesful status code: %s.", response.Status)
				log.Error("coral", "doRequest", err, "Processing request")
			} else {
				defer response.Body.Close()
				break
			}
		}

		// wait and retry to do the request
		time.Sleep(250 * time.Millisecond)
	}
	return err
}
