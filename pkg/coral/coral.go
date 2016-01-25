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
			log.Error("coral", "Addrow", err, "Error on sending request with row.")
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

	// only retry if there is a network Errorf

	// Retry retryTimes times if it fails to do the request
	for i := 0; i < retryTimes; i++ {
		response, err = client.Do(request)
		if err != nil {
			log.Error("coral", "doRequest", err, "Processing request")
		} else {
			if response.StatusCode != 200 {
				err = fmt.Errorf("Not succesful status code: %s.", response.Status)
				// wait and retry to do the request
				time.Sleep(250 * time.Millisecond)
			} else {
				defer response.Body.Close()
				break
			}
		}
	}
	return err
}
