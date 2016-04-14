package webservice

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ardanlabs/kit/log"
)

const (
	GET    string = "GET"
	POST   string = "POST"
	PUT    string = "PUT"
	DELETE string = "DELETE"

	retryTimes int = 3
)

//Response encapsulates a http response
type Response struct {
	Status     string
	Header     http.Header
	Body       string
	StatusCode int
}

func DoRequest(uuid string, method string, urlStr string, payload io.Reader) (*Response, error) {

	var err error
	request, err := http.NewRequest(method, urlStr, payload)
	if err != nil {
		log.Error(uuid, "coral.doRequest", err, "New http request.")
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	var response *http.Response

	// only retry if there is a network Errorf

	// Retry retryTimes times if it fails to do the request
	for i := 0; i < retryTimes; i++ {

		response, err = client.Do(request)

		if err != nil {
			log.Error(uuid, "coral.doRequest", err, "Sending request number %d to Pillar.", i)
		} else {

			defer response.Body.Close()
			resBody, _ := ioutil.ReadAll(response.Body)

			if response.StatusCode != 200 {
				err = fmt.Errorf("Not succesful status code: %s.", response.Status)
				// wait and retry to do the request
				time.Sleep(250 * time.Millisecond)
			} else {
				return &Response{
					response.Status,
					response.Header,
					string(resBody),
					response.StatusCode,
				}, nil
			}

		}
	}
	return nil, err
}
