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

// Config has the environment variables values
var Config struct {
	urlUser    string
	urlAsset   string
	urlComment string
}

func init() {

	logLevel := func() int {
		ll, err := cfg.Int("LOGGING_LEVEL")
		if err != nil {
			return log.USER
		}
		return ll
	}

	log.Init(os.Stderr, logLevel)
}

func setConfig() {

	// look for configuration urls
	Config.urlUser = os.Getenv("USER_URL")
	if Config.urlUser == "" {
		Config.urlUser = "http://localhost:8080/api/import/user"
		//log.Error("coral", "init", err, "Getting USER_URL env variable")
	}

	Config.urlAsset = os.Getenv("ASSET_URL")
	if Config.urlAsset == "" {
		Config.urlAsset = "http://localhost:8080/api/import/asset"
		//log.Error("coral", "init", err, "Getting ASSET_URL env variable")
	}

	Config.urlComment = os.Getenv("COMMENT_URL")
	if Config.urlComment == "" {
		Config.urlComment = "http://localhost:8080/api/import/comment"
		//log.Error("coral", "init", err, "Getting COMMENT_URL env variable")
	}
}

// AddRow send the row to pillar based on which collection is
func AddRow(data []byte, tableName string) error {

	setConfig()

	var err error

	switch tableName {
	case "user":
		err = addUser(data)
	case "comment":
		err = addComment(data)
	case "asset":
		err = addAsset(data)
	default:
		err = fmt.Errorf("No model %s in the coral systems.", tableName)
	}

	return err
}

// AddComment  calls the API to add the comment
func addComment(jcomment []byte) error {

	err := doRequest(methodPost, Config.urlComment, bytes.NewBuffer(jcomment))
	if err != nil {
		log.Error("coral", "addComment", err, "Error on sending request with comment.")
		return err
	}

	return nil
}

// AddAssets calls the API to add the asset
func addAsset(jasset []byte) error {

	err := doRequest(methodPost, Config.urlAsset, bytes.NewBuffer(jasset))
	if err != nil {
		return err
	}

	return nil
}

// AddUsers calls the API to add the user
func addUser(juser []byte) error {

	err := doRequest(methodPost, Config.urlUser, bytes.NewBuffer(juser))
	if err != nil {
		return err
	}

	return nil
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

	// Retry retryTimes times if it fails to do the request
	for i := 0; i < retryTimes; i++ {
		response, err := client.Do(request)
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
